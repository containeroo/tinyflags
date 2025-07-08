package tinyflags

// addScalar registers a scalar flag and returns its builder.
func addScalar[T any](f *FlagSet, name, short, usage string, val Value, ptr *T) *Flag[T] {
	bf := &baseFlag{
		name:  name,  // long name: --flag
		short: short, // short name: -f
		usage: usage, // help text
		value: val,   // parsed value
	}
	f.flags[name] = bf                      // register in lookup map
	f.registered = append(f.registered, bf) // preserve order

	return &Flag[T]{fs: f, bf: bf, ptr: ptr}
}

// addSlice registers a slice flag and returns its slice builder.
func addSlice[T any](f *FlagSet, name, short, usage string, val Value, ptr *T) *SliceFlag[T] {
	bf := &baseFlag{
		name:  name,
		short: short,
		usage: usage,
		value: val,
	}
	f.flags[name] = bf
	f.registered = append(f.registered, bf)

	return &SliceFlag[T]{Flag: Flag[T]{fs: f, bf: bf, ptr: ptr}}
}
