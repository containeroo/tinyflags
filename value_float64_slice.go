package tinyflags

import "strconv"

// Float64SliceP defines a float64 slice flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Float64SliceP(name, short string, def []float64, usage string) *SliceFlag[[]float64] {
	ptr := new([]float64)
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (float64, error) { return strconv.ParseFloat(s, 64) },
		func(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// Float64Slice defines a float64 slice flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Float64Slice(name string, def []float64, usage string) *SliceFlag[[]float64] {
	return f.Float64SliceP(name, "", def, usage)
}

// Float64SliceVarP defines a float64 slice flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) Float64SliceVarP(ptr *[]float64, name, short string, def []float64, usage string) *SliceFlag[[]float64] {
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (float64, error) { return strconv.ParseFloat(s, 64) },
		func(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// Float64SliceVar defines a float64 slice flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) Float64SliceVar(ptr *[]float64, name string, def []float64, usage string) *SliceFlag[[]float64] {
	return f.Float64SliceVarP(ptr, name, "", def, usage)
}
