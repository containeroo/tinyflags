package dynamic

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

type ScalarFlag[T any] struct {
	*builder.DynamicFlag[T]
	item *DynamicScalarValue[T]
}

func (f *ScalarFlag[T]) Choices(allowed ...T) *ScalarFlag[T] {
	f.item.setValidate(utils.AllowOnly(f.item.format, allowed))
	f.Allowed(utils.FormatList(f.item.format, allowed)...)
	return f
}

func (f *ScalarFlag[T]) Validate(fn func(T) error) *ScalarFlag[T] {
	f.item.setValidate(fn)
	return f
}

func (f *ScalarFlag[T]) Has(id string) bool {
	_, ok := f.item.values[id]
	return ok
}

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
		panic(fmt.Sprintf("missing required value for %s (%s)", f.item.field, id))
	}
	return val
}

func (f *ScalarFlag[T]) Values() map[string]T {
	return f.item.values
}

func (f *ScalarFlag[T]) ValuesAny() map[string]any {
	out := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		out[k] = v
	}
	return out
}
