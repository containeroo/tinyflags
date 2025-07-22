package tinyflags

import "github.com/containeroo/tinyflags/internal/scalar"

func (f *FlagSet) Bool(name string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVar(new(bool), name, def, usage)
}

func (f *FlagSet) BoolVar(ptr *bool, name, short string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVar(ptr, name, def, usage)
}
