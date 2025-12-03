package slice

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// SliceFlag is the user‐facing builder for slice flags.
type SliceFlag[T any] struct {
	builder.StaticFlag[[]T, *SliceFlag[T]]
	val *SliceValue[T]
}

// Delimiter sets the delimiter used to split input values.
func (f *SliceFlag[T]) Delimiter(sep string) *SliceFlag[T] {
	f.val.delimiter = sep
	return f
}

// Choices restricts allowed slice elements.
func (f *SliceFlag[T]) Choices(allowed ...T) *SliceFlag[T] {
	f.val.setValidate(utils.AllowOnly(f.val.format, allowed))
	f.Allowed(utils.FormatList(f.val.format, allowed)...)

	return f
}

// Validate lets you plug in arbitrary per‐element checks.
func (f *SliceFlag[T]) Validate(fn func(T) error) *SliceFlag[T] {
	f.val.setValidate(fn)
	return f
}

// Finalize sets a custom finalizer function for each element.
func (f *SliceFlag[T]) Finalize(fn func(T) T) *SliceFlag[T] {
	f.val.setFinalize(fn)
	return f
}

// StrictDelimiter rejects inputs that mix different separators.
func (f *SliceFlag[T]) StrictDelimiter() *SliceFlag[T] {
	f.val.setStrictDelimiter(true)
	return f
}

// AllowEmpty permits empty items in the slice (e.g., "a,,b").
func (f *SliceFlag[T]) AllowEmpty() *SliceFlag[T] {
	f.val.setAllowEmpty(true)
	return f
}

// Default returns the default value.
func (f *SliceFlag[T]) Default() []T {
	return f.val.def
}

// Changed returns true if the value was changed.
func (f *SliceFlag[T]) Changed() bool {
	return f.val.changed
}
