package tinyflags

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// Runnable represents a parsed command that can execute with cancellation support.
type Runnable interface {
	Run(context.Context) error
}

// Runner is the preferred name for one parsed command executor.
type Runner = Runnable

// Command represents a command or subcommand with local and persistent flags.
type Command struct {
	*FlagSet

	name         string
	summary      string
	handling     ErrorHandling
	requireChild bool
	parent       *Command
	globals      *FlagSet
	children     map[string]*Command
	order        []*Command
	selected     *Command
	builder      commandBuilder
}

type commandBuilder func() (Runnable, error)

// NewCommand creates a new root command.
func NewCommand(name string, handling ErrorHandling) *Command {
	local := NewFlagSet(name, handling)
	return &Command{
		FlagSet:  local,
		name:     name,
		handling: handling,
		globals:  local,
		children: make(map[string]*Command),
	}
}

// Command creates a nested subcommand with the given name and summary.
func (c *Command) Command(name string, summary string) *Command {
	fullName := c.FullName() + " " + name
	child := &Command{
		FlagSet:  NewFlagSet(fullName, c.handling),
		name:     name,
		summary:  summary,
		handling: c.handling,
		parent:   c,
		globals:  NewFlagSet(fullName, c.handling),
		children: make(map[string]*Command),
	}
	c.children[name] = child
	c.order = append(c.order, child)
	return child
}

// Globals returns the persistent flag set for this command subtree.
func (c *Command) Globals() *FlagSet {
	return c.globals
}

// HelpText renders help for the selected command when available, otherwise the receiver.
func (c *Command) HelpText() string {
	target := c
	if c != nil && c.selected != nil {
		target = c.selected
	}
	if target == nil {
		return ""
	}
	return renderCommandHelp(target)
}

// WriteHelp writes the rendered help text to w.
func (c *Command) WriteHelp(w io.Writer) error {
	if w == nil {
		return nil
	}
	_, err := io.WriteString(w, c.HelpText())
	return err
}

// RequireCommand enforces that one direct or nested child command must be selected.
func (c *Command) RequireCommand() *Command {
	c.requireChild = true
	return c
}

// BuildCommand registers a zero-argument builder that returns the runnable for this command.
func (c *Command) BuildCommand(builder any) *Command {
	c.setBuilder(wrapCommandBuilder(builder))
	return c
}

// Run registers one handler function plus deferred flag-backed arguments for this command.
func (c *Command) Run(handler any, bindings ...any) *Command {
	c.setBuilder(wrapCommandRunner(handler, bindings...))
	return c
}

// SelectedCommand returns the leaf command selected during the last parse.
func (c *Command) SelectedCommand() *Command {
	return c.selected
}

// Commands returns child commands in registration order.
func (c *Command) Commands() []*Command {
	return c.order
}

// Summary returns the short summary used in command listings.
func (c *Command) Summary() string {
	return c.summary
}

// Name returns the command segment name.
func (c *Command) Name() string {
	return c.name
}

// FullName returns the full command path.
func (c *Command) FullName() string {
	if c.parent == nil {
		return c.name
	}
	return c.parent.FullName() + " " + c.name
}

// Parse selects a command path and parses matching local and persistent flags.
func (c *Command) Parse(args []string) error {
	c.selected = c
	state := commandParseState{
		argsBySet: make(map[*FlagSet][]string),
	}

	current := c
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--help" || arg == "-h" {
			state.helpTarget = current
			continue
		}
		if arg == "--version" {
			state.versionRequested = true
			continue
		}

		if child, ok := current.children[arg]; ok && !strings.HasPrefix(arg, "-") {
			// Bare child names advance the command cursor instead of becoming positionals.
			current = child
			continue
		}

		if strings.HasPrefix(arg, "--") && arg != "--" {
			owner, flag := current.resolveLongArg(arg)
			if owner == nil || flag == nil {
				state.append(ownerOrCurrent(owner, current), arg)
				continue
			}

			state.append(owner, arg)
			// Route the following token with the same owner when the flag consumes a value.
			if !strings.Contains(arg, "=") && flagConsumesValue(flag) && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				state.append(owner, args[i])
			}
			continue
		}

		if strings.HasPrefix(arg, "-") && arg != "-" {
			consumed := current.routeShortArgs(arg, args, &i, &state)
			if consumed {
				continue
			}
		}

		state.append(current.FlagSet, arg)
	}

	c.selected = current
	if state.versionRequested {
		return c.FlagSet.Parse([]string{"--version"})
	}
	if state.helpTarget != nil {
		return RequestHelp(renderCommandHelp(state.helpTarget))
	}

	var errs []error
	for _, cmd := range c.commandPathTo(current) {
		for _, fs := range cmd.parseScopes() {
			if err := fs.Parse(state.argsBySet[fs]); err != nil {
				errs = append(errs, err)
				if c.handling != ContinueOnError {
					return err
				}
			}
		}
	}
	if err := c.missingRequiredCommand(current); err != nil {
		errs = append(errs, err)
		if c.handling != ContinueOnError {
			return err
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// ParseRunnable parses args and builds the runnable for the selected command.
func (c *Command) ParseRunnable(args []string) (Runnable, error) {
	return c.ParseRunner(args)
}

// ParseRunner parses args and builds the runner for the selected command.
func (c *Command) ParseRunner(args []string) (Runner, error) {
	if err := c.Parse(args); err != nil {
		return nil, err
	}

	selected := c.SelectedCommand()
	if selected == nil {
		return nil, fmt.Errorf("tinyflags: no command selected")
	}
	if selected.builder == nil {
		return nil, fmt.Errorf("tinyflags: no command runner registered for command %q", selected.FullName())
	}

	return selected.builder()
}

type commandParseState struct {
	argsBySet        map[*FlagSet][]string
	helpTarget       *Command
	versionRequested bool
}

// setBuilder stores the internal runner builder shared by all registration styles.
func (c *Command) setBuilder(builder commandBuilder) {
	c.builder = builder
}

// append records one routed argument for a specific flag set.
func (s *commandParseState) append(fs *FlagSet, arg string) {
	if fs == nil {
		return
	}
	s.argsBySet[fs] = append(s.argsBySet[fs], arg)
}

// parseScopes returns the flag sets that must parse for this command node.
func (c *Command) parseScopes() []*FlagSet {
	if c.parent == nil || c.globals == c.FlagSet {
		return []*FlagSet{c.FlagSet}
	}
	return []*FlagSet{c.globals, c.FlagSet}
}

// commandPathTo returns the command chain from the receiver to the target.
func (c *Command) commandPathTo(target *Command) []*Command {
	var path []*Command
	for cur := target; cur != nil; cur = cur.parent {
		path = append(path, cur)
		if cur == c {
			break
		}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func (c *Command) missingRequiredCommand(selected *Command) error {
	for _, cmd := range c.commandPathTo(selected) {
		if cmd == nil || !cmd.requireChild {
			continue
		}
		if cmd == selected && len(cmd.order) > 0 {
			return &UsageError{
				Err:  &CommandRequired{Command: cmd.FullName()},
				Help: renderCommandHelp(cmd),
			}
		}
	}
	return nil
}

// availableFlagSets returns local and inherited persistent flag sets in lookup order.
func (c *Command) availableFlagSets() []*FlagSet {
	scopes := []*FlagSet{c.FlagSet}
	if c.parent != nil && c.globals != c.FlagSet {
		scopes = append(scopes, c.globals)
	}
	for parent := c.parent; parent != nil; parent = parent.parent {
		if parent.globals != nil {
			scopes = append(scopes, parent.globals)
		}
	}
	return scopes
}

// resolveLongArg finds which flag set owns a long-form flag token.
func (c *Command) resolveLongArg(arg string) (*FlagSet, *core.BaseFlag) {
	name := strings.TrimPrefix(arg, "--")
	if eq := strings.Index(name, "="); eq >= 0 {
		name = name[:eq]
	}
	parts := strings.Split(name, ".")

	for _, fs := range c.availableFlagSets() {
		if fl := fs.impl.LookupFlag(name); fl != nil {
			return fs, fl
		}
		if len(parts) == 3 {
			if fl := lookupDynamicFlag(fs, parts[0], parts[2]); fl != nil {
				return fs, fl
			}
		}
	}
	return nil, nil
}

// resolveShortFlag finds which flag set owns a short-form flag token.
func (c *Command) resolveShortFlag(short string) (*FlagSet, *core.BaseFlag) {
	for _, fs := range c.availableFlagSets() {
		for _, fl := range fs.impl.OrderedStaticFlags() {
			if fl.Short == short {
				return fs, fl
			}
		}
	}
	return nil, nil
}

// routeShortArgs routes grouped short options to the owning flag sets.
func (c *Command) routeShortArgs(arg string, args []string, i *int, state *commandParseState) bool {
	shorts := strings.TrimPrefix(arg, "-")
	if shorts == "" {
		return false
	}

	for idx := 0; idx < len(shorts); idx++ {
		short := string(shorts[idx])
		owner, fl := c.resolveShortFlag(short)
		if owner == nil || fl == nil {
			state.append(c.FlagSet, arg)
			return true
		}

		if flagConsumesValue(fl) {
			if idx < len(shorts)-1 {
				state.append(owner, "-"+short+shorts[idx+1:])
				return true
			}
			state.append(owner, "-"+short)
			if *i+1 < len(args) && !strings.HasPrefix(args[*i+1], "-") {
				*i++
				state.append(owner, args[*i])
			}
			return true
		}

		state.append(owner, "-"+short)
	}

	return true
}

// flagConsumesValue reports whether a flag expects a following value token.
func flagConsumesValue(fl *core.BaseFlag) bool {
	if fl == nil || fl.Value == nil {
		return false
	}
	if _, ok := fl.Value.(core.Incrementable); ok {
		return false
	}
	if b, ok := fl.Value.(core.StrictBool); ok && !b.IsStrictBool() {
		return false
	}
	return true
}

// lookupDynamicFlag resolves a dynamic field inside one dynamic group.
func lookupDynamicFlag(fs *FlagSet, groupName, field string) *core.BaseFlag {
	for _, group := range fs.impl.DynamicGroups() {
		if group.Name() != groupName {
			continue
		}
		return group.LookupFlag(field)
	}
	return nil
}

// ownerOrCurrent falls back to the current command's local flag set when unset.
func ownerOrCurrent(owner *FlagSet, current *Command) *FlagSet {
	if owner != nil {
		return owner
	}
	return current.FlagSet
}

// renderCommandHelp renders local help plus any child command listing.
func renderCommandHelp(cmd *Command) string {
	helpText := renderLocalHelp(cmd.FlagSet)
	lines := strings.Split(strings.TrimRight(helpText, "\n"), "\n")
	if len(cmd.order) > 0 && len(lines) > 0 {
		lines[0] = renderUsageLine(cmd)
	}

	var b strings.Builder
	b.WriteString(strings.Join(lines, "\n"))
	if len(cmd.order) > 0 {
		b.WriteString("\n\nCommands:\n")
		width := longestCommandName(cmd.order)
		for _, child := range cmd.order {
			fmt.Fprintf(&b, "  %-*s  %s\n", width, child.name, child.summary)
		}
	}
	b.WriteString("\n")
	return b.String()
}

// renderLocalHelp renders help text for one command-local flag set.
func renderLocalHelp(fs *FlagSet) string {
	if fs == nil || fs.impl == nil {
		return ""
	}
	return fs.impl.RenderHelpText()
}

// renderUsageLine builds the usage line for a command help screen.
func renderUsageLine(cmd *Command) string {
	usage := "Usage: " + cmd.FullName()
	if hasAnyVisibleFlags(cmd.FlagSet) {
		usage += " [flags]"
	}
	if len(cmd.order) > 0 {
		usage += " <command>"
	}
	return usage
}

// hasAnyVisibleFlags reports whether a local command help screen has visible flags.
func hasAnyVisibleFlags(fs *FlagSet) bool {
	for _, fl := range fs.impl.OrderedStaticFlags() {
		if !fl.Hidden {
			return true
		}
	}
	for _, group := range fs.impl.DynamicGroups() {
		if group.IsHidden() {
			continue
		}
		for _, fl := range group.DynamicFlags() {
			if !fl.Hidden {
				return true
			}
		}
	}
	return false
}

// longestCommandName returns the widest child command name for aligned help output.
func longestCommandName(commands []*Command) int {
	width := 0
	for _, cmd := range commands {
		if len(cmd.name) > width {
			width = len(cmd.name)
		}
	}
	return width
}
