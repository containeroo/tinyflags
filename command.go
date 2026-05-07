package tinyflags

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// Command represents a command or subcommand with local and persistent flags.
type Command struct {
	*FlagSet

	name     string
	summary  string
	handling ErrorHandling
	parent   *Command
	globals  *FlagSet
	children map[string]*Command
	order    []*Command
	selected *Command
}

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

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

type commandParseState struct {
	argsBySet        map[*FlagSet][]string
	helpTarget       *Command
	versionRequested bool
}

func (s *commandParseState) append(fs *FlagSet, arg string) {
	if fs == nil {
		return
	}
	s.argsBySet[fs] = append(s.argsBySet[fs], arg)
}

func (c *Command) parseScopes() []*FlagSet {
	if c.parent == nil || c.globals == c.FlagSet {
		return []*FlagSet{c.FlagSet}
	}
	return []*FlagSet{c.globals, c.FlagSet}
}

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

func lookupDynamicFlag(fs *FlagSet, groupName, field string) *core.BaseFlag {
	for _, group := range fs.impl.DynamicGroups() {
		if group.Name() != groupName {
			continue
		}
		return group.LookupFlag(field)
	}
	return nil
}

func ownerOrCurrent(owner *FlagSet, current *Command) *FlagSet {
	if owner != nil {
		return owner
	}
	return current.FlagSet
}

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

func renderLocalHelp(fs *FlagSet) string {
	err := fs.Parse([]string{"--help"})
	if err == nil {
		return ""
	}
	return err.Error()
}

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

func longestCommandName(commands []*Command) int {
	width := 0
	for _, cmd := range commands {
		if len(cmd.name) > width {
			width = len(cmd.name)
		}
	}
	return width
}
