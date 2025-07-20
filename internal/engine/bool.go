package engine

import "github.com/containeroo/tinyflags/internal/scalar"

// BoolP defines a boolean flag with a short name and a default value.
func (f *FlagSet) BoolVarP(ptr *bool, name, short string, def bool, usage string) *scalar.BoolFlag {
	return scalar.NewBool(f, ptr, name, short, usage, def)
}

func (f *FlagSet) BoolVar(name, short string, def bool, usage string) *scalar.BoolFlag {
	return scalar.NewBool(f, new(bool), name, short, usage, def)
}

func (f *FlagSet) Bool(name string, def bool, usage string) *scalar.BoolFlag {
	return scalar.NewBool(f, new(bool), name, "", usage, def)
}

func (f *FlagSet) BoolP(name, short string, def bool, usage string) *scalar.BoolFlag {
	return scalar.NewBool(f, new(bool), name, short, usage, def)
}
