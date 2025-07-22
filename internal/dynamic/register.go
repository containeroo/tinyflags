package dynamic

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// defineDynamicScalar registers a scalar field under the group.
func defineDynamicScalar[T any](
	g *Group,
	field string,
	parse func(string) (T, error),
	format func(T) string,
) *ScalarFlag[T] {
	item := NewDynamicScalarValue(field, parse, format)
	if err := g.fs.RegisterDynamic(g.prefix, field, item); err != nil {
		panic(err)
	}
	bf := &core.BaseFlag{Name: field}

	return &ScalarFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        item,
	}
}

// defineDynamicSlice registers a slice field under the group.
func defineDynamicSlice[T any](
	g *Group,
	field string,
	parse func(string) (T, error),
	format func(T) string,
) *SliceFlag[T] {
	item := NewDynamicSliceValue(field, parse, format, g.fs.DefaultDelimiter())
	if err := g.fs.RegisterDynamic(g.prefix, field, item); err != nil {
		panic(err)
	}
	bf := &core.BaseFlag{Name: field}

	return &SliceFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        item,
	}
}
