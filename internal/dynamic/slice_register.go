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
	// Create a slice value with default delimiter from the flagset
	val := NewDynamicSliceValue(field, def, parse, format, g.fs.DefaultDelimiter())

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
