package tinyflags

import (
	"strconv"
)

// BoolValue embeds FlagBase[bool] and tracks whether strict boolean parsing is required.
type BoolValue struct {
	*FlagBase[bool]
	strict bool
}

// IsStrictBool reports whether the flag requires an explicit value (--flag=true/false).
func (b *BoolValue) IsStrictBool() bool {
	return b.strict
}

// BoolFlag provides fluent builder methods for boolean flags,
// including support for .Strict() to require explicit values.
type BoolFlag struct {
	*Flag[bool] // embed core builder methods like Env(), Required(), etc.
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

// Bool defines a boolean flag with a default value.
func (fs *FlagSet) BoolP(name string, short string, def bool, usage string) *BoolFlag {
	return fs.BoolVarP(new(bool), name, short, def, usage)
}

func (fs *FlagSet) BoolVar(ptr *bool, name string, def bool, usage string) *BoolFlag {
	return fs.BoolVarP(ptr, name, "", def, usage)
}

// BoolVar defines a boolean flag and stores the result in the given pointer.
func (fs *FlagSet) BoolVarP(ptr *bool, name string, short string, def bool, usage string) *BoolFlag {
	val := &BoolValue{
		FlagBase: NewFlagBase(
			ptr,
			def,
			strconv.ParseBool,
			strconv.FormatBool,
		),
	}
	flag := register(fs, name, short, usage, val.FlagBase, ptr)
	return &BoolFlag{Flag: flag, val: val}
}
