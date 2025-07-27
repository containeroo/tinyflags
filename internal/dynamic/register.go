package dynamic

import (
	"strings"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

type dynamicValueProvider[T any] interface {
	core.Value
	Base() *DynamicScalarValue[T]
}

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
		Usage: usage,
	}

	g.items[field] = core.GroupItem{
		Value: item,
		Flag:  bf,
	}

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
		Value: &dynamicSliceValue[[]T]{def: strings.Join(formatted, ",")},
		Usage: usage,
	}

	g.items[field] = core.GroupItem{
		Value: item,
		Flag:  bf,
	}

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
	strict := new(bool)

	// Wrap base value
	val := &BoolValue{
		DynamicScalarValue: NewDynamicScalarValue(field, def, parse, format),
		Strict:             strict,
	}

	// Register BaseFlag
	bf := &core.BaseFlag{
		Name:  field,
		Usage: usage,
		Value: &dynamicBoolValue[bool]{def: format(def), strictMode: strict}, // dummy CLI value
	}

	// Store in dynamic group
	g.items[field] = core.GroupItem{
		Value: val,
		Flag:  bf,
	}

	// Return fluent builder
	return &BoolFlag{
		DynamicFlag: builder.NewDynamicFlag[bool](g.fs, bf),
		item:        val.Base(),
		strictMode:  strict,
	}
}
