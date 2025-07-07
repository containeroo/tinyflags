package tinyflags

import "net"

// HostPortP defines a string flag that must represent a valid host:port pair,
// with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) HostPortP(name, short, def string, usage string) *Flag[string] {
	ptr := new(string)
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (string, error) {
			_, _, err := net.SplitHostPort(s)
			return s, err
		},
		func(s string) string { return s },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// HostPort defines a string flag that must represent a valid host:port pair,
// with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) HostPort(name, def string, usage string) *Flag[string] {
	return f.HostPortP(name, "", def, usage)
}

// HostPortVarP defines a string flag that must represent a valid host:port pair,
// with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) HostPortVarP(ptr *string, name, short string, def string, usage string) *Flag[string] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (string, error) {
			_, _, err := net.SplitHostPort(s)
			return s, err
		},
		func(s string) string { return s },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// HostPortVar defines a string flag that must represent a valid host:port pair,
// with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) HostPortVar(ptr *string, name string, def string, usage string) *Flag[string] {
	return f.HostPortVarP(ptr, name, "", def, usage)
}
