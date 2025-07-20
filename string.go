package tinyflags

import "github.com/containeroo/tinyflags/internal/scalar"

// String defines a scalar string flag with default value.
func (f *FlagSet) String(name, usage, def string) *scalar.ScalarFlag[string] {
	return f.impl.StringVarP(new(string), name, "", usage, def)
}

func (f *FlagSet) StringVar(ptr *string, name, short string, def string, usage string) *scalar.ScalarFlag[string] {
	return f.impl.StringVarP(ptr, name, short, usage, def)
}

func (f *FlagSet) StringP(name, short string, def string, usage string) *scalar.ScalarFlag[string] {
	return f.impl.StringVarP(new(string), name, short, usage, def)
}

func (f *FlagSet) StringVarP(ptr *string, name, short string, usage string, def string) *scalar.ScalarFlag[string] {
	return f.impl.StringVarP(ptr, name, short, usage, def)
}
