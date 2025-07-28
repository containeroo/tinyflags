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
	usage string,
	parse func(string) (T, error),
	format func(T) string,
) *ScalarFlag[T] {
	// Build internal value container
	val := NewDynamicScalarValue(field, def, parse, format)

	// Construct base CLI flag metadata
	bf := &core.BaseFlag{
		Name:  field,
		Usage: usage,
		Value: &placeholderValue{def: format(def)},
	}

	// Register the flag and its value in the group
	g.items[field] = core.GroupItem{Value: val, Flag: bf}
	g.itemOrder = append(g.itemOrder, bf)

	// Return typed wrapper for external use
	return &ScalarFlag[T]{
		DynamicFlag: builder.NewDynamicFlag[T](g.fs, bf),
		item:        val.Base(),
	}
}
