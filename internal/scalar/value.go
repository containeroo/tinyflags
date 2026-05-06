package scalar

import (
	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/utils"
)

// ScalarValue implements scalar flag parsing, formatting, and validation.
type ScalarValue[T any] struct {
	ptr      *T
	def      T
	changed  bool
	hooks    core.ValueHooks[T]
}

// NewScalarValue creates a new scalar value.
func NewScalarValue[T any](ptr *T, def T, parse func(string) (T, error), format func(T) string) *ScalarValue[T] {
	*ptr = def
	return &ScalarValue[T]{
		ptr:   ptr,
		def:   def,
		hooks: core.NewValueHooks(parse, format),
	}
}

// Set parses and stores one entry.
func (f *ScalarValue[T]) Set(s string) error {
	val, err := f.hooks.ParseValue(s)
	if err != nil {
		return err
	}
	*f.ptr = val
	f.changed = true
	return nil
}

// Get returns the stored value.
func (f *ScalarValue[T]) Get() any { return *f.ptr }

// Default returns the default value as string.
func (f *ScalarValue[T]) Default() string { return f.hooks.DefaultString(f.def) }

// Changed returns true if the value was changed.
func (f *ScalarValue[T]) Changed() bool { return f.changed }

// setValidate sets a per-item validation function.
func (f *ScalarValue[T]) setValidate(fn func(T) error) { f.hooks.SetValidate(fn) }

// setFinalize sets a per-item finalizer function.
func (f *ScalarValue[T]) setFinalize(fn func(T) T) { f.hooks.SetFinalize(fn) }

// setFinalizeDefaultValue enables running the finalizer on defaults when unset.
func (f *ScalarValue[T]) setFinalizeDefaultValue() { f.hooks.EnableFinalizeDefault() }

// Base returns the underlying value.
func (f *ScalarValue[T]) Base() *ScalarValue[T] { return f }

// ApplyDefaultFinalize applies the default-only finalizer when unset.
func (f *ScalarValue[T]) ApplyDefaultFinalize() {
	f.hooks.ApplyDefaultScalar(f.ptr, f.changed)
}

// ResetParseState restores the default value and clears changed/finalized state.
func (f *ScalarValue[T]) ResetParseState() {
	utils.ResetScalarState(f.ptr, f.def, &f.changed, &f.hooks.DefaultFinalized)
}
