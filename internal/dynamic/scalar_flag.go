package dynamic

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// ScalarFlag provides per-ID scalar flags.
// For example: --http.alpha.port=8080
type ScalarFlag[T any] struct {
	*builder.DynamicFlag[T]
	item *DynamicScalarValue[T]
}

// Choices restricts allowed scalar values.
func (f *ScalarFlag[T]) Choices(allowed ...T) *ScalarFlag[T] {
	f.item.setValidate(func(v T) error {
		for _, a := range allowed {
			if f.item.format(a) == f.item.format(v) {
				return nil
			}
		}
		return fmt.Errorf("must be one of %s", utils.FormatAllowed(allowed, f.item.format))
	})

	formatted := make([]string, len(allowed))
	for i, a := range allowed {
		formatted[i] = f.item.format(a)
	}

	f.Allowed(formatted...)

	return f
}

// Validate lets you plug in arbitrary per‐element checks.
func (f *ScalarFlag[T]) Validate(fn func(T) error) *ScalarFlag[T] {
	f.item.setValidate(fn)
	return f
}

// Has returns true if the flag was set
func (f *ScalarFlag[T]) Has(id string) bool {
	_, ok := f.item.values[id]
	return ok
}

// Fallback to default only happens here
func (f *ScalarFlag[T]) Get(id string) (T, bool) {
	val, ok := f.item.values[id]
	if !ok {
		return f.item.def, false
	}
	return val, true
}

func (f *ScalarFlag[T]) MustGet(id string) T {
	val, ok := f.Get(id)
	if !ok {
		panic(fmt.Sprintf("required flag not set: %s for id %s", f.item.field, id))
	}
	return val
}

// Values returns all stored values.
func (f *ScalarFlag[T]) Values() map[string]T {
	return f.item.values
}

// ValuesAny returns values as a generic map.
func (f *ScalarFlag[T]) ValuesAny() map[string]any {
	m := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		m[k] = v
	}
	return m
}
