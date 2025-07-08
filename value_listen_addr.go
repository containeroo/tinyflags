package tinyflags

import (
	"fmt"
	"net"
)

// ListenAddrP defines a string flag with the specified name, shorthand, default value, and usage string.
// The value is parsed using net.ResolveTCPAddr with "tcp" as the network.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrP(name, short, def, usage string) *Flag[string] {
	ptr := new(string)
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (string, error) {
			if _, err := net.ResolveTCPAddr("tcp", s); err != nil {
				return "", fmt.Errorf("invalid TCP address %q: %w", s, err)
			}
			return s, nil
		},
		func(s string) string { return s },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// ListenAddr defines a string flag with the specified name, default value, and usage string.
// The value is parsed using net.ResolveTCPAddr with "tcp" as the network.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddr(name string, def string, usage string) *Flag[string] {
	return f.ListenAddrP(name, "", def, usage)
}

// ListenAddrVarP defines a string flag with the specified name, shorthand, default value, and usage string.
// The value is parsed using net.ResolveTCPAddr with "tcp" as the network.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value
func (f *FlagSet) ListenAddrVarP(ptr *string, name, short string, def string, usage string) *Flag[string] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (string, error) {
			if _, err := net.ResolveTCPAddr("tcp", s); err != nil {
				return "", fmt.Errorf("invalid TCP address %q: %w", s, err)
			}
			return s, nil
		},
		func(s string) string { return s },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// ListenAddrVar defines a string flag with the specified name, default value, and usage string.
// The value is parsed using net.ResolveTCPAddr with "tcp" as the network.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value
func (f *FlagSet) ListenAddrVar(ptr *string, name string, def string, usage string) *Flag[string] {
	return f.ListenAddrVarP(ptr, name, "", def, usage)
}
