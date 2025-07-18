package tinyflags

type sliceValueProvider[T any] interface {
	Value
	FlagBaseProvider[[]T]
}

type scalarValueProvider[T any] interface {
	Value
	FlagBaseProvider[T]
}

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
		builderImpl: builderImpl[T]{
			fs:    f,
			bf:    bf,
			value: base,
			ptr:   ptr,
		},
	}
}

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

	return &SliceFlag[T]{
		Flag: Flag[[]T]{
			builderImpl: builderImpl[[]T]{
				fs:    f,
				bf:    bf,
				value: base,
				ptr:   ptr,
			},
		},
	}
}

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
