package tinyflags

func (fs *FlagSet) String(name, usage string, def string) *Flag[string] {
	return fs.StringVarP(new(string), name, "", def, usage)
}

func (fs *FlagSet) StringP(name, short string, def string, usage string) *Flag[string] {
	return fs.StringVarP(new(string), name, short, def, usage)
}

func (fs *FlagSet) StringVar(ptr *string, name string, def string, usage string) *Flag[string] {
	return fs.StringVarP(ptr, name, "", def, usage)
}

// String defines a scalar string flag (e.g. --name=value).
func (fs *FlagSet) StringVarP(ptr *string, name, short string, def string, usage string) *Flag[string] {
	value := NewFlagBase(
		ptr,
		def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)
	return register(fs, name, "", usage, value, ptr)
}

func (g *DynamicGroup) String(field, usage string) *DynamicFlag[string] {
	item := NewDynamicItem(
		field,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)

	g.items[field] = item

	// Register the dynamic pattern (e.g. http.*.address) in the main FlagSet
	g.fs.registerDynamic(g.prefix, field, item)

	bf := &baseFlag{
		name:  field,
		usage: usage,
	}

	return &DynamicFlag[string]{
		builderBase: builderBase[string]{
			fs:    g.fs,
			bf:    bf,
			value: nil,
			ptr:   nil,
		},
		item: item,
	}
}
