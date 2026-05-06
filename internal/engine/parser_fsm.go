package engine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// stateFn defines a function representing a parser state.
// It takes an argParser pointer and returns the next state function.
type stateFn func(*argParser) stateFn

// argParser holds parsing state and context for command-line argument parsing.
type argParser struct {
	args  []string // CLI arguments to parse
	index int      // current index in args
	fs    *FlagSet // reference to the defined flag set
	out   []string // collected positional arguments
	err   error    // first error encountered, if any
	errs  []error  // collected errors (ContinueOnError)
}

// next returns the next argument and advances the index.
// If all args are consumed, it returns ok=false.
func (p *argParser) next() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		p.index++
		ok = true
	}
	return arg, ok
}

// peek returns the next argument without advancing the index.
func (p *argParser) peek() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		ok = true
	}
	return arg, ok
}

// run executes the parser state machine starting from stateStart.
// It continues until no further state is returned or an error is encountered.
func (p *argParser) run() error {
	state := argParserStateStart
	for state != nil {
		state = state(p)
		if p.err != nil {
			if p.fs.errorHandling == ContinueOnError {
				p.errs = append(p.errs, p.err)
				p.err = nil
				if state == nil {
					state = argParserStateStart
				}
				continue
			}
			return p.err
		}
	}
	if len(p.errs) > 0 {
		return errors.Join(p.errs...)
	}
	return nil
}

// runArgParserFSM initializes the argument parser and runs it.
// It returns any remaining positional arguments and a parsing error if any.
func runArgParserFSM(fs *FlagSet, args []string) ([]string, error) {
	p := &argParser{
		fs:   fs,
		args: args,
	}
	err := p.run()
	return p.out, err
}

// stateStart determines what kind of argument we're looking at (flag, positional, etc.)
// and returns the appropriate handler state function.
func argParserStateStart(p *argParser) stateFn {
	arg, ok := p.next()
	if !ok {
		return nil // no more arguments → done
	}

	switch {
	case arg == "--":
		// treat the rest as positional arguments
		p.out = append(p.out, p.args[p.index:]...)
		p.index = len(p.args)
		return nil

	case strings.HasPrefix(arg, "--"):
		// long form: --flag or --flag=value
		return argParserStateLong(arg)

	case strings.HasPrefix(arg, "-") && len(arg) > 1:
		// short form: -f, -abc, or -fvalue
		return argParserStateShort(arg)

	default:
		// positional argument
		p.out = append(p.out, arg)
		return argParserStateStart
	}
}

// handleUnknown routes an unknown flag through the configured handler.
func handleUnknown(p *argParser, name string) stateFn {
	if p.fs.unknownFlag == nil {
		p.err = fmt.Errorf("unknown flag: %s", name)
		return nil
	}
	if err := p.fs.unknownFlag(name); err != nil {
		p.err = err
		return nil
	}
	return argParserStateStart
}

// argParserStateLong handles arguments of the form --flag or --flag=value.
func argParserStateLong(arg string) stateFn {
	return func(p *argParser) stateFn {
		nameval := strings.TrimPrefix(arg, "--")
		name, val, hasVal := splitFlagArg(nameval)

		switch {
		case isDynamicFlag(name):
			return handleDynamic(name, val, hasVal, arg)

		case isKnownStaticFlag(p, name):
			return handleStatic(name, val, hasVal)

		default:
			return handleUnknown(p, "--"+name)
		}
	}
}

// handleDynamic processes a dynamic flag using the provided name, value, and format.
func handleDynamic(name, val string, hasVal bool, raw string) stateFn {
	return func(p *argParser) stateFn {
		item, id := lookupDynamic(p, name, raw)
		if p.err != nil {
			return nil
		}

		// Support non-strict bools like --group.id.flag
		if handled := tryDynBool(item, id); handled {
			return argParserStateStart
		}

		item.GetAny(name)

		// Handle --group.id.flag=value
		if hasVal {
			p.err = trySetDynamic(item, id, val, name)
			return argParserStateStart
		}

		// Handle --group.id.flag value
		if handled := handleDynamicValue(p, item, id, name); !handled {
			return nil
		}

		return argParserStateStart
	}
}

// handleDynamicValue parses and sets the next token as the value of a dynamic flag.
func handleDynamicValue(p *argParser, item core.DynamicValue, id, name string) bool {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		p.err = fmt.Errorf("missing value for flag: --%s", name)
		return false
	}

	p.next()
	p.err = trySetDynamic(item, id, next, name)
	return true
}

// handleStatic processes a known static flag by name.
func handleStatic(name, val string, hasVal bool) stateFn {
	return func(p *argParser) stateFn {
		fl := p.fs.staticFlagsMap[name]

		// Handle non-strict bools like -v
		if handled := tryBool(fl); handled {
			return argParserStateStart
		}
		// Support implicit increment if applicable
		if handled := tryCounter(p, fl); handled {
			return argParserStateStart
		}
		// Handle --flag=value
		if hasVal {
			p.err = trySet(fl.Value, val, "invalid value for flag --%s: %w", name)
			return argParserStateStart
		}
		// Handle --flag value
		if handled := tryLongValue(p, fl, name); handled {
			return argParserStateStart
		}

		p.err = fmt.Errorf("missing value for flag: --%s", name)
		return nil
	}
}

// argParserStateShort handles grouped short flags like -abc or -p8080.
func argParserStateShort(arg string) stateFn {
	return func(p *argParser) stateFn {
		shorts := strings.TrimPrefix(arg, "-")

		for i := 0; i < len(shorts); i++ {
			char := string(shorts[i])
			flag := findShortFlag(p.fs, char)
			if flag == nil {
				if next := handleUnknown(p, "-"+char); next == nil {
					return nil
				}
				continue
			}

			// Handle non-strict bools like -v
			if handled := tryBool(flag); handled {
				continue
			}
			// Handle counters like -vvvv
			if handled := tryCounter(p, flag); handled {
				continue
			}
			// Handle -p8080 (combined)
			if handled := tryShortCombined(p, flag, i, shorts, char); handled {
				break
			}

			// Handle: -p 8080
			p.err = tryShortValue(p, flag, char)
			break
		}

		return argParserStateStart
	}
}

// findShortFlag searches for a flag matching the given short letter.
func findShortFlag(fs *FlagSet, short string) *core.BaseFlag {
	for _, fl := range fs.staticFlagsMap {
		if fl.Short == short {
			return fl
		}
	}
	return nil
}

// tryBool attempts to set a non-strict boolean flag to true.
func tryBool(flag *core.BaseFlag) bool {
	if b, ok := flag.Value.(core.StrictBool); ok && !b.IsStrictBool() {
		flag.Value.Set("true") // nolint:errcheck
		return true
	}
	return false
}

// tryDynBool attempts to set a non-strict dynamic boolean flag to true.
func tryDynBool(item core.DynamicValue, id string) bool {
	if b, ok := item.(core.StrictBool); ok && !b.IsStrictBool() {
		item.Set(id, "true") // nolint
		return true
	}
	return false
}

// tryCounter calls Increment on a flag if it supports it.
func tryCounter(p *argParser, flag *core.BaseFlag) bool {
	if inc, ok := flag.Value.(core.Incrementable); ok {
		p.err = inc.Increment()
		return true
	}
	return false
}

// tryShortCombined parses short flag with immediate value like -p8080.
func tryShortCombined(p *argParser, flag *core.BaseFlag, i int, shorts string, char string) bool {
	if i < len(shorts)-1 {
		val := shorts[i+1:]
		p.err = trySet(flag.Value, val, "invalid value for flag -%s: %w", char)
		return true
	}
	return false
}

// tryLongValue consumes the next argument as a value for a long flag.
func tryLongValue(p *argParser, flag *core.BaseFlag, name string) bool {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return false
	}

	p.next()
	p.err = trySet(flag.Value, next, "invalid value for flag --%s: %w", name)
	return true
}

// tryShortValue consumes the next argument as a value for a short flag.
func tryShortValue(p *argParser, flag *core.BaseFlag, short string) error {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return fmt.Errorf("missing value for flag: -%s", flag.Short)
	}
	p.next()
	return trySet(flag.Value, next, "invalid value for flag -%s: %w", short)
}

// trySet attempts to set a value and wraps any error using the format string.
func trySet(value core.Value, input string, format string, label string) error {
	if err := value.Set(input); err != nil {
		return fmt.Errorf(format, label, err)
	}
	return nil
}

// trySetDynamic attempts to set a dynamic flag value and wraps any error.
func trySetDynamic(item core.DynamicValue, id, val, label string) error {
	if err := item.Set(id, val); err != nil {
		return fmt.Errorf("invalid value for dynamic flag --%s: %w", label, err)
	}
	return nil
}

// splitFlagArg splits --flag=value into ("flag", "value", true).
// Returns ("flag", "", false) if no '=' is present.
func splitFlagArg(s string) (name, val string, hasVal bool) {
	if i := strings.Index(s, "="); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, "", false
}

// isDynamicFlag checks if the flag name follows the group.id.flag pattern.
func isDynamicFlag(name string) bool {
	parts := strings.Split(name, ".")
	return len(parts) == 3 // group.id.field
}

// isKnownStaticFlag reports whether a flag is registered statically.
func isKnownStaticFlag(p *argParser, name string) bool {
	_, ok := p.fs.staticFlagsMap[name]
	return ok
}

// lookupDynamic locates a dynamic flag and its instance ID from the full name.
// Returns the dynamic value and ID, or sets p.err if lookup fails.
func lookupDynamic(p *argParser, name string, raw string) (core.DynamicValue, string) {
	parts := strings.Split(name, ".")
	if len(parts) != 3 {
		p.err = fmt.Errorf("invalid dynamic flag: --%s", name)
		return nil, ""
	}
	groupName, id, field := parts[0], parts[1], parts[2]

	group, ok := p.fs.dynamicGroupsMap[groupName]
	if !ok {
		p.err = fmt.Errorf("unknown dynamic group %q in flag %s", groupName, raw)
		return nil, ""
	}

	item, ok := group.Items()[field]
	if !ok {
		p.err = fmt.Errorf("unknown dynamic field %q in flag %s", field, raw)
		return nil, ""
	}

	return item.Value, id
}
