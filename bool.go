package tinyflags

import (
	"strconv"
)

// BoolValue holds the internal state of a boolean flag and whether it is strict.
type BoolValue struct {
	*ValueImpl[bool]
	strict bool
}

// IsStrictBool reports whether the flag requires an explicit value (--flag=true/false).
func (b *BoolValue) IsStrictBool() bool {
	return b.strict
}

// BoolFlag provides fluent builder methods for boolean flags,
// including support for .Strict() to require explicit values.
type BoolFlag struct {
	*Flag[bool] // embeds core builder methods like Env(), Required(), etc.
	val         *BoolValue
}

// Strict marks this boolean flag as requiring an explicit value.
func (b *BoolFlag) Strict() *BoolFlag {
	b.val.strict = true
	return b
}

// Bool defines a boolean flag with a default value.
func (fs *FlagSet) Bool(name string, def bool, usage string) *BoolFlag {
	return fs.BoolVarP(new(bool), name, "", def, usage)
}

// BoolP defines a boolean flag with a short name and a default value.
func (fs *FlagSet) BoolP(name, short string, def bool, usage string) *BoolFlag {
	return fs.BoolVarP(new(bool), name, short, def, usage)
}

// BoolVar defines a boolean flag and binds it to the given pointer.
func (fs *FlagSet) BoolVar(ptr *bool, name string, def bool, usage string) *BoolFlag {
	return fs.BoolVarP(ptr, name, "", def, usage)
}

// BoolVarP defines a boolean flag with a short name and binds it to the given pointer.
func (fs *FlagSet) BoolVarP(ptr *bool, name, short string, def bool, usage string) *BoolFlag {
	val := &BoolValue{
		ValueImpl: NewValueImpl(
			ptr,
			def,
			strconv.ParseBool,
			strconv.FormatBool,
		),
	}
	flag := addScalar(fs, name, short, usage, val, ptr)
	return &BoolFlag{Flag: flag, val: val}
}
