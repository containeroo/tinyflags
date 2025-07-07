package tinyflags

import "strconv"

// Float32P defines a float32 flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Float32P(name, short string, def float32, usage string) *Flag[float32] {
	ptr := new(float32)
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (float32, error) {
			v, err := strconv.ParseFloat(s, 32)
			return float32(v), err
		},
		func(f float32) string {
			return strconv.FormatFloat(float64(f), 'f', -1, 32)
		},
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// Float32 defines a float32 flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Float32(name string, def float32, usage string) *Flag[float32] {
	return f.Float32P(name, "", def, usage)
}

// Float32VarP defines a float32 flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) Float32VarP(ptr *float32, name, short string, def float32, usage string) *Flag[float32] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (float32, error) {
			v, err := strconv.ParseFloat(s, 32)
			return float32(v), err
		},
		func(f float32) string {
			return strconv.FormatFloat(float64(f), 'f', -1, 32)
		},
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// Float32Var defines a float32 flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) Float32Var(ptr *float32, name string, def float32, usage string) *Flag[float32] {
	return f.Float32VarP(ptr, name, "", def, usage)
}
