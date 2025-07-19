package tinyflags

// StringSlice defines a string slice flag with a default value.
func (f *FlagSet) StringSlice(name string, def []string, usage string) *SliceFlag[string] {
	return f.StringSliceVarP(new([]string), name, "", def, usage)
}

// StringSliceP defines a string slice flag with a short name.
func (f *FlagSet) StringSliceP(name, short string, def []string, usage string) *SliceFlag[string] {
	return f.StringSliceVarP(new([]string), name, short, def, usage)
}

// StringSliceVar defines a string slice flag and binds it to a variable.
func (f *FlagSet) StringSliceVar(ptr *[]string, name string, def []string, usage string) *SliceFlag[string] {
	return f.StringSliceVarP(ptr, name, "", def, usage)
}

// StringSliceVarP defines a string slice flag with a short name and binds it to a variable.
func (f *FlagSet) StringSliceVarP(ptr *[]string, name, short string, def []string, usage string) *SliceFlag[string] {
	val := NewSliceValueImpl(
		ptr,
		def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}
