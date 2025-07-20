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
	name, short, usage string,
	val sliceValueProvider[T],
	ptr *[]T,
) *SliceFlag[T] {
	base := val.Base()

	bf := &core.BaseFlag{
		Name:  name,
		Short: short,
		Usage: usage,
		Value: val,
	}

	reg.RegisterFlag(name, bf)

	return &SliceFlag[T]{
		DefaultFlag: builder.DefaultFlag[[]T]{
			Registry: reg,
			BF:       bf,
			Ptr:      ptr,
		},
		val: base,
	}
}
