package tinyflags

import "strconv"

// Float32SliceP defines a float32 slice flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Float32SliceP(name, short string, def []float32, usage string) *SliceFlag[[]float32] {
	ptr := new([]float32)
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (float32, error) {
			v, err := strconv.ParseFloat(s, 32)
			return float32(v), err
		},
		func(f float32) string {
			return strconv.FormatFloat(float64(f), 'f', -1, 32)
		},
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// Float32Slice defines a float32 slice flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Float32Slice(name string, def []float32, usage string) *SliceFlag[[]float32] {
	return f.Float32SliceP(name, "", def, usage)
}

// Float32SliceVarP defines a float32 slice flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) Float32SliceVarP(ptr *[]float32, name, short string, def []float32, usage string) *SliceFlag[[]float32] {
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (float32, error) {
			v, err := strconv.ParseFloat(s, 32)
			return float32(v), err
		},
		func(f float32) string {
			return strconv.FormatFloat(float64(f), 'f', -1, 32)
		},
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// Float32SliceVar defines a float32 slice flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) Float32SliceVar(ptr *[]float32, name string, def []float32, usage string) *SliceFlag[[]float32] {
	return f.Float32SliceVarP(ptr, name, "", def, usage)
}
