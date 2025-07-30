package slice

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// SliceFlag is the user‐facing builder for slice flags.
type SliceFlag[T any] struct {
	builder.StaticFlag[[]T]
	val *SliceValue[T]
}

// Delimiter sets the delimiter used to split input values.
func (f *SliceFlag[T]) Delimiter(sep string) *SliceFlag[T] {
	f.val.delimiter = sep
	return f
}

// Choices restricts allowed slice elements.
func (f *SliceFlag[T]) Choices(allowed ...T) *SliceFlag[T] {
	f.val.setFinalize(utils.AllowOnly(f.val.format, allowed))
	f.Allowed(utils.FormatList(f.val.format, allowed)...)

	return f
}

// Finalize lets you plug in arbitrary per‐element checks.
func (f *SliceFlag[T]) Finalize(fn func(T) error) *SliceFlag[T] {
	f.val.setFinalize(fn)
	return f
}
