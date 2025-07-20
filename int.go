package tinyflags

import "github.com/containeroo/tinyflags/internal/scalar"

// Int defines a scalar int flag with default value.
func (f *FlagSet) Int(name string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.impl.IntVarP(new(int), name, "", def, usage)
}

func (f *FlagSet) IntVar(ptr *int, name, short string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.impl.IntVarP(ptr, name, short, def, usage)
}

func (f *FlagSet) IntP(name, short string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.impl.IntVarP(new(int), name, short, def, usage)
}

func (f *FlagSet) IntVarP(ptr *int, name, short string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.impl.IntVarP(ptr, name, short, def, usage)
}
