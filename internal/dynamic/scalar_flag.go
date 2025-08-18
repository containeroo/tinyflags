package dynamic

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// ScalarFlag represents a typed dynamic scalar flag with per-ID values.
type ScalarFlag[T any] struct {
	*builder.DynamicFlag[T]                        // Embedded base flag metadata
	item                    *DynamicScalarValue[T] // Underlying value and parser
}

// Choices restricts the allowed values to the provided list.
func (f *ScalarFlag[T]) Choices(allowed ...T) *ScalarFlag[T] {
	f.item.setValidate(utils.AllowOnly(f.item.format, allowed))
	f.Allowed(utils.FormatList(f.item.format, allowed)...)
	return f
}

// Validate adds a custom validation function for values.
func (f *ScalarFlag[T]) Validate(fn func(T) error) *ScalarFlag[T] {
	f.item.setValidate(fn)
	return f
}

// Finalize adds a custom finalizer function for values.
func (f *ScalarFlag[T]) Finalize(fn func(T) T) *ScalarFlag[T] {
	f.item.setFinalize(fn)
	return f
}

// Default returns the default value.
func (f *ScalarFlag[T]) Default() T {
	return f.item.def
}

// Changed returns true if the value was changed.
func (f *ScalarFlag[T]) Changed() bool {
	return f.item.changed
}

// Has reports whether a value was set for the given ID.
func (f *ScalarFlag[T]) Has(id string) bool {
	_, ok := f.item.values[id]
	return ok
}

// Get returns the value for the given ID and whether it exists.
func (f *ScalarFlag[T]) Get(id string) (T, bool) {
	val, ok := f.item.values[id]
	if !ok {
		return f.item.def, false
	}
	return val, true
}

// MustGet returns the value or panics if missing.
func (f *ScalarFlag[T]) MustGet(id string) T {
	val, ok := f.Get(id)
	if !ok {
		panic(fmt.Sprintf("missing required value for %s (%s)", f.item.field, id))
	}
	return val
}

// Values returns all parsed values keyed by ID.
func (f *ScalarFlag[T]) Values() map[string]T {
	return f.item.values
}

// ValuesAny returns all values as a map of any.
func (f *ScalarFlag[T]) ValuesAny() map[string]any {
	out := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		out[k] = v
	}
	return out
}
