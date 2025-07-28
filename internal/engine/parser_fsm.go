package engine

import (
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
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
	state := stateStart
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

// stateStart determines what kind of argument we're looking at (flag, positional, etc.)
// and returns the appropriate handler state function.
func stateStart(p *parser) stateFn {
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
		return stateLong(arg)

	case strings.HasPrefix(arg, "-") && len(arg) > 1:
		// short form: -f, -abc, or -fvalue
		return stateShort(arg)

	default:
		// positional argument
		p.out = append(p.out, arg)
		return stateStart
	}
}

// stateLong handles arguments of the form --flag or --flag=value.
func stateLong(arg string) stateFn {
	return func(p *parser) stateFn {
		nameval := strings.TrimPrefix(arg, "--")
		name, val, hasVal := splitFlagArg(nameval)

		switch {
		case isDynamicFlag(name):
			return handleDynamic(name, val, hasVal)

		case isKnownStaticFlag(p, name):
			return handleStatic(name, val, hasVal)

		default:
			p.err = fmt.Errorf("unknown flag: --%s", name)
			return nil
		}
	}
}

func handleDynamic(name, val string, hasVal bool) stateFn {
	return func(p *parser) stateFn {
		item, id := lookupDynamic(p, name)
		if p.err != nil {
			return nil
		}

		// Handle non-strict bools like --group.id.flag
		if handled := tryDynBool(item, id); handled {
			return stateStart
		}

		// Handle --group.id.flag=value
		if hasVal {
			p.err = trySetDynamic(item, id, val, name)
			return stateStart
		}

		// Hanlde --group.id.flag value
		if handled := handleDynamicValue(p, item, id, name); !handled {
			return nil
		}

		return stateStart
	}
}

func handleDynamicValue(p *parser, item core.DynamicValue, id, name string) bool {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		p.err = fmt.Errorf("missing value for flag: --%s", name)
		return false
	}

	p.next()
	p.err = trySetDynamic(item, id, next, name)
	return true
}

func handleStatic(name, val string, hasVal bool) stateFn {
	return func(p *parser) stateFn {
		fl := p.fs.staticFlagsMap[name]

		// Handle non-strict bools like -v
		if handled := tryBool(fl); handled {
			return stateStart
		}

		// Support implicit increment if applicable
		if handled := tryCounter(p, fl); handled {
			return stateStart
		}

		// Handle --flag=value
		if hasVal {
			p.err = trySet(fl.Value, val, "invalid value for flag --%s: %s.", name)
			return stateStart
		}

		// Handle --flag value
		if handled := tryLongValue(p, fl, name); handled {
			return stateStart
		}

		p.err = fmt.Errorf("missing value for flag: --%s", name)
		return nil
	}
}

// stateShort handles grouped short flags like -abc or single ones like -f value.
func stateShort(arg string) stateFn {
	return func(p *parser) stateFn {
		shorts := strings.TrimPrefix(arg, "-")

		for i := 0; i < len(shorts); i++ {
			char := string(shorts[i])
			flag := findShortFlag(p.fs, char)
			if flag == nil {
				p.err = fmt.Errorf("unknown short flag: -%s", char)
				return nil
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

		return stateStart
	}
}

func findShortFlag(fs *FlagSet, short string) *core.BaseFlag {
	for _, fl := range fs.staticFlagsMap {
		if fl.Short == short {
			return fl
		}
	}
	return nil
}

func tryBool(flag *core.BaseFlag) bool {
	if b, ok := flag.Value.(core.StrictBool); ok && !b.IsStrictBool() {
		flag.Value.Set("true") // nolint:errcheck
		return true
	}
	return false
}

func tryDynBool(item core.DynamicValue, id string) bool {
	if b, ok := item.(core.StrictBool); ok && !b.IsStrictBool() {
		item.Set(id, "true") // nolint
		return true
	}
	return false
}

func tryCounter(p *parser, flag *core.BaseFlag) bool {
	if inc, ok := flag.Value.(core.Incrementable); ok {
		p.err = inc.Increment()
		return true
	}
	return false
}

func tryShortCombined(p *parser, flag *core.BaseFlag, i int, shorts string, char string) bool {
	if i < len(shorts)-1 {
		val := shorts[i+1:]
		p.err = trySet(flag.Value, val, "invalid value for flag -%s: %w", char)
		return true
	}
	return false
}

func tryLongValue(p *parser, flag *core.BaseFlag, name string) bool {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return false
	}

	p.next()
	p.err = trySet(flag.Value, next, "invalid value for flag --%s: got %s.", name)
	return true
}

func tryShortValue(p *parser, flag *core.BaseFlag, short string) error {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return fmt.Errorf("missing value for flag: -%s", flag.Short)
	}
	p.next()
	return trySet(flag.Value, next, "invalid value for flag -%s: %w", short)
}

// trySet attempts to set the given value using input.
// If setting fails, it wraps the error using the provided format and label.
func trySet(value core.Value, input string, format string, label string) error {
	if err := value.Set(input); err != nil {
		return fmt.Errorf(format, label, err)
	}
	return nil
}

func trySetDynamic(item core.DynamicValue, id, val, label string) error {
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
	_, ok := p.fs.staticFlagsMap[name]
	return ok
}

// lookupDynamic locates the dynamic item and returns its parser and ID.
func lookupDynamic(p *parser, name string) (core.DynamicValue, string) {
	parts := strings.Split(name, ".")
	if len(parts) != 3 {
		p.err = fmt.Errorf("invalid dynamic flag: --%s", name)
		return nil, ""
	}
	groupName, id, field := parts[0], parts[1], parts[2]

	group, ok := p.fs.dynamicGroupsMap[groupName]
	if !ok {
		p.err = fmt.Errorf("unknown dynamic group: %s", name)
		return nil, ""
	}

	item, ok := group.Items()[field]
	if !ok {
		p.err = fmt.Errorf("unknown dynamic field: %s", name)
		return nil, ""
	}

	return item.Value, id
}
