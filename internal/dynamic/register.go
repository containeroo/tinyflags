package dynamic

import (
	"strings"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// registerDynamicScalar registers a scalar field under the group.
func registerDynamicScalar[T any](
	g *Group,
	field string,
	def T,
	usage string,
	parse func(string) (T, error),
	format func(T) string,
) *ScalarFlag[T] {
	item := NewDynamicScalarValue(field, def, parse, format)

	bf := &core.BaseFlag{
		Name:  field,
		Value: &dynamicHelpValue[T]{def: format(def)},
		Usage: usage,
	}
	g.items[field] = item
	g.fs.RegisterFlag(field, bf)

	return &ScalarFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        item,
	}
}

// registerDynamicSlice registers a slice field under the group.
func registerDynamicSlice[T any](
	g *Group,
	field string,
	def []T,
	usage string,
	parse func(string) (T, error),
	format func(T) string,
) *SliceFlag[T] {
	item := NewDynamicSliceValue(field, def, parse, format, g.fs.DefaultDelimiter())

	formatted := make([]string, len(def))
	for i, val := range def {
		formatted[i] = format(val)
	}

	bf := &core.BaseFlag{
		Name:  field,
		Value: &dynamicHelpValue[[]T]{def: strings.Join(formatted, ",")},
		Usage: usage,
	}

	g.items[field] = item
	g.fs.RegisterFlag(field, bf)

	return &SliceFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        item,
	}
}

// registerDynamicBool registers a boolean field under the group.
func registerDynamicBool(
	g *Group,
	field string,
	def bool,
	usage string,
	parse func(string) (bool, error),
	format func(bool) string,
) *BoolFlag {
	item := NewDynamicScalarValue(field, def, parse, format)

	bf := &core.BaseFlag{
		Name:  field,
		Usage: usage,
	}
	g.items[field] = &BoolValue{DynamicScalarValue: item, strictMode: false}
	g.fs.RegisterFlag(field, bf)

	return &BoolFlag{
		DynamicFlag: builder.NewDynamicFlag[bool](g.fs, bf),
		item:        item,
	}
}
