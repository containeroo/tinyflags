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
	name, short, usage string,
	val scalarValueProvider[T],
	ptr *T,
) *ScalarFlag[T] {
	base := val.Base()

	bf := &core.BaseFlag{
		Name:  name,
		Short: short,
		Usage: usage,
		Value: val,
	}

	reg.RegisterFlag(name, bf)

	return &ScalarFlag[T]{
		DefaultFlag: builder.DefaultFlag[T]{
			Registry: reg,
			BF:       bf,
			Ptr:      ptr,
		},
		val: base,
	}
}
