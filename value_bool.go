package tinyflags

import "strconv"

// BoolValue wraps BaseValue[bool] and implements StrictBool.
type BoolValue struct {
	*FlagItem[bool]
	strict bool
}

// IsStrictBool reports true if the flag must be passed explicitly as true/false.
func (b *BoolValue) IsStrictBool() bool { return b.strict }

// boolFlag extends Flag[bool] with .Strict() support.
type boolFlag struct {
	*Flag[bool]
	val *BoolValue
}

// Strict marks the bool flag as requiring an explicit value.
func (b *boolFlag) Strict() *boolFlag {
	b.val.strict = true
	return b
}

// BoolP defines a bool flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) BoolP(name, short string, def bool, usage string) *boolFlag {
	ptr := new(bool)
	val := &BoolValue{
		FlagItem: NewFlagItem(ptr, def, strconv.ParseBool, strconv.FormatBool),
	}
	builder := addScalar(f, name, short, usage, val, ptr)
	return &boolFlag{Flag: builder, val: val}
}

// Bool defines a bool flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Bool(name string, def bool, usage string) *boolFlag {
	return f.BoolP(name, "", def, usage)
}

// BoolVarP defines a bool flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) BoolVarP(ptr *bool, name, short string, def bool, usage string) *boolFlag {
	val := &BoolValue{
		FlagItem: NewFlagItem(ptr, def, strconv.ParseBool, strconv.FormatBool),
	}
	builder := addScalar(f, name, short, usage, val, ptr)
	return &boolFlag{Flag: builder, val: val}
}

// BoolVar defines a bool flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) BoolVar(ptr *bool, name string, def bool, usage string) *boolFlag {
	return f.BoolVarP(ptr, name, "", def, usage)
}
