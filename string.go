package tinyflags

// String defines a scalar string flag with a default value.
func (fs *FlagSet) String(name string, def string, usage string) *Flag[string] {
	return fs.StringVarP(new(string), name, "", def, usage)
}

// StringP defines a string flag with a short name.
func (fs *FlagSet) StringP(name, short string, def string, usage string) *Flag[string] {
	return fs.StringVarP(new(string), name, short, def, usage)
}

// StringVar defines a string flag and binds it to a variable.
func (fs *FlagSet) StringVar(ptr *string, name string, def string, usage string) *Flag[string] {
	return fs.StringVarP(ptr, name, "", def, usage)
}

// StringVarP defines a scalar string flag with a short name and binds it to a variable.
func (fs *FlagSet) StringVarP(ptr *string, name, short, def, usage string) *Flag[string] {
	val := NewValueImpl(
		ptr,
		def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)
	return addScalar(fs, name, short, usage, val, ptr)
}

func (g *DynamicGroup) String(field, usage string) *DynamicFlag[string] {
	item := NewDynamicItem(
		field,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)

	g.items[field] = item

	// Register the dynamic pattern (e.g. http.*.address) in the main FlagSet
	addDynamic(g.fs, g.prefix, field, item)

	bf := &baseFlag{
		name:  field,
		usage: usage,
	}

	return &DynamicFlag[string]{
		builderImpl: builderImpl[string]{
			fs:    g.fs,
			bf:    bf,
			value: nil,
			ptr:   nil,
		},
		item: item,
	}
}
