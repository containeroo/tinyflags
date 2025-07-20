package tinyflags

import "github.com/containeroo/tinyflags/internal/slice"

func (f *FlagSet) StringSlice(name string, def []string, usage string) *slice.SliceFlag[string] {
	return f.impl.StringSliceVarP(new([]string), name, "", usage, def)
}

func (f *FlagSet) StringSliceP(name, short string, def []string, usage string) *slice.SliceFlag[string] {
	return f.impl.StringSliceVarP(new([]string), name, short, usage, def)
}

func (f *FlagSet) StringSliceVar(name string, def []string, usage string) *slice.SliceFlag[string] {
	return f.impl.StringSliceVarP(new([]string), name, "", usage, def)
}

func (f *FlagSet) StringSliceVarP(ptr *[]string, name, short string, usage string, def []string) *slice.SliceFlag[string] {
	return f.impl.StringSliceVarP(ptr, name, short, usage, def)
}
