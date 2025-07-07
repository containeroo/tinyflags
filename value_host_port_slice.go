package tinyflags

import "net"

// HostPortSliceP defines a string slice flag where each element must be a valid host:port pair,
// with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) HostPortSliceP(name, short string, def []string, usage string) *SliceFlag[[]string] {
	ptr := new([]string)
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (string, error) {
			_, _, err := net.SplitHostPort(s)
			return s, err
		},
		func(s string) string { return s },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// HostPortSlice defines a string slice flag where each element must be a valid host:port pair,
// with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) HostPortSlice(name string, def []string, usage string) *SliceFlag[[]string] {
	return f.HostPortSliceP(name, "", def, usage)
}

// HostPortSliceVarP defines a string slice flag where each element must be a valid host:port pair,
// with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) HostPortSliceVarP(ptr *[]string, name, short string, def []string, usage string) *SliceFlag[[]string] {
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (string, error) {
			_, _, err := net.SplitHostPort(s)
			return s, err
		},
		func(s string) string { return s },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// HostPortSliceVar defines a string slice flag where each element must be a valid host:port pair,
// with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) HostPortSliceVar(ptr *[]string, name string, def []string, usage string) *SliceFlag[[]string] {
	return f.HostPortSliceVarP(ptr, name, "", def, usage)
}
