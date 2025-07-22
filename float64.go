package tinyflags

import "github.com/containeroo/tinyflags/internal/scalar"

// Float64 defines a scalar float64 flag with default value.
func (f *FlagSet) Float64(name string, def float64, usage string) *scalar.ScalarFlag[float64] {
	return f.impl.Float64Var(new(float64), name, def, usage)
}

func (f *FlagSet) Float64Var(ptr *float64, name string, def float64, usage string) *scalar.ScalarFlag[float64] {
	return f.impl.Float64Var(ptr, name, def, usage)
}
