package slice

import (
	"fmt"
	"strings"
)

// SliceValue implements slice flag parsing and validation.
type SliceValue[T any] struct {
	ptr       *[]T
	def       []T
	value     []T
	changed   bool
	delimiter string
	parse     func(string) (T, error)
	format    func(T) string
	validate  func(T) error
}

// NewSliceValue creates a new slice value.
func NewSliceValue[T any](
	ptr *[]T,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *SliceValue[T] {
	*ptr = append([]T{}, def...)
	return &SliceValue[T]{
		ptr:       ptr,
		def:       def,
		delimiter: delimiter,
		parse:     parse,
		format:    format,
	}
}

// Set parses and stores the slice from a delimited string for a given ID.
func (v *SliceValue[T]) Set(s string) error {
	if !v.changed {
		*v.ptr = nil
	}
	parts := strings.Split(s, v.delimiter)
	for _, raw := range parts {
		val, err := v.parse(strings.TrimSpace(raw))
		if err != nil {
			return fmt.Errorf("invalid slice item %q: %w", raw, err)
		}
		if v.validate != nil {
			if err := v.validate(val); err != nil {
				return fmt.Errorf("invalid value %q: %w", raw, err)
			}
		}
		*v.ptr = append(*v.ptr, val)
	}
	v.value = *v.ptr
	v.changed = true
	return nil
}

// Get returns the parsed slice for the given ID.
func (v *SliceValue[T]) Get() any {
	return *v.ptr
}

// Default returns the default value as string.
func (f *SliceValue[T]) Default() string {
	out := make([]string, 0, len(f.def))
	for _, v := range f.def {
		out = append(out, f.format(v))
	}
	return strings.Join(out, f.delimiter)
}

// Changed returns true if the value was changed.
func (v *SliceValue[T]) Changed() bool {
	return v.changed
}

// setFinalize sets a per-item validation function.
func (v *SliceValue[T]) setFinalize(fn func(T) error) { v.validate = fn }

// Base returns the underlying value.
func (v *SliceValue[T]) Base() *SliceValue[T] { return v }

// isSlice is a no-op marker method to implement core.SliceMarker.
func (v *SliceValue[T]) IsSlice() {}
