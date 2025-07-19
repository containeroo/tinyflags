package tinyflags

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
