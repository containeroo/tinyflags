package tinyflags

import "github.com/containeroo/tinyflags/internal/scalar"

// Int defines a scalar int flag with default value.
func (f *FlagSet) Int(name string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.impl.IntVar(new(int), name, def, usage)
}

func (f *FlagSet) IntVar(ptr *int, name string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.impl.IntVar(ptr, name, def, usage)
}
