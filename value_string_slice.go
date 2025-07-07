package tinyflags

// StringSliceP defines a string slice flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) StringSliceP(name, short string, def []string, usage string) *SliceFlag[[]string] {
	ptr := new([]string)
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// StringSlice defines a string slice flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) StringSlice(name string, def []string, usage string) *SliceFlag[[]string] {
	return f.StringSliceP(name, "", def, usage)
}

// StringSliceP defines a string flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) StringSliceVarP(ptr *[]string, name, short string, def []string, usage string) *SliceFlag[[]string] {
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// StringSlice defines a string slice flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) StringSliceVar(ptr *[]string, name string, def []string, usage string) *SliceFlag[[]string] {
	return f.StringSliceVarP(ptr, name, "", def, usage)
}
