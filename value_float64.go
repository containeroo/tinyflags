package tinyflags

import "strconv"

// Float64P defines a float64 flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Float64P(name, short string, def float64, usage string) *Flag[float64] {
	ptr := new(float64)
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (float64, error) { return strconv.ParseFloat(s, 64) },
		func(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// Float64 defines a float64 flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Float64(name string, def float64, usage string) *Flag[float64] {
	return f.Float64P(name, "", def, usage)
}

// Float64VarP defines a float64 flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) Float64VarP(ptr *float64, name, short string, def float64, usage string) *Flag[float64] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (float64, error) { return strconv.ParseFloat(s, 64) },
		func(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// Float64Var defines a float64 flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) Float64Var(ptr *float64, name string, def float64, usage string) *Flag[float64] {
	return f.Float64VarP(ptr, name, "", def, usage)
}
