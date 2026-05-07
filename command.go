package tinyflags

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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

	name     string
	summary  string
	handling ErrorHandling
	parent   *Command
	globals  *FlagSet
	children map[string]*Command
	order    []*Command
	selected *Command
	builder  commandBuilder
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

// BuildCommand registers a zero-argument builder that returns the runnable for this command.
func (c *Command) BuildCommand(builder any) *Command {
	c.builder = wrapCommandBuilder(builder)
	return c
}

// Run registers one handler function plus deferred flag-backed arguments for this command.
func (c *Command) Run(handler any, bindings ...any) *Command {
	c.builder = wrapCommandRunner(handler, bindings...)
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

var (
	runnableType = reflect.TypeOf((*Runnable)(nil)).Elem()
	contextType  = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType    = reflect.TypeOf((*error)(nil)).Elem()
)

// wrapCommandBuilder validates one builder function and adapts it to the internal runner shape.
func wrapCommandBuilder(builder any) commandBuilder {
	if builder == nil {
		panic("tinyflags: command builder cannot be nil")
	}

	value := reflect.ValueOf(builder)
	typ := value.Type()
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("tinyflags: command builder must be a function, got %T", builder))
	}
	if typ.NumIn() != 0 {
		panic("tinyflags: command builder must not accept arguments")
	}
	if typ.NumOut() < 1 || typ.NumOut() > 2 {
		panic("tinyflags: command builder must return Runnable or (Runnable, error)")
	}
	if !typ.Out(0).Implements(runnableType) {
		panic(fmt.Sprintf("tinyflags: command builder must return a Runnable, got %s", typ.Out(0)))
	}
	if typ.NumOut() == 2 && !typ.Out(1).Implements(errorType) {
		panic(fmt.Sprintf("tinyflags: command builder second return must be an error, got %s", typ.Out(1)))
	}

	return func() (Runnable, error) {
		results := value.Call(nil)
		if len(results) == 2 && !results[1].IsZero() {
			return nil, results[1].Interface().(error)
		}
		if isNilBuilderResult(results[0]) {
			return nil, fmt.Errorf("tinyflags: command builder returned a nil Runnable")
		}

		runnable, ok := results[0].Interface().(Runnable)
		if !ok {
			return nil, fmt.Errorf("tinyflags: command builder returned %T, which does not implement Runnable", results[0].Interface())
		}
		return runnable, nil
	}
}

// wrapCommandRunner validates one handler plus bound arguments and adapts them to a runnable builder.
func wrapCommandRunner(handler any, bindings ...any) commandBuilder {
	value := reflect.ValueOf(handler)
	if !value.IsValid() {
		panic("tinyflags: command handler cannot be nil")
	}
	spec := parseRunHandler(value.Type())

	return func() (Runnable, error) {
		args, err := resolveRunBindings(spec, bindings)
		if err != nil {
			return nil, err
		}
		return commandHandlerRunner{
			handler:       value,
			args:          args,
			injectContext: spec.injectContext,
			returnsError:  spec.returnsError,
		}, nil
	}
}

type runHandlerSpec struct {
	paramTypes    []reflect.Type
	injectContext bool
	returnsError  bool
}

// parseRunHandler validates one registered command handler signature.
func parseRunHandler(handlerType reflect.Type) runHandlerSpec {
	if handlerType == nil || handlerType.Kind() != reflect.Func {
		panic("tinyflags: command handler must be a function")
	}

	spec := runHandlerSpec{}
	start := 0
	if handlerType.NumIn() > 0 && handlerType.In(0).Implements(contextType) {
		spec.injectContext = true
		start = 1
	}
	for i := start; i < handlerType.NumIn(); i++ {
		spec.paramTypes = append(spec.paramTypes, handlerType.In(i))
	}

	if handlerType.NumOut() > 1 {
		panic("tinyflags: command handler must return no values or one error")
	}
	if handlerType.NumOut() == 1 {
		if !handlerType.Out(0).Implements(errorType) {
			panic(fmt.Sprintf("tinyflags: command handler return type must be error, got %s", handlerType.Out(0)))
		}
		spec.returnsError = true
	}

	return spec
}

type commandHandlerRunner struct {
	handler       reflect.Value
	args          []reflect.Value
	injectContext bool
	returnsError  bool
}

// Run executes one registered command handler with its parsed argument values.
func (r commandHandlerRunner) Run(ctx context.Context) error {
	callArgs := make([]reflect.Value, 0, len(r.args)+1)
	if r.injectContext {
		callArgs = append(callArgs, reflect.ValueOf(ctx))
	}
	callArgs = append(callArgs, r.args...)

	results := r.handler.Call(callArgs)
	if !r.returnsError || len(results) == 0 || results[0].IsZero() {
		return nil
	}
	return results[0].Interface().(error)
}

// resolveRunBindings freezes one handler invocation's bound arguments after parsing.
func resolveRunBindings(spec runHandlerSpec, bindings []any) ([]reflect.Value, error) {
	if len(bindings) != len(spec.paramTypes) {
		return nil, fmt.Errorf("tinyflags: command handler expects %d bound arguments, got %d", len(spec.paramTypes), len(bindings))
	}

	args := make([]reflect.Value, 0, len(bindings))
	for i, binding := range bindings {
		value, err := resolveRunBinding(binding, spec.paramTypes[i])
		if err != nil {
			return nil, fmt.Errorf("tinyflags: binding %d: %w", i+1, err)
		}
		args = append(args, value)
	}
	return args, nil
}

// resolveRunBinding converts one registered binding into the concrete parameter value a handler needs.
func resolveRunBinding(binding any, paramType reflect.Type) (reflect.Value, error) {
	if binding == nil {
		return reflect.Value{}, fmt.Errorf("nil binding is not supported for %s", paramType)
	}

	value := reflect.ValueOf(binding)
	if resolved, ok := assignRunBinding(value, paramType); ok {
		return freezeRunBindingValue(resolved), nil
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}, fmt.Errorf("nil pointer binding cannot satisfy %s", paramType)
		}
		if resolved, ok := assignRunBinding(value.Elem(), paramType); ok {
			return freezeRunBindingValue(resolved), nil
		}
	}

	return reflect.Value{}, fmt.Errorf("cannot use %s as %s", value.Type(), paramType)
}

// assignRunBinding normalizes one reflect value to the handler parameter type when compatible.
func assignRunBinding(value reflect.Value, paramType reflect.Type) (reflect.Value, bool) {
	if value.Type().AssignableTo(paramType) {
		return value, true
	}
	if value.Type().ConvertibleTo(paramType) {
		return value.Convert(paramType), true
	}
	return reflect.Value{}, false
}

// freezeRunBindingValue copies one resolved binding so later flag mutations do not leak into one parsed runner.
func freezeRunBindingValue(value reflect.Value) reflect.Value {
	frozen := reflect.New(value.Type()).Elem()
	frozen.Set(value)
	return frozen
}

// isNilBuilderResult reports whether a reflect value contains a nil value for nilable kinds.
func isNilBuilderResult(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
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
	err := fs.Parse([]string{"--help"})
	if err == nil {
		return ""
	}
	return err.Error()
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
