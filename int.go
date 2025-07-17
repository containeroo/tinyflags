package tinyflags

import "strconv"

func (fs *FlagSet) Int(name string, def int, usage string) *Flag[int] {
	return fs.IntVarP(new(int), name, "", def, usage)
}

func (fs *FlagSet) IntP(name string, short string, def int, usage string) *Flag[int] {
	return fs.IntVarP(new(int), name, short, def, usage)
}

func (fs *FlagSet) IntVar(ptr *int, name string, def int, usage string) *Flag[int] {
	return fs.IntVarP(ptr, name, "", def, usage)
}

// Int defines a scalar int flag (e.g. --port=8080).
func (fs *FlagSet) IntVarP(ptr *int, name, short string, def int, usage string) *Flag[int] {
	value := NewFlagBase(
		ptr,
		def,
		strconv.Atoi,
		func(i int) string { return strconv.Itoa(i) },
	)
	return register(fs, name, short, usage, value, ptr)
}

// Int creates a dynamic int flag under the group (e.g. --group.id.port=8080).
func (g *DynamicGroup) Int(field, usage string) *DynamicFlag[int] {
	item := NewDynamicItem(
		field,
		strconv.Atoi,
		func(i int) string { return strconv.Itoa(i) },
	)

	g.items[field] = item

	// Register the dynamic pattern (e.g. group.*.port) in the main FlagSet
	g.fs.registerDynamic(g.prefix, field, item)

	bf := &baseFlag{
		name:  field,
		usage: usage,
	}

	return &DynamicFlag[int]{
		builderBase: builderBase[int]{
			fs:    g.fs,
			bf:    bf,
			value: nil,
			ptr:   nil,
		},
		item: item,
	}
}
