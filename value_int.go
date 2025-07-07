package tinyflags

import "strconv"

// IntP defines an int flag with specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IntP(name, short string, def int, usage string) *Flag[int] {
	ptr := new(int)
	val := NewFlagItem(ptr, def, strconv.Atoi, strconv.Itoa)
	return addScalar(f, name, short, usage, val, ptr)
}

// Int defines an int flag with specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Int(name string, def int, usage string) *Flag[int] {
	return f.IntP(name, "", def, usage)
}

// IntVarP defines an int flag with specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IntVarP(ptr *int, name, short string, def int, usage string) *Flag[int] {
	val := NewFlagItem(ptr, def, strconv.Atoi, strconv.Itoa)
	return addScalar(f, name, short, usage, val, ptr)
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IntVar(ptr *int, name string, def int, usage string) *Flag[int] {
	return f.IntVarP(ptr, name, "", def, usage)
}
