package tinyflags

import (
	"os"
)

// FileP defines a file flag with the specified name, shorthand, default value, and usage string.
// Each input string is interpreted as a file path and opened using os.Open.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) FileP(name, short string, def *os.File, usage string) *Flag[*os.File] {
	ptr := new(*os.File)
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (*os.File, error) {
			return os.Open(s)
		},
		func(file *os.File) string {
			if file == nil {
				return ""
			}
			return file.Name()
		},
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// File defines a file flag with the specified name, default value, and usage string.
// Each input string is interpreted as a file path and opened using os.Open.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) File(name string, def *os.File, usage string) *Flag[*os.File] {
	return f.FileP(name, "", def, usage)
}

// FileVarP defines a file flag with the specified name, shorthand, default value, and usage string.
// Each input string is interpreted as a file path and opened using os.Open.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) FileVarP(ptr **os.File, name, short string, def *os.File, usage string) *Flag[*os.File] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (*os.File, error) {
			return os.Open(s)
		},
		func(file *os.File) string {
			if file == nil {
				return ""
			}
			return file.Name()
		},
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// FileVar defines a file flag with the specified name, default value, and usage string.
// Each input string is interpreted as a file path and opened using os.Open.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) FileVar(ptr **os.File, name string, def *os.File, usage string) *Flag[*os.File] {
	return f.FileVarP(ptr, name, "", def, usage)
}
