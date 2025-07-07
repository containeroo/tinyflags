package tinyflags

import "time"

// DurationSliceP defines a duration flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) DurationSliceP(name, short string, def []time.Duration, usage string) *SliceFlag[[]time.Duration] {
	ptr := new([]time.Duration)
	val := NewSliceItem(
		ptr,
		def,
		time.ParseDuration,
		func(d time.Duration) string { return d.String() },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// DurationSlice defines a duration slice flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) DurationSlice(name string, def []time.Duration, usage string) *SliceFlag[[]time.Duration] {
	return f.DurationSliceP(name, "", def, usage)
}

// DurationSliceVarP defines a duration slice flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) DurationSliceVarP(ptr *[]time.Duration, name, short string, def []time.Duration, usage string) *SliceFlag[[]time.Duration] {
	val := NewSliceItem(
		ptr,
		def,
		time.ParseDuration,
		func(d time.Duration) string { return d.String() },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// DurationSliceVar defines a duration slice flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) DurationSliceVar(ptr *[]time.Duration, name string, def []time.Duration, usage string) *SliceFlag[[]time.Duration] {
	val := NewSliceItem(
		ptr,
		def,
		time.ParseDuration,
		func(d time.Duration) string { return d.String() },
		f.defaultDelimiter,
	)
	return addSlice(f, name, "", usage, val, ptr)
}
