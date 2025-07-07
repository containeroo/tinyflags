package tinyflags

import "fmt"

// FlagItem is a generic flag value holder for scalar types.
// It implements the Value interface and tracks the parsed value,
// default value, and whether the value has been changed.
type FlagItem[T any] struct {
	ptr       *T                      // Target pointer to store the parsed value
	def       T                       // Default value
	changed   bool                    // True if the value was set via CLI or environment
	parse     func(string) (T, error) // Function to parse a string into type T
	format    func(T) string          // Function to format a value of type T as string
	validator func(T) error           // optional validation callback
	allowed   []T                     // optional: for error reporting
}

// NewFlagItem constructs a new BaseValue for scalar flags.
// The parse and format functions must be provided to convert values to/from strings.
func NewFlagItem[T any](ptr *T, def T, parse func(string) (T, error), format func(T) string) *FlagItem[T] {
	*ptr = def
	bv := &FlagItem[T]{
		ptr:     ptr,
		def:     def,
		parse:   parse,
		format:  format,
		changed: false,
	}
	return bv
}

// Set parses and sets the value from the given string.
func (v *FlagItem[T]) Set(s string) error {
	val, err := v.parse(s)
	if err != nil {
		return err
	}
	if v.validator != nil {
		if err := v.validator(val); err != nil {
			return fmt.Errorf("invalid value %q: %w", s, err)
		}
	}
	*v.ptr = val
	v.changed = true
	return nil
}

// Get returns the current value stored in the pointer.
func (v *FlagItem[T]) Get() any {
	return *v.ptr
}

// Default returns the default value formatted as a string.
func (v *FlagItem[T]) Default() string {
	return v.format(v.def)
}

// IsChanged reports whether the value was explicitly set.
func (v *FlagItem[T]) IsChanged() bool {
	return v.changed
}

// SetValidator sets a validation callback for this flag.
func (v *FlagItem[T]) SetValidator(fn func(T) error, allowed []T) {
	v.validator = fn
	v.allowed = allowed
}
