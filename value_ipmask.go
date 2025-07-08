package tinyflags

import (
	"fmt"
	"net"
	"strings"
)

// IPMaskP defines a net.IPMask flag with name, shorthand, default value, and usage string.
// The input must be in dotted decimal format (e.g., "255.255.255.0").
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IPMaskP(name, short string, def net.IPMask, usage string) *Flag[net.IPMask] {
	ptr := new(net.IPMask)

	val := NewFlagItem(
		ptr,
		def,
		func(s string) (net.IPMask, error) {
			parts := strings.Split(s, ".")
			if len(parts) != 4 {
				return nil, fmt.Errorf("invalid IP mask: %s", s)
			}
			ip := net.ParseIP(s)
			if ip == nil {
				return nil, fmt.Errorf("invalid IP format: %s", s)
			}
			return net.IPMask(ip.To4()), nil
		},
		func(m net.IPMask) string {
			if m == nil {
				return ""
			}
			return net.IP(m).String()
		},
	)

	return addScalar(f, name, short, usage, val, ptr)
}

// IPMask defines a net.IPMask flag with name, default value, and usage string.
// The input must be in dotted decimal format (e.g., "255.255.255.0").
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IPMask(name string, def net.IPMask, usage string) *Flag[net.IPMask] {
	return f.IPMaskP(name, "", def, usage)
}

// IPMaskVarP defines a net.IPMask flag with name, shorthand, default value, and usage string.
// The input must be in dotted decimal format (e.g., "255.255.255.0").
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IPMaskVarP(ptr *net.IPMask, name, short string, def net.IPMask, usage string) *Flag[net.IPMask] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (net.IPMask, error) {
			parts := strings.Split(s, ".")
			if len(parts) != 4 {
				return nil, fmt.Errorf("invalid IP mask: %s", s)
			}
			ip := net.ParseIP(s)
			if ip == nil {
				return nil, fmt.Errorf("invalid IP format: %s", s)
			}
			return net.IPMask(ip.To4()), nil
		},
		func(m net.IPMask) string {
			if m == nil {
				return ""
			}
			return net.IP(m).String()
		},
	)

	return addScalar(f, name, short, usage, val, ptr)
}

// IPMaskVar defines a net.IPMask flag with name, default value, and usage string.
// The input must be in dotted decimal format (e.g., "255.255.255.0").
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IPMaskVar(ptr *net.IPMask, name string, def net.IPMask, usage string) *Flag[net.IPMask] {
	return f.IPMaskVarP(ptr, name, "", def, usage)
}
