package tinyflags

import "strconv"

// Int defines a scalar int flag with a default value.
func (f *FlagSet) Int(name string, def int, usage string) *Flag[int] {
	return f.IntVarP(new(int), name, "", def, usage)
}

// IntP defines a scalar int flag with a short name.
func (f *FlagSet) IntP(name, short string, def int, usage string) *Flag[int] {
	return f.IntVarP(new(int), name, short, def, usage)
}

// IntVar defines a scalar int flag and binds it to a variable.
func (f *FlagSet) IntVar(ptr *int, name string, def int, usage string) *Flag[int] {
	return f.IntVarP(ptr, name, "", def, usage)
}

// IntVarP defines a scalar int flag with a short name and binds it to a variable.
func (f *FlagSet) IntVarP(ptr *int, name, short string, def int, usage string) *Flag[int] {
	val := NewValueImpl(
		ptr,
		def,
		strconv.Atoi,
		strconv.Itoa,
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// Int creates a dynamic int flag under the group (e.g. --group.id.port=8080).
func (g *DynamicGroup) Int(field, usage string) *DynamicFlag[int] {
	item := NewDynamicItemImpl(
		field,
		strconv.Atoi,
		strconv.Itoa,
	)

	g.items[field] = item

	return addDynamic(g.fs, g.prefix, field, usage, item)
}
