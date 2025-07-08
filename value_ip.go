package tinyflags

import "net"

// IPP defines a net.IP flag with the specified name, shorthand, default value, and usage.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IPP(name, short string, def net.IP, usage string) *Flag[net.IP] {
	ptr := new(net.IP)
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (net.IP, error) { return net.ParseIP(s), nil },
		func(ip net.IP) string { return ip.String() },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// IP defines a net.IP flag with the specified name, default value, and usage.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IP(name string, def net.IP, usage string) *Flag[net.IP] {
	return f.IPP(name, "", def, usage)
}

// IPVarP defines a net.IP flag with the specified name, shorthand, default value, and usage.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IPVarP(ptr *net.IP, name, short string, def net.IP, usage string) *Flag[net.IP] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (net.IP, error) { return net.ParseIP(s), nil },
		func(ip net.IP) string { return ip.String() },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// IPVar defines a net.IP flag with the specified name, default value, and usage.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IPVar(ptr *net.IP, name string, def net.IP, usage string) *Flag[net.IP] {
	return f.IPVarP(ptr, name, "", def, usage)
}
