package dynamic

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// SliceFlag provides per-ID slice flags.
// e.g., --http.alpha.tags=a,b or --node.node1.labels=env,prod.
type SliceFlag[T any] struct {
	*builder.DynamicFlag[T]
	item *DynamicSliceValue[T]
}

// Delimiter sets the separator between items.
func (f *SliceFlag[T]) Delimiter(sep string) *SliceFlag[T] { f.item.setDelimiter(sep); return f }

// Choices restricts allowed slice elements.
func (f *SliceFlag[T]) Choices(allowed ...T) *SliceFlag[T] {
	f.item.setValidate(func(v T) error {
		for _, a := range allowed {
			if f.item.format(a) == f.item.format(v) {
				return nil
			}
		}
		return fmt.Errorf("must be one of [%s]", utils.FormatAllowed(allowed, f.item.format))
	})

	formatted := make([]string, len(allowed))
	for i, a := range allowed {
		formatted[i] = f.item.format(a)
	}

	f.Allowed(formatted...)

	return f
}

// Validate lets you plug in arbitrary per‐element checks.
func (f *SliceFlag[T]) Validate(fn func(T) error) *SliceFlag[T] {
	f.item.setValidate(fn)
	return f
}

// Has returns true if the flag was set
func (f *SliceFlag[T]) Has(id string) bool {
	_, ok := f.item.values[id]
	return ok
}

// Get returns the parsed slice for the given ID.
// Fallback to default only happens here
func (f *SliceFlag[T]) Get(id string) ([]T, bool) {
	val, ok := f.item.values[id]
	if !ok {
		return f.item.def, false
	}
	return val, true
}

// Get returns the parsed slice for the given ID.
func (f *SliceFlag[T]) MustGet(id string) []T {
	val, ok := f.Get(id)
	if !ok {
		panic(fmt.Sprintf("required flag not set: %s for id %s", f.item.field, id))
	}
	return val
}

// Values returns all instance values.
func (f *SliceFlag[T]) Values() map[string][]T {
	return f.item.values
}

// ValuesAny returns values as a map[string]any for interface compatibility.
func (f *SliceFlag[T]) ValuesAny() map[string]any {
	out := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		out[k] = v
	}
	return out
}
