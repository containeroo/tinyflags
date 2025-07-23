package engine

import "github.com/containeroo/tinyflags/internal/scalar"

func (f *FlagSet) BoolVar(ptr *bool, name string, def bool, usage string) *scalar.BoolFlag {
	return scalar.NewBool(f, ptr, name, usage, def)
}

func (f *FlagSet) Bool(name string, def bool, usage string) *scalar.BoolFlag {
	return scalar.NewBool(f, new(bool), name, usage, def)
}
