package engine

import "github.com/containeroo/tinyflags/internal/slice"

// String defines a scalar string flag with default value.
func (f *FlagSet) StringSliceVar(ptr *[]string, name string, usage string, def []string) *slice.SliceFlag[string] {
	return defineSlice(f, ptr, name, usage, def, func(s string) (string, error) { return s, nil }, func(s string) string { return s })
}
