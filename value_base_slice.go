package tinyflags

import (
	"fmt"
	"strings"
)

// SliceFlagItem is a generic flag value holder for slice types.
// It implements the Value interface for flags that accept multiple values,
// optionally split by a custom delimiter.
type SliceFlagItem[T any] struct {
	ptr       *[]T                    // Target pointer to store the parsed slice
	def       []T                     // Default slice value
	changed   bool                    // True if the value was set via CLI or environment
	parse     func(string) (T, error) // Function to parse a single element
	format    func(T) string          // Function to format a single element
	delimiter string                  // Delimiter used to split input string into slice elements
	validator func(T) bool            // optional validation callback
	allowed   []T                     // optional: for help/choices output
}

// NewSliceItem constructs a new BaseSliceValue.
// The default value is defensively copied into the target pointer.
func NewSliceItem[T any](
	ptr *[]T,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *SliceFlagItem[T] {
	*ptr = append([]T{}, def...) // defensive copy to avoid modifying input
	return &SliceFlagItem[T]{
		ptr:       ptr,
		def:       def,
		parse:     parse,
		format:    format,
		delimiter: delimiter,
	}
}

// isSlice implements the SliceFlag marker interface for BaseSliceValue.
func (v *SliceFlagItem[T]) isSlice() {}

// Set splits the input string by the delimiter and parses each item.
func (v *SliceFlagItem[T]) Set(s string) error {
	if !v.changed {
		*v.ptr = nil // Clear default only on first Set
	}
	for _, p := range strings.Split(s, v.delimiter) {
		val, err := v.parse(strings.TrimSpace(p))
		if err != nil {
			return fmt.Errorf("invalid slice item %q: %w", p, err)
		}
		if v.validator != nil && !v.validator(val) {
			return fmt.Errorf("invalid value %q: must be one of [%s]", s, formatAllowed(v.allowed, v.format))
		}
		*v.ptr = append(*v.ptr, val)
	}
	v.changed = true
	return nil
}

// Get returns the current slice value.
func (v *SliceFlagItem[T]) Get() any {
	return *v.ptr
}

// Default returns the default slice as a single string, joined by the delimiter.
func (v *SliceFlagItem[T]) Default() string {
	formatted := make([]string, 0, len(v.def))
	for _, item := range v.def {
		formatted = append(formatted, v.format(item))
	}
	return strings.Join(formatted, v.delimiter)
}

// IsChanged reports whether the slice value was explicitly set.
func (v *SliceFlagItem[T]) IsChanged() bool {
	return v.changed
}

// SetValidator sets a validation callback for this flag.
func (v *SliceFlagItem[T]) SetValidator(fn func(T) bool, allowed []T) {
	v.validator = fn
	v.allowed = allowed
}

// SetDelimiter sets the delimiter used for splitting input strings into slice elements.
func (v *SliceFlagItem[T]) SetDelimiter(d string) {
	v.delimiter = d
}
