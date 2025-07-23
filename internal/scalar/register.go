package scalar

import (
	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// scalarValueProvider is the interface for scalar flags.
type scalarValueProvider[T any] interface {
	core.Value
	Base() *ScalarValue[T]
}

// RegisterScalar registers a scalar flag.
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
		StaticFlag: builder.NewStaticFlag(reg, bf, ptr),
		val:        base,
	}
}
