package tinyflags

import "github.com/containeroo/tinyflags/internal/slice"

func (f *FlagSet) StringSlice(name string, def []string, usage string) *slice.SliceFlag[string] {
	return f.impl.StringSliceVar(new([]string), name, usage, def)
}

func (f *FlagSet) StringSliceVar(ptr *[]string, name string, usage string, def []string) *slice.SliceFlag[string] {
	return f.impl.StringSliceVar(ptr, name, usage, def)
}
