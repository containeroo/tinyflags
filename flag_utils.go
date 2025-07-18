package tinyflags

// sliceValueProvider is a constraint for slice flag values that provide their FlagBase.
type sliceValueProvider[T any] interface {
	Value
	FlagBaseProvider[[]T]
}

// scalarValueProvider is a constraint for scalar flag values that provide their FlagBase.
type scalarValueProvider[T any] interface {
	Value
	FlagBaseProvider[T]
}

// addScalar registers a scalar flag and returns its builder.
func addScalar[T any](f *FlagSet, name, short, usage string, val scalarValueProvider[T], ptr *T,
) *Flag[T] {
	base := val.Base()

	bf := &baseFlag{
		name:  name,
		short: short,
		usage: usage,
		value: val,
	}
	f.flags[name] = bf
	f.registered = append(f.registered, bf)

	return &Flag[T]{
		builderBase: builderBase[T]{
			fs:    f,
			bf:    bf,
			value: base,
			ptr:   ptr,
		},
	}
}

// addSlice registers a slice flag and returns its builder.
func addSlice[T any](
	f *FlagSet,
	name, short, usage string,
	val sliceValueProvider[T],
	ptr *[]T,
) *SliceFlag[T] {
	base := val.Base()

	bf := &baseFlag{
		name:  name,
		short: short,
		usage: usage,
		value: val,
	}

	f.flags[name] = bf
	f.registered = append(f.registered, bf)

	return &SliceFlag[T]{ // ← return correct wrapper
		Flag: Flag[[]T]{
			builderBase: builderBase[[]T]{
				fs:    f,
				bf:    bf,
				value: base,
				ptr:   ptr,
			},
		},
	}
}

// addDynamic registers a dynamic flag (e.g., --group.id.field) in the FlagSet.
// It stores the item under fs.dynamic[group][field] for later parsing.
func addDynamic(
	fs *FlagSet,
	group string,
	field string,
	item DynamicValue,
) {
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
}
