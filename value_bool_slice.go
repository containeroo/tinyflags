package tinyflags

import "strconv"

// BoolSliceP defines a bool slice flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().

// BoolSlice defines a bool slice flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().

// BoolSliceP defines a bool flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().

// BoolSlice defines a bool slice flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().

func (f *FlagSet) BoolSliceP(name, short string, def []bool, usage string) *SliceFlag[[]bool] {
	ptr := new([]bool)
	val := NewSliceItem(ptr, def, strconv.ParseBool, strconv.FormatBool, f.defaultDelimiter)
	return addSlice(f, name, short, usage, val, ptr)
}

// BoolSlice defines a bool slice flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) BoolSlice(name string, def []bool, usage string) *SliceFlag[[]bool] {
	return f.BoolSliceP(name, "", def, usage)
}

// BoolSliceP defines a bool flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) BoolSliceVarP(ptr *[]bool, name, short string, def []bool, usage string) *SliceFlag[[]bool] {
	val := NewSliceItem(ptr, def, strconv.ParseBool, strconv.FormatBool, f.defaultDelimiter)
	return addSlice(f, name, short, usage, val, ptr)
}

// BoolSlice defines a bool slice flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) BoolSliceVar(ptr *[]bool, name string, def []bool, usage string) *SliceFlag[[]bool] {
	return f.BoolSliceVarP(ptr, name, "", def, usage)
}
