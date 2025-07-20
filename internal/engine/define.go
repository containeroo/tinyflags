package engine

import (
	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/scalar"
	"github.com/containeroo/tinyflags/internal/slice"
)

// defineScalar creates and registers a scalar flag for the given FlagSet.
// It supports optional pointer binding via the `ptr` argument.
func defineScalar[T any](
	r core.Registry,
	ptr *T,
	name, short, usage string,
	def T,
	parse func(string) (T, error),
	format func(T) string,
) *scalar.ScalarFlag[T] {
	val := scalar.NewScalarValue(ptr, def, parse, format)
	return scalar.RegisterScalar(r, name, short, usage, val, ptr)
}

// defineSlice creates and registers a slice flag for the given Registry.
// It supports optional pointer binding via the `ptr` argument.
func defineSlice[T any](
	r core.Registry,
	ptr *[]T,
	name, short, usage string,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
) *slice.SliceFlag[T] {
	val := slice.NewSliceValue(ptr, def, parse, format, r.DefaultDelimiter())
	return slice.RegisterSlice(r, name, short, usage, val, ptr)
}
