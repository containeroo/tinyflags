package dynamic

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/utils"
)

// registerDynamicSlice registers a slice field under the group.
func registerDynamicSlice[T any](
	g *Group,
	field string,
	def []T,
	usage string,
	parse func(string) (T, error),
	format func(T) string,
	trimSpace bool,
) *SliceFlag[T] {
	return registerSliceValue(g, field, def, usage, format,
		NewDynamicSliceValue(field, def, parse, format, g.fs.DefaultDelimiter(), trimSpace))
}

// RegisterSlice registers a custom parser that expands each input chunk into elements.
func RegisterSlice[T any](g *Group, field string, def []T, usage string,
	parse func(string) ([]T, error), format func(T) string,
) *SliceFlag[T] {
	val := NewDynamicSliceValue(field, def, nil, format, g.fs.DefaultDelimiter(), true)
	val.parse = parse
	return registerSliceValue(g, field, def, usage, format, val)
}

func registerSliceValue[T any](g *Group, field string, def []T, usage string,
	format func(T) string, val *DynamicSliceValue[T],
) *SliceFlag[T] {

	// Construct CLI-facing flag placeholder with default value
	bf := &core.BaseFlag{
		Name:  field,
		Usage: usage,
		Value: &slicePlaceholder{def: utils.JoinFormatted(def, format)},
	}

	// Register flag and value in the group
	g.items[field] = core.GroupItem{Value: val, Flag: bf}
	g.itemOrder = append(g.itemOrder, bf)

	// Return wrapper with typed access
	return &SliceFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        val,
	}
}
