package tinyflags

// StringSlice defines a string slice flag with a default value.
func (fs *FlagSet) StringSlice(name string, def []string, usage string) *SliceFlag[string] {
	return fs.StringSliceVarP(new([]string), name, "", def, usage)
}

// StringSliceP defines a string slice flag with a short name.
func (fs *FlagSet) StringSliceP(name, short string, def []string, usage string) *SliceFlag[string] {
	return fs.StringSliceVarP(new([]string), name, short, def, usage)
}

// StringSliceVar defines a string slice flag and binds it to a variable.
func (fs *FlagSet) StringSliceVar(ptr *[]string, name string, def []string, usage string) *SliceFlag[string] {
	return fs.StringSliceVarP(ptr, name, "", def, usage)
}

// StringSliceVarP defines a string slice flag with a short name and binds it to a variable.
func (fs *FlagSet) StringSliceVarP(ptr *[]string, name, short string, def []string, usage string) *SliceFlag[string] {
	val := NewSliceValueImpl(
		ptr,
		def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
		fs.defaultDelimiter,
	)
	return addSlice(fs, name, short, usage, val, ptr)
}

// StringSlice defines a dynamic string slice flag under the group (e.g. --http.alpha.tags=one,two).
func (g *DynamicGroup) StringSlice(field, usage string) *DynamicSliceFlag[string] {
	item := NewDynamicSliceItem(
		field,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
		g.fs.defaultDelimiter, // assumes FlagSet has defaultDelimiter set
	)

	g.items[field] = item

	// Register the dynamic pattern (e.g. http.*.tags)
	addDynamic(g.fs, g.prefix, field, item)

	bf := &baseFlag{
		name:  field,
		usage: usage,
	}

	return &DynamicSliceFlag[string]{
		builderImpl: builderImpl[[]string]{
			fs:  g.fs,
			bf:  bf,
			ptr: nil, // dynamic flags don't use a pre-bound pointer
			// value remains nil for dynamic slice flags
		},
		item: item,
	}
}
