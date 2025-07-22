package scalar

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

type scalarValueProvider[T any] interface {
	core.Value
	Base() *ScalarValue[T]
}

func RegisterScalar[T any](
	reg core.Registry,
	name, usage string,
	val scalarValueProvider[T],
	ptr *T,
) *ScalarFlag[T] {
	base := val.Base()

	bf := &core.BaseFlag{
		Name:  name,
		Usage: usage,
		Value: val,
	}
	reg.RegisterFlag(name, bf)

	return &ScalarFlag[T]{
		StaticFlag: builder.NewDefaultFlag(reg, bf, ptr),
		val:        base,
	}
}
