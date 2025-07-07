package tinyflags

import "strconv"

// IntSliceP defines a []int slice flag with the specified name, shorthand, default value, and usage.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IntSliceP(name, short string, def []int, usage string) *SliceFlag[[]int] {
	ptr := new([]int)
	val := NewSliceItem(ptr, def, strconv.Atoi, strconv.Itoa, f.defaultDelimiter)
	return addSlice(f, name, short, usage, val, ptr)
}

// IntSlice defines a []int slice flag with the specified name, default value, and usage.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IntSlice(name string, def []int, usage string) *SliceFlag[[]int] {
	return f.IntSliceP(name, "", def, usage)
}

// IntSliceVarP defines a []int slice flag with the specified name, shorthand, default value, and usage.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IntSliceVarP(ptr *[]int, name, short string, def []int, usage string) *SliceFlag[[]int] {
	val := NewSliceItem(ptr, def, strconv.Atoi, strconv.Itoa, f.defaultDelimiter)
	return addSlice(f, name, short, usage, val, ptr)
}

// IntSliceVar defines a []int slice flag with the specified name, default value, and usage.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IntSliceVar(ptr *[]int, name string, def []int, usage string) *SliceFlag[[]int] {
	return f.IntSliceVarP(ptr, name, "", def, usage)
}
