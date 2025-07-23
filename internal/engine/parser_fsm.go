package engine

import (
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// stateFn represents one step of the parse state machine.
// It returns the next stateFn, or nil to finish.
type stateFn func(*parser) stateFn

// parser holds the mutable state while parsing CLI args.
type parser struct {
	args  []string // all CLI arguments
	index int      // next index in args
	fs    *FlagSet // flag definitions & storage
	out   []string // collected positional args
	err   error    // first error encountered
}

// parseArgsWithFSM runs the FSM over args for a given FlagSet.
func parseArgsWithFSM(fs *FlagSet, args []string) ([]string, error) {
	p := &parser{fs: fs, args: args}
	if err := p.run(); err != nil {
		return nil, err
	}
	return p.out, nil
}

// run repeatedly invokes stateFns until completion or error.
func (p *parser) run() error {
	state := stateStartFn
	for state != nil {
		state = state(p)
		if p.err != nil {
			return p.err
		}
	}
	return nil
}

// next consumes and returns the next argument, or ok=false if done.
func (p *parser) next() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		p.index++
		ok = true
	}
	return
}

// peek returns the next argument without consuming it.
func (p *parser) peek() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		ok = true
	}
	return
}

// stateStartFn classifies the next token as long-flag, short-flag, or positional.
func stateStartFn(p *parser) stateFn {
	arg, ok := p.next()
	if !ok {
		return nil // done
	}
	switch {
	case arg == "--":
		// All following tokens are positional
		p.out = append(p.out, p.args[p.index:]...)
		p.index = len(p.args)
		return nil
	case strings.HasPrefix(arg, "--"):
		// long form: --flag or --flag=value
		return stateLongFlagFn(arg)
	case strings.HasPrefix(arg, "-") && len(arg) > 1:
		// short form: -f, -abc, or -fvalue
		return stateShortFlagFn(arg)
	default:
		// positional argument
		p.out = append(p.out, arg)
		return stateStartFn
	}
}

// stateLongFlagFn handles "--flag" or "--flag=value".
func stateLongFlagFn(arg string) stateFn {
	return func(p *parser) stateFn {
		// strip leading dashes and split on "=" if present
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

// stateShortFlagFn handles "-f", "-abc", "-fvalue", or "-f value".
func stateShortFlagFn(arg string) stateFn {
	return func(p *parser) stateFn {
		shorts := strings.TrimPrefix(arg, "-")
		for i := 0; i < len(shorts); i++ {
			ch := string(shorts[i])
			fl := findShortFlag(p.fs, ch)
			if fl == nil {
				p.err = fmt.Errorf("unknown short flag: -%s", ch)
				return nil
			}
			// boolean shorthand
			if tryBool(fl) {
				continue
			}
			// counter/increment
			if tryCounter(p, fl) {
				continue
			}
			// combined value: -fvalue
			if tryShortCombined(p, fl, i, shorts, ch) {
				break
			}
			// next token as value: -f value
			p.err = tryShortValue(p, fl, ch)
			break
		}
		return stateStartFn
	}
}

// handleDynamicFlag parses flags like "--group.id.field".
func handleDynamicFlag(name, rawVal string, hasVal bool) stateFn {
	return func(p *parser) stateFn {
		item, id := lookupDynamic(p, name)
		if p.err != nil {
			return nil
		}

		tryShort := func() bool {
			return tryDynamicBool(item, id) // only boolean shorthand; counters don’t apply
		}
		doSet := func(v string) error {
			return trySetDynamic(item, id, v, name)
		}

		handled, err := consumeValue(hasVal, rawVal, doSet, tryShort, p)
		if handled {
			p.err = err
			return stateStartFn
		}
		p.err = fmt.Errorf("missing value for flag: --%s", name)
		return nil
	}
}

// handleStaticFlag parses flags like "--flag".
// handleStaticFlag parses a static flag (--flag or --flag=value)
func handleStaticFlag(name, rawVal string, hasVal bool) stateFn {
	return func(p *parser) stateFn {
		fl := p.fs.flags[name]

		// tryShort now encompasses both bool‐shorthand and counters
		tryShort := func() bool {
			if tryBool(fl) {
				return true
			}
			if tryCounter(p, fl) {
				return true
			}
			return false
		}

		// doSet is the same form as dynamic: wrap Value.Set
		doSet := func(v string) error {
			return trySet(fl.Value, v, "invalid value for flag --%s: %s.", name)
		}

		handled, err := consumeValue(hasVal, rawVal, doSet, tryShort, p)
		if handled {
			p.err = err
			return stateStartFn
		}

		p.err = fmt.Errorf("missing value for flag: --%s", name)
		return nil
	}
}

// consumeValue handles in-order:
//   - explicit =value (hasVal==true)
//   - shorthand via tryShort()
//   - next-token as value
func consumeValue(
	hasVal bool,
	rawVal string,
	doSet func(string) error,
	tryShort func() bool,
	p *parser,
) (handled bool, err error) {
	if hasVal {
		return true, doSet(rawVal)
	}
	if tryShort() {
		return true, nil
	}
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return false, nil
	}
	p.next()
	return true, doSet(next)
}

// findShortFlag returns the BaseFlag matching a single-letter alias.
func findShortFlag(fs *FlagSet, short string) *core.BaseFlag {
	for _, fl := range fs.flags {
		if fl.Short == short {
			return fl
		}
	}
	return nil
}

// tryBool handles non-strict boolean shorthand for static flags.
func tryBool(fl *core.BaseFlag) bool {
	if sb, ok := fl.Value.(core.StrictBool); ok && !sb.IsStrictBool() {
		fl.Value.Set("true") // nolint:errcheck
		return true
	}
	return false
}

// lookupDynamic finds the DynamicValue for a flag named "group.id.field".
// On failure it sets p.err and returns nil.
func lookupDynamic(p *parser, name string) (item core.DynamicValue, id string) {
	parts := strings.Split(name, ".")
	if len(parts) != 3 {
		p.err = fmt.Errorf("invalid dynamic flag: --%s", name)
		return nil, ""
	}
	group, id, field := parts[0], parts[1], parts[2]

	groupMap, ok := p.fs.dynamic[group]
	if !ok {
		p.err = fmt.Errorf("unknown dynamic group: --%s", name)
		return nil, ""
	}
	item, ok = groupMap[field]
	if !ok {
		p.err = fmt.Errorf("unknown dynamic field: --%s", name)
		return nil, ""
	}
	return item, id
}

// tryDynamicBool handles non-strict boolean shorthand for dynamic flags.
func tryDynamicBool(item core.DynamicValue, id string) bool {
	if sb, ok := item.(core.StrictBool); ok && !sb.IsStrictBool() {
		item.Set(id, "true") // nolint:errcheck
		return true
	}
	return false
}

// tryCounter handles "counter" flags that increment on each occurrence.
func tryCounter(p *parser, fl *core.BaseFlag) bool {
	if inc, ok := fl.Value.(core.Incrementable); ok {
		p.err = inc.Increment()
		return true
	}
	return false
}

// tryShortCombined handles "-fvalue" in a grouped short-flag string.
func tryShortCombined(p *parser, fl *core.BaseFlag, idx int, shorts, char string) bool {
	if idx < len(shorts)-1 {
		p.err = trySet(fl.Value, shorts[idx+1:], "invalid value for flag -%s: %w", char)
		return true
	}
	return false
}

// tryShortValue handles "-f value" by peeking and consuming the next token.
func tryShortValue(p *parser, fl *core.BaseFlag, short string) error {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return fmt.Errorf("missing value for flag: -%s", short)
	}
	p.next()
	return trySet(fl.Value, next, "invalid value for flag -%s: %w", short)
}

// trySet wraps Value.Set with contextual error formatting.
func trySet(val core.Value, in, fmtStr, label string) error {
	if err := val.Set(in); err != nil {
		return fmt.Errorf(fmtStr, label, err)
	}
	return nil
}

// trySetDynamic wraps DynamicValue.Set with contextual error formatting.
func trySetDynamic(item core.DynamicValue, id, in, label string) error {
	if err := item.Set(id, in); err != nil {
		return fmt.Errorf("invalid value for dynamic flag --%s: %w", label, err)
	}
	return nil
}

// splitFlagArg splits "name=value" into (name, value, true),
// or ("name","",false) if no "=" is present.
func splitFlagArg(s string) (name, val string, hasVal bool) {
	if i := strings.Index(s, "="); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, "", false
}

// isDynamicFlag returns true for names like "group.id.field".
func isDynamicFlag(name string) bool {
	parts := strings.Split(name, ".")
	return len(parts) == 3
}

// isKnownStaticFlag returns true if name matches a registered static flag.
func isKnownStaticFlag(p *parser, name string) bool {
	_, ok := p.fs.flags[name]
	return ok
}
