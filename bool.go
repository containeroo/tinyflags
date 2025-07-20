package tinyflags

import "github.com/containeroo/tinyflags/internal/scalar"

func (f *FlagSet) Bool(name string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVarP(new(bool), name, "", def, usage)
}

func (f *FlagSet) BoolP(name, short string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVarP(new(bool), name, short, def, usage)
}

func (f *FlagSet) BoolVar(name, short string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVarP(new(bool), name, short, def, usage)
}

func (f *FlagSet) BoolVarP(ptr *bool, name, short string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVarP(ptr, name, short, def, usage)
}
