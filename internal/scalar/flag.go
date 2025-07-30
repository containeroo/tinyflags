package scalar

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// ScalarFlag is the user-facing scalar flag builder.
type ScalarFlag[T any] struct {
	builder.StaticFlag[T]
	val *ScalarValue[T]
}

// Choices restricts allowed scalar values.
func (f *ScalarFlag[T]) Choices(allowed ...T) *ScalarFlag[T] {
	f.val.setValidate(utils.AllowOnly(f.val.format, allowed))
	f.Allowed(utils.FormatList(f.val.format, allowed)...)
	return f
}

// Validate lets you plug in arbitrary per‐element checks.
func (f *ScalarFlag[T]) Validate(fn func(T) error) *ScalarFlag[T] {
	f.val.setValidate(fn)
	return f
}

// Finalize sets a custom finalizer function for each
func (f *ScalarFlag[T]) Finalize(fn func(T) T) *ScalarFlag[T] {
	f.val.setFinalize(fn)
	return f
}
