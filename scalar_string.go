package tinyflags

// String defines a scalar string flag with a default value.
func (f *FlagSet) String(name string, def string, usage string) *Flag[string] {
	return f.StringVarP(new(string), name, "", def, usage)
}

// StringP defines a string flag with a short name.
func (f *FlagSet) StringP(name, short string, def string, usage string) *Flag[string] {
	return f.StringVarP(new(string), name, short, def, usage)
}

// StringVar defines a string flag and binds it to a variable.
func (f *FlagSet) StringVar(ptr *string, name string, def string, usage string) *Flag[string] {
	return f.StringVarP(ptr, name, "", def, usage)
}

// StringVarP defines a scalar string flag with a short name and binds it to a variable.
func (f *FlagSet) StringVarP(ptr *string, name, short, def, usage string) *Flag[string] {
	val := NewScalarValueImpl(
		ptr,
		def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)
	return addScalar(f, name, short, usage, val, ptr)
}
