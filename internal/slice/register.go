package slice

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// ValueProvider is the interface for slice flags.
type ValueProvider[T any] interface {
	core.Value
	Base() *SliceValue[T]
}

// RegisterSlice registers a slice flag.
func RegisterSlice[T any](
	reg core.Registry,
	name, usage string,
	val ValueProvider[T],
	ptr *[]T,
) *SliceFlag[T] {
	bf := &core.BaseFlag{
		Name:  name,
		Usage: usage,
		Value: val,
	}

	reg.RegisterFlag(name, bf)

	flag := &SliceFlag[T]{}
	flag.StaticFlag = builder.NewStaticFlag(reg, bf, ptr, flag)
	flag.val = val.Base()
	return flag
}
