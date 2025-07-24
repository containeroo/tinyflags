package dynamic

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// registerDynamicScalar registers a scalar field under the group.
func registerDynamicScalar[T any](
	g *Group,
	field string,
	def T,
	parse func(string) (T, error),
	format func(T) string,
) *ScalarFlag[T] {
	item := NewDynamicScalarValue(field, def, parse, format)
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
	def []T,
	parse func(string) (T, error),
	format func(T) string,
) *SliceFlag[T] {
	item := NewDynamicSliceValue(field, def, parse, format, g.fs.DefaultDelimiter())
	bf := &core.BaseFlag{Name: field}
	g.items[field] = item

	return &SliceFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        item,
	}
}

// registerDynamicBool registers a dynamic boolean field under the group.
func registerDynamicBool(
	g *Group,
	field string,
	def bool,
	parse func(string) (bool, error),
	format func(bool) string,
) *BoolFlag {
	item := NewDynamicScalarValue(field, def, parse, format)

	// BoolValue wraps item to expose IsStrictBool
	flagVal := &BoolValue{
		item:       item,
		strictMode: new(bool), // pointer for later mutation by .Strict()
	}

	g.items[field] = flagVal // Store in dynamic group registry

	// Also create and register BaseFlag
	bf := &core.BaseFlag{Name: field}
	g.fs.RegisterFlag(field, bf)

	// Return the user-facing BoolFlag
	return &BoolFlag{
		DynamicFlag: builder.NewDynamicFlag[bool](g.fs, bf),
		item:        item,
		strictMode:  *flagVal.strictMode,
	}
}
