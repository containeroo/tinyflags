package tinyflags

// StringP defines a string flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) StringP(name, short string, def string, usage string) *Flag[string] {
	ptr := new(string)
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// String defines a string flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) String(name string, def string, usage string) *Flag[string] {
	return f.StringP(name, "", def, usage)
}

// StringVarP defines a string flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) StringVarP(ptr *string, name, short string, def string, usage string) *Flag[string] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// StringVar defines a string flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) StringVar(ptr *string, name string, def string, usage string) *Flag[string] {
	return f.StringVarP(ptr, name, "", def, usage)
}
