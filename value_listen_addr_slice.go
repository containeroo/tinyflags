package tinyflags

import (
	"fmt"
	"net"
)

// ListenAddrSliceP defines a []string slice flag with a name, shorthand, default value, and usage string.
// Each element is parsed using net.ResolveTCPAddr with "tcp" as the network.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrSliceP(name, short string, def []string, usage string) *SliceFlag[[]string] {
	ptr := new([]string)
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (string, error) {
			if _, err := net.ResolveTCPAddr("tcp", s); err != nil {
				return "", fmt.Errorf("invalid TCP address %q: %w", s, err)
			}
			return s, nil
		},
		func(s string) string { return s },
		f.defaultDelimiter,
	)

	return addSlice(f, name, short, usage, val, ptr)
}

// ListenAddrSlice defines a []string slice flag without a shorthand.
// Each element is parsed using net.ResolveTCPAddr with "tcp" as the network.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrSlice(name string, def []string, usage string) *SliceFlag[[]string] {
	return f.ListenAddrSliceP(name, "", def, usage)
}

// ListenAddrSliceVarP defines a []string slice flag with a name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrSliceVarP(ptr *[]string, name, short string, def []string, usage string) *SliceFlag[[]string] {
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (string, error) {
			if _, err := net.ResolveTCPAddr("tcp", s); err != nil {
				return "", fmt.Errorf("invalid TCP address %q: %w", s, err)
			}
			return s, nil
		},
		func(s string) string { return s },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// ListenAddrSliceVar defines a []string slice flag without a shorthand.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrSliceVar(ptr *[]string, name string, def []string, usage string) *SliceFlag[[]string] {
	return f.ListenAddrSliceVarP(ptr, name, "", def, usage)
}
