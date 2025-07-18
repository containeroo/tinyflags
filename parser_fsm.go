package tinyflags

import (
	"fmt"
	"strings"
)

// stateFn defines a function representing a parser state.
// It takes a parser pointer and returns the next state function.
type stateFn func(*parser) stateFn

// parser holds parsing state and context for command-line argument parsing.
type parser struct {
	args  []string // CLI arguments to parse
	index int      // current index in args
	fs    *FlagSet // reference to the defined flag set
	out   []string // collected positional arguments
	err   error    // first error encountered, if any
}

// next returns the next argument and advances the index.
// If all args are consumed, it returns ok=false.
func (p *parser) next() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		p.index++
		ok = true
	}
	return
}

// peek returns the next argument without advancing the index.
func (p *parser) peek() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		ok = true
	}
	return
}

// run executes the parser state machine starting from stateStartFn
// It continues until there are no more states (state == nil).
func (p *parser) run() error {
	state := stateStartFn
	for state != nil {
		state = state(p)
		if p.err != nil {
			return p.err
		}
	}
	return p.err
}

// parseArgsWithFSM initializes the parser and runs it.
// It returns any positional arguments and the first error encountered.
func parseArgsWithFSM(fs *FlagSet, args []string) ([]string, error) {
	p := &parser{
		fs:   fs,
		args: args,
	}
	err := p.run()
	return p.out, err
}

// stateStartFn determines what kind of argument we're looking at (flag, positional, etc.)
// and returns the appropriate handler state function.
func stateStartFn(p *parser) stateFn {
	arg, ok := p.next()
	if !ok {
		return nil // no more arguments → done
	}

	if arg == "--" {
		// "--" means treat the rest as positional arguments
		p.out = append(p.out, p.args[p.index:]...)
		p.index = len(p.args)
		return nil
	}

	if strings.HasPrefix(arg, "--") {
		// Long-form flag like --debug or --port=8080
		return stateLongFlagFn(arg)
	}

	if strings.HasPrefix(arg, "-") && len(arg) > 1 {
		// Short-form flag like -d or grouped flags like -abc
		return stateShortFlagFn(arg)
	}

	// Otherwise, it's a positional argument
	p.out = append(p.out, arg)
	return stateStartFn
}

// stateLongFlagFn handles arguments of the form --flag or --flag=value.
func stateLongFlagFn(arg string) stateFn {
	return func(p *parser) stateFn {
		nameval := strings.TrimPrefix(arg, "--")
		name, val, hasVal := splitFlagArg(nameval)

		switch {
		case isDynamicFlag(name):
			return handleDynamicFlag(name, val, hasVal)

		case isKnownStaticFlag(p, name):
			return handleStaticFlag(name, val, hasVal)

		default:
			p.err = fmt.Errorf("unknown flag: --%s", name)
			return nil
		}
	}
}

func handleDynamicFlag(name, val string, hasVal bool) stateFn {
	return func(p *parser) stateFn {
		parts := strings.Split(name, ".")
		group, id, field := parts[0], parts[1], parts[2]

		groupFields, ok := p.fs.dynamic[group]
		if !ok {
			p.err = fmt.Errorf("unknown dynamic group: --%s", name)
			return nil
		}
		item, ok := groupFields[field]
		if !ok {
			p.err = fmt.Errorf("unknown dynamic field: --%s", name)
			return nil
		}

		if hasVal {
			p.err = trySetDynamic(item, id, val, name)
			return stateStartFn
		}

		next, ok := p.peek()
		if !ok || strings.HasPrefix(next, "-") {
			p.err = fmt.Errorf("missing value for flag: --%s", name)
			return nil
		}

		p.next()
		p.err = trySetDynamic(item, id, next, name)
		return stateStartFn
	}
}

func handleStaticFlag(name, val string, hasVal bool) stateFn {
	return func(p *parser) stateFn {
		fl := p.fs.flags[name]

		// Handle booleans (with or without strict mode)
		if bv, ok := fl.value.(StrictBool); ok {
			if !bv.IsStrictBool() {
				fl.value.Set("true") // nolint:errcheck
				return stateStartFn
			}
		}

		if hasVal {
			p.err = trySet(fl.value, val, "invalid value for flag --%s: %s.", name)
			return stateStartFn
		}

		next, ok := p.peek()
		if !ok || strings.HasPrefix(next, "-") {
			p.err = fmt.Errorf("missing value for flag: --%s", name)
			return nil
		}

		p.next()
		p.err = trySet(fl.value, next, "invalid value for flag --%s: got %s.", name)
		return stateStartFn
	}
}

// stateShortFlagFn handles grouped short flags like -abc or single ones like -f value.
func stateShortFlagFn(arg string) stateFn {
	return func(p *parser) stateFn {
		shorts := strings.TrimPrefix(arg, "-")

		for i := 0; i < len(shorts); i++ {
			char := string(shorts[i])
			flag := findShortFlag(p.fs, char)
			if flag == nil {
				p.err = fmt.Errorf("unknown short flag: -%s", char)
				return nil
			}

			// case: -v   (non-strict bool)
			if handled := handleShortBool(flag); handled {
				continue
			}

			// case: -p8080 (combined)
			if i < len(shorts)-1 {
				val := shorts[i+1:]
				p.err = trySet(flag.value, val, "invalid value for flag -%s: %w", char)
				break
			}

			// case: -p 8080 (value is next arg)
			p.err = handleShortValue(p, flag, char)
			break
		}
		return stateStartFn
	}
}

func findShortFlag(fs *FlagSet, short string) *baseFlag {
	for _, fl := range fs.flags {
		if fl.short == short {
			return fl
		}
	}
	return nil
}

func handleShortBool(flag *baseFlag) bool {
	if b, ok := flag.value.(StrictBool); ok && !b.IsStrictBool() {
		flag.value.Set("true") // nolint:errcheck
		return true
	}
	return false
}

func handleShortValue(p *parser, flag *baseFlag, short string) error {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return fmt.Errorf("missing value for flag: -%s", short)
	}
	p.next()
	return trySet(flag.value, next, "invalid value for flag -%s: %w", short)
}

// trySet attempts to set the given value using input.
// If setting fails, it wraps the error using the provided format and label.
func trySet(value Value, input string, format string, label string) error {
	if err := value.Set(input); err != nil {
		return fmt.Errorf(format, label, err)
	}
	return nil
}

func trySetDynamic(item DynamicValue, id, val, label string) error {
	if err := item.Set(id, val); err != nil {
		return fmt.Errorf("invalid value for dynamic flag --%s: %w", label, err)
	}
	return nil
}

// splitFlagArg splits a string like "flag=value" into ("flag", "value", true).
// If there's no '=', it returns ("flag", "", false).
func splitFlagArg(s string) (name, val string, hasVal bool) {
	if i := strings.Index(s, "="); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, "", false
}

func isDynamicFlag(name string) bool {
	parts := strings.Split(name, ".")
	return len(parts) == 3 // group.id.field
}

func isKnownStaticFlag(p *parser, name string) bool {
	_, ok := p.fs.flags[name]
	return ok
}
