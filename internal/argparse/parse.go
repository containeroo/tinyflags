package argparse

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// Config supplies the callbacks and behavior needed by the argument parser.
type Config struct {
	ContinueOnError   bool                                                    // Continue parsing and join errors instead of stopping at the first error.
	LookupStaticFlag  func(string) *core.BaseFlag                             // Resolves a long static flag name.
	LookupShortFlag   func(string) *core.BaseFlag                             // Resolves a short static flag name.
	LookupDynamicFlag func(string, string) (core.DynamicValue, string, error) // Resolves a dynamic flag and its ID.
	HandleUnknownFlag func(string) error                                      // Handles an unknown flag when one is encountered.
	RecordOrigin      func(name, key string)                                  // Records the source key that set a flag.
}

// stateFn consumes parser input and returns the next state.
type stateFn func(*parser) stateFn

// parser holds one stateful argument parsing pass.
type parser struct {
	config Config   // Callbacks and behavior supplied by the owning flag set.
	args   []string // Input arguments being parsed.
	index  int      // Index of the next unconsumed input argument.
	out    []string // Positional arguments collected during parsing.
	err    error    // Error produced by the current parser state.
	errs   []error  // Errors accumulated when continuation is enabled.
}

// Parse tokenizes args and applies callbacks to populate flag values.
func Parse(config Config, args []string) ([]string, error) {
	p := &parser{
		config: config,
		args:   args,
	}
	err := p.run()
	return p.out, err
}

// next returns and consumes the next input argument.
func (p *parser) next() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		p.index++
		ok = true
	}
	return arg, ok
}

// peek returns the next input argument without consuming it.
func (p *parser) peek() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		ok = true
	}
	return arg, ok
}

// run advances the state machine until input is exhausted or parsing fails.
func (p *parser) run() error {
	state := stateStart
	for state != nil {
		state = state(p)
		if p.err != nil {
			if p.config.ContinueOnError {
				p.errs = append(p.errs, p.err)
				p.err = nil
				if state == nil {
					state = stateStart
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

// stateStart classifies the next token and dispatches to the appropriate parser state.
func stateStart(p *parser) stateFn {
	arg, ok := p.next()
	if !ok {
		return nil
	}

	switch {
	case arg == "--":
		p.out = append(p.out, p.args[p.index:]...)
		p.index = len(p.args)
		return nil
	case strings.HasPrefix(arg, "--"):
		return stateLong(arg)
	case strings.HasPrefix(arg, "-") && len(arg) > 1:
		return stateShort(arg)
	default:
		p.out = append(p.out, arg)
		return stateStart
	}
}

// handleUnknown delegates an unknown flag to the optional callback.
func handleUnknown(p *parser, name string) stateFn {
	if p.config.HandleUnknownFlag == nil {
		p.err = fmt.Errorf("unknown flag %s", name)
		return nil
	}
	if err := p.config.HandleUnknownFlag(name); err != nil {
		p.err = err
		return nil
	}
	return stateStart
}

// stateLong parses one long flag token.
func stateLong(arg string) stateFn {
	return func(p *parser) stateFn {
		nameval := strings.TrimPrefix(arg, "--")
		name, val, hasVal := splitFlagArg(nameval)

		switch {
		case isDynamicFlag(name):
			return handleDynamic(name, val, hasVal, arg)
		case p.config.LookupStaticFlag(name) != nil:
			return handleStatic(name, val, hasVal)
		default:
			return handleUnknown(p, "--"+name)
		}
	}
}

// handleDynamic resolves and sets one dynamic flag.
func handleDynamic(name, val string, hasVal bool, raw string) stateFn {
	return func(p *parser) stateFn {
		item, id, err := p.config.LookupDynamicFlag(name, raw)
		if err != nil {
			p.err = err
			return nil
		}

		if handled := tryDynamicBool(item, id); handled {
			recordOrigin(p, name, "--"+name)
			return stateStart
		}

		if hasVal {
			p.err = trySetDynamic(item, id, val, name)
			if p.err == nil {
				recordOrigin(p, name, "--"+name)
			}
			return stateStart
		}

		if handled := handleDynamicValue(p, item, id, name); !handled {
			return nil
		}
		if p.err == nil {
			recordOrigin(p, name, "--"+name)
		}

		return stateStart
	}
}

// handleDynamicValue consumes and sets a separate value for a dynamic flag.
func handleDynamicValue(p *parser, item core.DynamicValue, id, name string) bool {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		p.err = fmt.Errorf("missing value for flag --%s", name)
		return false
	}

	p.next()
	p.err = trySetDynamic(item, id, next, name)
	return true
}

// handleStatic resolves and sets one static flag.
func handleStatic(name, val string, hasVal bool) stateFn {
	return func(p *parser) stateFn {
		flag := p.config.LookupStaticFlag(name)

		if handled := tryBool(flag); handled {
			recordOrigin(p, flag.Name, "--"+name)
			return stateStart
		}
		if handled := tryCounter(p, flag); handled {
			if p.err == nil {
				recordOrigin(p, flag.Name, "--"+name)
			}
			return stateStart
		}
		if hasVal {
			p.err = trySet(flag.Value, val, "invalid value for flag --%s: %w", name)
			if p.err == nil {
				recordOrigin(p, flag.Name, "--"+name)
			}
			return stateStart
		}
		if handled := tryLongValue(p, flag, name); handled {
			if p.err == nil {
				recordOrigin(p, flag.Name, "--"+name)
			}
			return stateStart
		}

		p.err = fmt.Errorf("missing value for flag --%s", name)
		return nil
	}
}

// stateShort parses one or more combined short flag tokens.
func stateShort(arg string) stateFn {
	return func(p *parser) stateFn {
		shorts := strings.TrimPrefix(arg, "-")

		for i := 0; i < len(shorts); i++ {
			char := string(shorts[i])
			flag := p.config.LookupShortFlag(char)
			if flag == nil {
				if next := handleUnknown(p, "-"+char); next == nil {
					return nil
				}
				continue
			}

			if handled := tryBool(flag); handled {
				recordOrigin(p, flag.Name, "-"+char)
				continue
			}
			if handled := tryCounter(p, flag); handled {
				if p.err == nil {
					recordOrigin(p, flag.Name, "-"+char)
				}
				continue
			}
			if handled := tryShortCombined(p, flag, i, shorts, char); handled {
				if p.err == nil {
					recordOrigin(p, flag.Name, "-"+char)
				}
				break
			}

			p.err = tryShortValue(p, flag, char)
			if p.err == nil {
				recordOrigin(p, flag.Name, "-"+char)
			}
			break
		}

		return stateStart
	}
}

// tryBool applies the implicit true value for a non-strict boolean flag.
func tryBool(flag *core.BaseFlag) bool {
	if flag == nil {
		return false
	}
	if b, ok := flag.Value.(core.StrictBool); ok && !b.IsStrictBool() {
		flag.Value.Set("true") // nolint:errcheck
		return true
	}
	return false
}

// tryDynamicBool applies the implicit true value for a non-strict dynamic boolean.
func tryDynamicBool(item core.DynamicValue, id string) bool {
	if b, ok := item.(core.StrictBool); ok && !b.IsStrictBool() {
		item.Set(id, "true") // nolint:errcheck
		return true
	}
	return false
}

// tryCounter increments a counter flag when it supports incrementing.
func tryCounter(p *parser, flag *core.BaseFlag) bool {
	if inc, ok := flag.Value.(core.Incrementable); ok {
		p.err = inc.Increment()
		return true
	}
	return false
}

// tryShortCombined treats the remainder of a short flag group as this flag's value.
func tryShortCombined(p *parser, flag *core.BaseFlag, i int, shorts string, char string) bool {
	if i < len(shorts)-1 {
		val := shorts[i+1:]
		p.err = trySet(flag.Value, val, "invalid value for flag -%s: %w", char)
		return true
	}
	return false
}

// tryLongValue consumes the following token as a long flag's value.
func tryLongValue(p *parser, flag *core.BaseFlag, name string) bool {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return false
	}

	p.next()
	p.err = trySet(flag.Value, next, "invalid value for flag --%s: %w", name)
	return true
}

// tryShortValue consumes the following token as a short flag's value.
func tryShortValue(p *parser, flag *core.BaseFlag, short string) error {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return fmt.Errorf("missing value for flag -%s", flag.Short)
	}
	p.next()
	return trySet(flag.Value, next, "invalid value for flag -%s: %w", short)
}

// trySet assigns input to a static value and decorates validation errors.
func trySet(value core.Value, input string, format string, label string) error {
	if err := value.Set(input); err != nil {
		return fmt.Errorf(format, label, err)
	}
	return nil
}

// trySetDynamic assigns val to one dynamic value and decorates validation errors.
func trySetDynamic(item core.DynamicValue, id, val, label string) error {
	if err := item.Set(id, val); err != nil {
		return fmt.Errorf("invalid value for flag --%s: %w", label, err)
	}
	return nil
}

// recordOrigin reports a successful flag assignment to the owning parser.
func recordOrigin(p *parser, name, key string) {
	if p.config.RecordOrigin != nil {
		p.config.RecordOrigin(name, key)
	}
}

// splitFlagArg separates a flag name from an optional equals-delimited value.
func splitFlagArg(s string) (name, val string, hasVal bool) {
	if i := strings.Index(s, "="); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, "", false
}

// isDynamicFlag reports whether name has the group.ID.field shape.
func isDynamicFlag(name string) bool {
	return len(strings.Split(name, ".")) == 3
}
