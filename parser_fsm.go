package tinyflags

import (
	"errors"
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

		// Look up the flag by name
		fl := p.fs.flags[name]
		if fl == nil {
			p.err = errors.New("unknown flag: --" + name)
			return nil
		}

		// Handle booleans without explicit values: --flag ⇒ true
		if f, ok := fl.value.(*BoolValue); ok && !f.IsStrictBool() {
			fl.value.Set("true") // nolint:errcheck
			return stateStartFn
		}

		// If flag is of the form --flag=value
		if hasVal {
			if err := fl.value.Set(val); err != nil {
				p.err = fmt.Errorf("invalid value for flag --%s: %w", name, err)
			}
			return stateStartFn
		}

		// Otherwise, consume the next arg as the value
		next, ok := p.peek()
		if !ok || strings.HasPrefix(next, "-") {
			p.err = errors.New("missing value for flag: --" + name)
			return nil
		}

		p.next() // consume the next argument
		if err := fl.value.Set(next); err != nil {
			p.err = fmt.Errorf("invalid value for flag --%s: %w", name, err)
		}
		return stateStartFn
	}
}

// stateShortFlagFn handles grouped short flags like -abc or single ones like -f value.
func stateShortFlagFn(arg string) stateFn {
	return func(p *parser) stateFn {
		shorts := strings.TrimPrefix(arg, "-")

		for i := 0; i < len(shorts); i++ {
			char := string(shorts[i])

			// Find the flag with matching short name
			var target *baseFlag
			for _, fl := range p.fs.flags {
				if fl.short == char {
					target = fl
					break
				}
			}
			if target == nil {
				p.err = errors.New("unknown short flag: -" + char)
				return nil
			}

			// Handle -f where f is a non-strict bool flag
			if b, ok := target.value.(*BoolValue); ok && !b.IsStrictBool() {
				target.value.Set("true") // nolint:errcheck
				continue                 // check next short flag in group
			}

			// Handle -p8080 (combined value)
			if i < len(shorts)-1 {
				val := shorts[i+1:]
				if err := target.value.Set(val); err != nil {
					p.err = fmt.Errorf("invalid value for flag -%s: %w", char, err)
				}
				break // done with this flag group
			}

			// Handle -p 8080 (value in next argument)
			next, ok := p.peek()
			if !ok || strings.HasPrefix(next, "-") {
				p.err = errors.New("missing value for flag: -" + char)
				return nil
			}
			p.next() // consume next arg
			if err := target.value.Set(next); err != nil {
				p.err = fmt.Errorf("invalid value for flag -%s: %w", char, err)
			}
			break
		}
		return stateStartFn
	}
}

// splitFlagArg splits a string like "flag=value" into ("flag", "value", true).
// If there's no '=', it returns ("flag", "", false).
func splitFlagArg(s string) (name, val string, hasVal bool) {
	if i := strings.Index(s, "="); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, "", false
}
