package slice

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

type sliceValueProvider[T any] interface {
	core.Value
	Base() *SliceValue[T]
}

func RegisterSlice[T any](
	reg core.Registry,
	name, usage string,
	val sliceValueProvider[T],
	ptr *[]T,
) *SliceFlag[T] {
	base := val.Base()

	bf := &core.BaseFlag{
		Name:  name,
		Usage: usage,
		Value: val,
	}

	reg.RegisterFlag(name, bf)

	return &SliceFlag[T]{
		StaticFlag: builder.NewDefaultFlag(reg, bf, ptr),
		val:        base,
	}
}
