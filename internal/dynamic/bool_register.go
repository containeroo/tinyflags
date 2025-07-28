package dynamic

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// registerDynamicBool registers a dynamic boolean flag under the group.
func registerDynamicBool(
	g *Group,
	field string,
	def bool,
	usage string,
	parse func(string) (bool, error),
	format func(bool) string,
) *BoolFlag {
	// Shared pointer to track strict mode behavior
	strict := new(bool)

	// Create a BoolValue that wraps a dynamic scalar value and strict marker
	val := &BoolValue{
		DynamicScalarValue: NewDynamicScalarValue(field, def, parse, format),
		strictMode:         strict,
	}

	// Create placeholder BaseFlag used for help and CLI parsing
	bf := &core.BaseFlag{
		Name:  field,
		Usage: usage,
		Value: &boolPlaceholder{def: format(def), strictMode: strict},
	}

	// Register flag and value in the group
	g.items[field] = core.GroupItem{Value: val, Flag: bf}
	g.itemOrder = append(g.itemOrder, bf)

	// Return wrapped BoolFlag for external access
	return &BoolFlag{
		DynamicFlag: builder.NewDynamicFlag[bool](g.fs, bf),
		item:        val.Base(),
		strictMode:  strict,
	}
}
