package dynamic

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

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
	val := &BoolValue{
		DynamicScalarValue: NewDynamicScalarValue(field, def, parse, format),
		Strict:             strict,
	}

	bf := &core.BaseFlag{
		Name:  field,
		Usage: usage,
		Value: &boolPlaceholder{def: format(def), strict: strict},
	}
	g.items[field] = core.GroupItem{Value: val, Flag: bf}
	g.itemOrder = append(g.itemOrder, bf)
	return &BoolFlag{
		DynamicFlag: builder.NewDynamicFlag[bool](g.fs, bf),
		item:        val.Base(),
		strictMode:  strict,
	}
}
