package argparse

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// Config supplies the callbacks and behavior needed by the argument parser.
type Config struct {
	ContinueOnError   bool
	LookupStaticFlag  func(string) *core.BaseFlag
	LookupShortFlag   func(string) *core.BaseFlag
	LookupDynamicFlag func(string, string) (core.DynamicValue, string, error)
	HandleUnknownFlag func(string) error
}

type stateFn func(*parser) stateFn

type parser struct {
	config Config
	args   []string
	index  int
	out    []string
	err    error
	errs   []error
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

func (p *parser) next() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		p.index++
		ok = true
	}
	return arg, ok
}

func (p *parser) peek() (arg string, ok bool) {
	if p.index < len(p.args) {
		arg = p.args[p.index]
		ok = true
	}
	return arg, ok
}

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

func handleDynamic(name, val string, hasVal bool, raw string) stateFn {
	return func(p *parser) stateFn {
		item, id, err := p.config.LookupDynamicFlag(name, raw)
		if err != nil {
			p.err = err
			return nil
		}

		if handled := tryDynamicBool(item, id); handled {
			return stateStart
		}

		item.GetAny(name)

		if hasVal {
			p.err = trySetDynamic(item, id, val, name)
			return stateStart
		}

		if handled := handleDynamicValue(p, item, id, name); !handled {
			return nil
		}

		return stateStart
	}
}

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

func handleStatic(name, val string, hasVal bool) stateFn {
	return func(p *parser) stateFn {
		flag := p.config.LookupStaticFlag(name)

		if handled := tryBool(flag); handled {
			return stateStart
		}
		if handled := tryCounter(p, flag); handled {
			return stateStart
		}
		if hasVal {
			p.err = trySet(flag.Value, val, "invalid value for flag --%s: %w", name)
			return stateStart
		}
		if handled := tryLongValue(p, flag, name); handled {
			return stateStart
		}

		p.err = fmt.Errorf("missing value for flag --%s", name)
		return nil
	}
}

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
				continue
			}
			if handled := tryCounter(p, flag); handled {
				continue
			}
			if handled := tryShortCombined(p, flag, i, shorts, char); handled {
				break
			}

			p.err = tryShortValue(p, flag, char)
			break
		}

		return stateStart
	}
}

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

func tryDynamicBool(item core.DynamicValue, id string) bool {
	if b, ok := item.(core.StrictBool); ok && !b.IsStrictBool() {
		item.Set(id, "true") // nolint:errcheck
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
	p.err = trySet(flag.Value, next, "invalid value for flag --%s: %w", name)
	return true
}

func tryShortValue(p *parser, flag *core.BaseFlag, short string) error {
	next, ok := p.peek()
	if !ok || strings.HasPrefix(next, "-") {
		return fmt.Errorf("missing value for flag -%s", flag.Short)
	}
	p.next()
	return trySet(flag.Value, next, "invalid value for flag -%s: %w", short)
}

func trySet(value core.Value, input string, format string, label string) error {
	if err := value.Set(input); err != nil {
		return fmt.Errorf(format, label, err)
	}
	return nil
}

func trySetDynamic(item core.DynamicValue, id, val, label string) error {
	if err := item.Set(id, val); err != nil {
		return fmt.Errorf("invalid value for flag --%s: %w", label, err)
	}
	return nil
}

func splitFlagArg(s string) (name, val string, hasVal bool) {
	if i := strings.Index(s, "="); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, "", false
}

func isDynamicFlag(name string) bool {
	return len(strings.Split(name, ".")) == 3
}
