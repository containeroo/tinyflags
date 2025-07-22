package tinyflags

import "github.com/containeroo/tinyflags/internal/scalar"

// Float32 defines a scalar float32 flag with default value.
func (f *FlagSet) Float32(name string, def float32, usage string) *scalar.ScalarFlag[float32] {
	return f.impl.Float32Var(new(float32), name, def, usage)
}

func (f *FlagSet) Float32Var(ptr *float32, name, short string, def float32, usage string) *scalar.ScalarFlag[float32] {
	return f.impl.Float32Var(ptr, name, def, usage)
}
