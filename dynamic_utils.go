package tinyflags

// addDynamic registers a dynamic scalar flag and returns its builder.
func addDynamic[T any](
	fs *FlagSet,
	group string,
	field string,
	usage string,
	item *DynamicItemImpl[T],
) *DynamicFlag[T] {
	if fs.dynamic == nil {
		fs.dynamic = make(map[string]map[string]DynamicValue)
	}
	if _, ok := fs.dynamic[group]; !ok {
		fs.dynamic[group] = make(map[string]DynamicValue)
	}
	if _, exists := fs.dynamic[group][field]; exists {
		panic("addDynamic: duplicate dynamic flag registration for " + group + "." + field)
	}
	fs.dynamic[group][field] = item

	bf := &baseFlag{
		name:  field,
		usage: usage,
	}

	return &DynamicFlag[T]{
		fs:   fs,
		bf:   bf,
		item: item,
	}
}

// addDynamicSlice registers a dynamic slice flag and returns its builder.
func addDynamicSlice[T any](
	fs *FlagSet,
	group string,
	field string,
	usage string,
	item *DynamicSliceItemImpl[T],
) *DynamicSliceFlag[T] {
	if fs.dynamic == nil {
		fs.dynamic = make(map[string]map[string]DynamicValue)
	}
	if _, ok := fs.dynamic[group]; !ok {
		fs.dynamic[group] = make(map[string]DynamicValue)
	}
	if _, exists := fs.dynamic[group][field]; exists {
		panic("addDynamicSlice: duplicate dynamic flag registration for " + group + "." + field)
	}
	fs.dynamic[group][field] = item

	bf := &baseFlag{
		name:  field,
		usage: usage,
	}

	return &DynamicSliceFlag[T]{
		fs:   fs,
		bf:   bf,
		item: item,
	}
}

// addDynamicBool registers a dynamic bool flag and returns its builder.
func addDynamicBool(
	fs *FlagSet,
	group string,
	field string,
	usage string,
	item *DynamicItemImpl[bool],
) *DynamicBoolFlag {
	if fs.dynamic == nil {
		fs.dynamic = make(map[string]map[string]DynamicValue)
	}
	if _, ok := fs.dynamic[group]; !ok {
		fs.dynamic[group] = make(map[string]DynamicValue)
	}
	if _, exists := fs.dynamic[group][field]; exists {
		panic("addDynamicBool: duplicate dynamic flag registration for " + group + "." + field)
	}
	fs.dynamic[group][field] = item

	bf := &baseFlag{
		name:  field,
		usage: usage,
	}

	return &DynamicBoolFlag{
		fs:   fs,
		bf:   bf,
		item: item,
	}
}
