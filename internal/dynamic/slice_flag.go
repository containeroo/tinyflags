package dynamic

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

type SliceFlag[T any] struct {
	*builder.DynamicFlag[T]
	item *DynamicSliceValue[T]
}

func (f *SliceFlag[T]) Delimiter(sep string) *SliceFlag[T] {
	f.item.setDelimiter(sep)
	return f
}

func (f *SliceFlag[T]) Choices(allowed ...T) *SliceFlag[T] {
	f.item.setValidate(utils.AllowOnly(f.item.format, allowed))
	f.Allowed(utils.FormatList(f.item.format, allowed)...)
	return f
}

func (f *SliceFlag[T]) Validate(fn func(T) error) *SliceFlag[T] {
	f.item.setValidate(fn)
	return f
}

func (f *SliceFlag[T]) Has(id string) bool {
	_, ok := f.item.values[id]
	return ok
}

func (f *SliceFlag[T]) Get(id string) ([]T, bool) {
	val, ok := f.item.values[id]
	if !ok {
		return f.item.def, false
	}
	return val, true
}

func (f *SliceFlag[T]) MustGet(id string) []T {
	val, ok := f.Get(id)
	if !ok {
		panic(fmt.Sprintf("missing required value for %s (%s)", f.item.field, id))
	}
	return val
}

func (f *SliceFlag[T]) Values() map[string][]T {
	return f.item.values
}

func (f *SliceFlag[T]) ValuesAny() map[string]any {
	out := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		out[k] = v
	}
	return out
}
