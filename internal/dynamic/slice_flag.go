package dynamic

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// SliceFlag represents a dynamic slice flag with per-ID values.
type SliceFlag[T any] struct {
	*builder.DynamicFlag[T]                       // Embedded flag metadata
	item                    *DynamicSliceValue[T] // Underlying slice value store
}

// Delimiter sets the string delimiter for parsing slice values.
func (f *SliceFlag[T]) Delimiter(sep string) *SliceFlag[T] {
	f.item.setDelimiter(sep)
	return f
}

// Choices restricts allowed slice values to the given list.
func (f *SliceFlag[T]) Choices(allowed ...T) *SliceFlag[T] {
	f.item.setValidate(utils.AllowOnly(f.item.format, allowed))
	f.Allowed(utils.FormatList(f.item.format, allowed)...)
	return f
}

// Validate sets a custom validation function for each element.
func (f *SliceFlag[T]) Validate(fn func(T) error) *SliceFlag[T] {
	f.item.setValidate(fn)
	return f
}

// Finalize sets a custom finalizer function for each element.
func (f *SliceFlag[T]) Finalize(fn func(T) T) *SliceFlag[T] {
	f.item.setFinalize(fn)
	return f
}

// FinalizeWithID sets a finalizer that receives the instance ID.
func (f *SliceFlag[T]) FinalizeWithID(fn func(string, T) T) *SliceFlag[T] {
	f.item.setFinalizeWithID(fn)
	return f
}

// Default returns the default value.
func (f *SliceFlag[T]) Default() []T {
	return f.item.def
}

// Changed returns true if the value was changed.
func (f *SliceFlag[T]) Changed() bool {
	return f.item.changed
}

// Has reports whether a value is set for the given ID.
func (f *SliceFlag[T]) Has(id string) bool {
	_, ok := f.item.values[id]
	return ok
}

// Get returns the slice for a given ID and whether it exists.
func (f *SliceFlag[T]) Get(id string) ([]T, bool) {
	val, ok := f.item.values[id]
	if !ok {
		return f.item.def, false
	}
	return val, true
}

// MustGet returns the value or panics if it is not set.
func (f *SliceFlag[T]) MustGet(id string) []T {
	val, ok := f.Get(id)
	if !ok {
		panic(fmt.Sprintf("missing required value for %s (%s)", f.item.field, id))
	}
	return val
}

// Values returns all parsed values keyed by ID.
func (f *SliceFlag[T]) Values() map[string][]T {
	return f.item.values
}

// ValuesAny returns all values as a map of any.
func (f *SliceFlag[T]) ValuesAny() map[string]any {
	out := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		out[k] = v
	}
	return out
}

// StrictDelimiter rejects inputs that mix different separators.
func (f *SliceFlag[T]) StrictDelimiter() *SliceFlag[T] {
	f.item.setStrictDelimiter(true)
	return f
}

// AllowEmpty permits empty items in the slice (e.g., "a,,b").
func (f *SliceFlag[T]) AllowEmpty() *SliceFlag[T] {
	f.item.setAllowEmpty(true)
	return f
}
