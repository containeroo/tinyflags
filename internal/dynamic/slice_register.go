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
) *SliceFlag[T] {
	val := NewDynamicSliceValue(field, def, parse, format, g.fs.DefaultDelimiter())

	bf := &core.BaseFlag{
		Name:  field,
		Usage: usage,
		Value: &slicePlaceholder{def: utils.JoinFormatted(def, format)},
	}
	g.items[field] = core.GroupItem{Value: val, Flag: bf}
	g.itemOrder = append(g.itemOrder, bf)
	return &SliceFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        val,
	}
}
