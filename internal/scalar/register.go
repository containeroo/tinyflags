package scalar

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// ValueProvider is the interface for scalar flag values.
type ValueProvider[T any] interface {
	core.Value
	Base() *ScalarValue[T]
}

// RegisterScalar registers a scalar flag.
func RegisterScalar[T any](
	reg core.Registry,
	name, usage string,
	val ValueProvider[T],
	ptr *T,
) *ScalarFlag[T] {
	bf := &core.BaseFlag{
		Name:  name,
		Usage: usage,
		Value: val,
	}
	reg.RegisterFlag(name, bf)

	flag := &ScalarFlag[T]{}
	flag.scalarFlagBase = scalarFlagBase[T, *ScalarFlag[T]]{
		StaticFlag: builder.NewStaticFlag(reg, bf, ptr, flag),
		val:        val.Base(),
		self:       flag,
	}
	return flag
}
