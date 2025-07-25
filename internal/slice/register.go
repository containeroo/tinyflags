package slice

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// sliceValueProvider is the interface for slice flags.
type sliceValueProvider[T any] interface {
	core.Value
	Base() *SliceValue[T]
}

// RegisterSlice registers a slice flag.
func RegisterSlice[T any](
	reg core.Registry,
	name, usage string,
	val sliceValueProvider[T],
	ptr *[]T,
) *SliceFlag[T] {
	bf := &core.BaseFlag{
		Name:  name,
		Usage: usage,
		Value: val,
	}

	reg.RegisterFlag(name, bf)

	return &SliceFlag[T]{
		StaticFlag: builder.NewStaticFlag(reg, bf, ptr),
		val:        val.Base(),
	}
}
