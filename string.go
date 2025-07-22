package tinyflags

import "github.com/containeroo/tinyflags/internal/scalar"

// String defines a scalar string flag with default value.
func (f *FlagSet) String(name, def, usage string) *scalar.ScalarFlag[string] {
	return f.impl.StringVar(new(string), name, usage, def)
}

func (f *FlagSet) StringVar(ptr *string, name string, def string, usage string) *scalar.ScalarFlag[string] {
	return f.impl.StringVar(ptr, name, usage, def)
}
