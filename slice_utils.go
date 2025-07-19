package tinyflags

type sliceValueProvider[T any] interface {
	Value
	FlagBaseProvider[[]T]
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
