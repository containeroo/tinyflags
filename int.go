package tinyflags

import "strconv"

// Int defines a scalar int flag with a default value.
func (fs *FlagSet) Int(name string, def int, usage string) *Flag[int] {
	return fs.IntVarP(new(int), name, "", def, usage)
}

// IntP defines a scalar int flag with a short name.
func (fs *FlagSet) IntP(name, short string, def int, usage string) *Flag[int] {
	return fs.IntVarP(new(int), name, short, def, usage)
}

// IntVar defines a scalar int flag and binds it to a variable.
func (fs *FlagSet) IntVar(ptr *int, name string, def int, usage string) *Flag[int] {
	return fs.IntVarP(ptr, name, "", def, usage)
}

// IntVarP defines a scalar int flag with a short name and binds it to a variable.
func (fs *FlagSet) IntVarP(ptr *int, name, short string, def int, usage string) *Flag[int] {
	val := NewFlagValue(
		ptr,
		def,
		strconv.Atoi,
		strconv.Itoa,
	)
	return addScalar(fs, name, short, usage, val, ptr)
}

// Int creates a dynamic int flag under the group (e.g. --group.id.port=8080).
func (g *DynamicGroup) Int(field, usage string) *DynamicFlag[int] {
	item := NewDynamicItem(
		field,
		strconv.Atoi,
		strconv.Itoa,
	)

	g.items[field] = item

	// Register the dynamic pattern (e.g. group.*.port) in the main FlagSet
	addDynamic(g.fs, g.prefix, field, item)

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
