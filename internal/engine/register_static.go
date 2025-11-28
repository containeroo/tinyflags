package engine

import (
	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/scalar"
	"github.com/containeroo/tinyflags/internal/slice"
)

// RegisterStaticScalar centralizes scalar flag registration.
func RegisterStaticScalar[T any](
	reg core.Registry,
	ptr *T,
	name, usage string,
	def T,
	parse func(string) (T, error),
	format func(T) string,
) *scalar.ScalarFlag[T] {
	val := scalar.NewScalarValue(ptr, def, parse, format)
	return scalar.RegisterScalar(reg, name, usage, val, ptr)
}

// RegisterStaticSlice centralizes slice flag registration.
func RegisterStaticSlice[T any](
	reg core.Registry,
	ptr *[]T,
	name, usage string,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *slice.SliceFlag[T] {
	val := slice.NewSliceValue(ptr, def, parse, format, delimiter)
	return slice.RegisterSlice(reg, name, usage, val, ptr)
}
