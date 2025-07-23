package engine

import (
	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/scalar"
	"github.com/containeroo/tinyflags/internal/slice"
)

// registerScalar creates and registers a scalar flag for the given FlagSet.
// It supports optional pointer binding via the `ptr` argument.
func registerScalar[T any](
	r core.Registry,
	ptr *T,
	name, usage string,
	def T,
	parse func(string) (T, error),
	format func(T) string,
) *scalar.ScalarFlag[T] {
	val := scalar.NewScalarValue(ptr, def, parse, format)
	return scalar.RegisterScalar(r, name, usage, val, ptr)
}

// registerSlice creates and registers a slice flag for the given Registry.
// It supports optional pointer binding via the `ptr` argument.
func registerSlice[T any](
	r core.Registry,
	ptr *[]T,
	name, usage string,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
) *slice.SliceFlag[T] {
	val := slice.NewSliceValue(ptr, def, parse, format, r.DefaultDelimiter())
	return slice.RegisterSlice(r, name, usage, val, ptr)
}
