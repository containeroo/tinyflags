package engine

import "github.com/containeroo/tinyflags/internal/slice"

// String defines a scalar string flag with default value.
func (f *FlagSet) StringSliceVarP(ptr *[]string, name, short string, usage string, def []string) *slice.SliceFlag[string] {
	return defineSlice(f, ptr, name, short, usage, def, func(s string) (string, error) { return s, nil }, func(s string) string { return s })
}
