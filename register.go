package tinyflags

// register creates and registers a new typed scalar flag in the FlagSet.
func register[T any](
	fs *FlagSet,
	name, short, usage string,
	value *FlagBase[T],
	ptr *T,
) *Flag[T] {
	bf := &baseFlag{
		name:  name,
		short: short,
		usage: usage,
		value: value,
	}

	if fs.flags == nil {
		fs.flags = make(map[string]*baseFlag)
	}
	fs.flags[name] = bf

	return &Flag[T]{
		builderBase: builderBase[T]{
			fs:    fs,
			bf:    bf,
			value: value,
			ptr:   ptr,
		},
	}
}

// registerDynamic lets the FlagSet recognize dynamic flags like --group.id.field
func (fs *FlagSet) registerDynamic(group, field string, item DynamicValue) {
	if fs.dynamic == nil {
		fs.dynamic = make(map[string]map[string]DynamicValue)
	}
	if fs.dynamic[group] == nil {
		fs.dynamic[group] = make(map[string]DynamicValue)
	}
	fs.dynamic[group][field] = item
}
