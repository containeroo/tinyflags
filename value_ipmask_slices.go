package tinyflags

import (
	"fmt"
	"net"
	"strings"
)

// IPMaskSliceP defines a []net.IPMask slice flag with name, shorthand, default value, and usage string.
// Each element must be in dotted decimal format (e.g., "255.255.255.0").
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IPMaskSliceP(name, short string, def []net.IPMask, usage string) *SliceFlag[[]net.IPMask] {
	ptr := new([]net.IPMask)

	val := NewSliceItem(
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
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// IPMaskSlice defines a []net.IPMask slice flag with name, default value, and usage string.
// Each element must be in dotted decimal format (e.g., "255.255.255.0").
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IPMaskSlice(name string, def []net.IPMask, usage string) *SliceFlag[[]net.IPMask] {
	return f.IPMaskSliceP(name, "", def, usage)
}

// IPMaskSliceVarP defines a []net.IPMask slice flag with name, shorthand, default value, and usage string.
// Each element must be in dotted decimal format (e.g., "255.255.255.0").
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IPMaskSliceVarP(ptr *[]net.IPMask, name, short string, def []net.IPMask, usage string) *SliceFlag[[]net.IPMask] {
	val := NewSliceItem(
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
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// IPMaskSliceVar defines a []net.IPMask slice flag with name, default value, and usage string.
// Each element must be in dotted decimal format (e.g., "255.255.255.0").
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IPMaskSliceVar(ptr *[]net.IPMask, name string, def []net.IPMask, usage string) *SliceFlag[[]net.IPMask] {
	return f.IPMaskSliceVarP(ptr, name, "", def, usage)
}
