package tinyflags

import (
	"fmt"
	"strings"
)

// stateFn defines a function representing a parser state.
type stateFn func(*parser) stateFn

// parser holds parsing state and context for command-line argument parsing.
type parser struct {
	args  []string
	index int
	fs    *FlagSet
	out   []string
	err   error
}

// next returns the next argument and advances the index.
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

// run executes the parser state machine.
func (p *parser) run() error {
	state := stateStartFn
	for state != nil {
		state = state(p)
	}
	return p.err
}

func parseArgsWithFSM(fs *FlagSet, args []string) ([]string, error) {
	p := &parser{fs: fs, args: args}
	err := p.run()
	return p.out, err
}

func stateStartFn(p *parser) stateFn {
	arg, ok := p.next()
	if !ok {
		return nil
	}

	if arg == "--" {
		p.out = append(p.out, p.args[p.index:]...)
		p.index = len(p.args)
		return nil
	}

	if strings.HasPrefix(arg, "--") {
		return stateLongFlagFn(arg)
	}

	if strings.HasPrefix(arg, "-") && len(arg) > 1 {
		return stateShortFlagFn(arg)
	}

	p.out = append(p.out, arg)
	return stateStartFn
}

func stateLongFlagFn(arg string) stateFn {
	return func(p *parser) stateFn {
		nameval := strings.TrimPrefix(arg, "--")
		name, val, hasVal := splitFlagArg(nameval)

		// Dynamic: --group.id.flag
		// Dynamic: --group.id.field
		if dot := strings.Index(name, "."); dot != -1 {
			groupName := name[:dot]
			rest := name[dot+1:]

			// dynamic format: id.field
			if dot2 := strings.Index(rest, "."); dot2 != -1 {
				id := rest[:dot2]
				field := rest[dot2+1:]

				groupFields, ok := p.fs.dynamic[groupName]
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

		// Static
		fl := p.fs.flags[name]
		if fl == nil {
			p.err = fmt.Errorf("unknown flag: --%s", name)
			return nil
		}

		// Non-strict bool
		if b, ok := fl.value.(*BoolValue); ok && !b.IsStrictBool() {
			fl.value.Set("true") // nolint:errcheck
			return stateStartFn
		}

		if hasVal {
			p.err = trySet(fl.value, val, "invalid value for flag --%s: %w", name)
			return stateStartFn
		}

		next, ok := p.peek()
		if !ok || strings.HasPrefix(next, "-") {
			p.err = fmt.Errorf("missing value for flag: --%s", name)
			return nil
		}
		p.next()
		p.err = trySet(fl.value, next, "invalid value for flag --%s: %w", name)
		return stateStartFn
	}
}

func stateShortFlagFn(arg string) stateFn {
	return func(p *parser) stateFn {
		shorts := strings.TrimPrefix(arg, "-")
		for i := 0; i < len(shorts); i++ {
			char := string(shorts[i])
			var target *baseFlag
			for _, fl := range p.fs.flags {
				if fl.short == char {
					target = fl
					break
				}
			}
			if target == nil {
				p.err = fmt.Errorf("unknown short flag: -%s", char)
				return nil
			}

			if b, ok := target.value.(*BoolValue); ok && !b.IsStrictBool() {
				target.value.Set("true") // nolint:errcheck
				continue
			}

			if i < len(shorts)-1 {
				val := shorts[i+1:]
				p.err = trySet(target.value, val, "invalid value for flag -%s: %w", char)
				break
			}

			next, ok := p.peek()
			if !ok || strings.HasPrefix(next, "-") {
				p.err = fmt.Errorf("missing value for flag: -%s", char)
				return nil
			}
			p.next()
			p.err = trySet(target.value, next, "invalid value for flag -%s: %w", char)
			break
		}
		return stateStartFn
	}
}

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

func splitFlagArg(s string) (name, val string, hasVal bool) {
	if i := strings.Index(s, "="); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, "", false
}
