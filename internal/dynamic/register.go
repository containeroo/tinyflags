package dynamic

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// registerDynamicScalar registers a scalar field under the group.
func registerDynamicScalar[T any](
	g *Group,
	field string,
	parse func(string) (T, error),
	format func(T) string,
) *ScalarFlag[T] {
	item := NewDynamicScalarValue(field, parse, format)
	bf := &core.BaseFlag{Name: field}
	g.items[field] = item

	return &ScalarFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        item,
	}
}

// registerDynamicSlice registers a slice field under the group.
func registerDynamicSlice[T any](
	g *Group,
	field string,
	parse func(string) (T, error),
	format func(T) string,
) *SliceFlag[T] {
	item := NewDynamicSliceValue(field, parse, format, g.fs.DefaultDelimiter())
	bf := &core.BaseFlag{Name: field}
	g.items[field] = item

	return &SliceFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        item,
	}
}
