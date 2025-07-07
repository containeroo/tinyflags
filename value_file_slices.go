package tinyflags

import "os"

// FileSliceP defines a file slice flag with the specified name, shorthand, default value, and usage string.
// Each input string is interpreted as a file path and opened using os.Open.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) FileSliceP(name, short string, def []*os.File, usage string) *SliceFlag[[]*os.File] {
	ptr := new([]*os.File)
	val := NewSliceItem(
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
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// FileSlice defines a file slice flag with the specified name, default value, and usage string.
// Each input string is interpreted as a file path and opened using os.Open.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) FileSlice(name string, def []*os.File, usage string) *SliceFlag[[]*os.File] {
	return f.FileSliceP(name, "", def, usage)
}

// FileSliceVarP defines a file slice flag with the specified name, shorthand, default value, and usage string.
// Each input string is interpreted as a file path and opened using os.Open.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) FileSliceVarP(ptr *[]*os.File, name, short string, def []*os.File, usage string) *SliceFlag[[]*os.File] {
	val := NewSliceItem(
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
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// FileSliceVar defines a file slice flag with the specified name, default value, and usage string.
// Each input string is interpreted as a file path and opened using os.Open.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) FileSliceVar(ptr *[]*os.File, name string, def []*os.File, usage string) *SliceFlag[[]*os.File] {
	return f.FileSliceVarP(ptr, name, "", def, usage)
}
