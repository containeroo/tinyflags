package scalar

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// scalarFlagBase provides shared builder logic while preserving the concrete flag type.
type scalarFlagBase[T any, Self any] struct {
	builder.StaticFlag[T, Self]
	val  *ScalarValue[T]
	self Self
}

// Choices restricts allowed scalar values.
func (f *scalarFlagBase[T, Self]) Choices(allowed ...T) Self {
	f.val.setValidate(utils.AllowOnly(f.val.format, allowed))
	f.Allowed(utils.FormatList(f.val.format, allowed)...)
	return f.self
}

// Validate lets you plug in arbitrary per‐element checks.
func (f *scalarFlagBase[T, Self]) Validate(fn func(T) error) Self {
	f.val.setValidate(fn)
	return f.self
}

// Finalize sets a custom finalizer function for each value.
func (f *scalarFlagBase[T, Self]) Finalize(fn func(T) T) Self {
	f.val.setFinalize(fn)
	return f.self
}

// FinalizeDefaultValue runs the finalizer for defaults when the flag is unset.
func (f *scalarFlagBase[T, Self]) FinalizeDefaultValue() Self {
	f.val.setFinalizeDefaultValue()
	return f.self
}

// Default returns the default value.
func (f *scalarFlagBase[T, Self]) Default() T {
	return f.val.def
}

// Changed returns true if the value was changed.
func (f *scalarFlagBase[T, Self]) Changed() bool {
	return f.val.changed
}

// ScalarFlag is the user-facing scalar flag builder.
type ScalarFlag[T any] struct {
	scalarFlagBase[T, *ScalarFlag[T]]
}
