package tinyflags

import "net"

// IPSliceP defines a []net.IP slice flag with name, shorthand, default value, and usage string.
// It returns the slice flag for chaining.
func (f *FlagSet) IPSliceP(name, short string, def []net.IP, usage string) *SliceFlag[[]net.IP] {
	ptr := new([]net.IP)
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (net.IP, error) { return net.ParseIP(s), nil },
		func(ip net.IP) string { return ip.String() },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// IPSlice defines a []net.IP slice flag with name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) IPSlice(name string, def []net.IP, usage string) *SliceFlag[[]net.IP] {
	return f.IPSliceP(name, "", def, usage)
}

// IPSliceVarP defines a []net.IP slice flag with name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IPSliceVarP(ptr *[]net.IP, name, short string, def []net.IP, usage string) *SliceFlag[[]net.IP] {
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (net.IP, error) { return net.ParseIP(s), nil },
		func(ip net.IP) string { return ip.String() },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// IPSliceVar defines a []net.IP slice flag with name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) IPSliceVar(ptr *[]net.IP, name string, def []net.IP, usage string) *SliceFlag[[]net.IP] {
	return f.IPSliceVarP(ptr, name, "", def, usage)
}
