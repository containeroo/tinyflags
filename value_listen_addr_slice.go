package tinyflags

import "net"

// ListenAddrSliceP defines a []*net.TCPAddr slice flag with a name, shorthand, default value, and usage string.
// Each element is parsed using net.ResolveTCPAddr with "tcp" as the network.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrSliceP(name, short string, def []*net.TCPAddr, usage string) *SliceFlag[[]*net.TCPAddr] {
	var ptr *[]*net.TCPAddr
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (*net.TCPAddr, error) {
			return net.ResolveTCPAddr("tcp", s)
		},
		func(addr *net.TCPAddr) string {
			if addr == nil {
				return ""
			}
			return addr.String()
		},
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// ListenAddrSlice defines a []*net.TCPAddr slice flag without a shorthand.
// Each element is parsed using net.ResolveTCPAddr with "tcp" as the network.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrSlice(name string, def []*net.TCPAddr, usage string) *SliceFlag[[]*net.TCPAddr] {
	return f.ListenAddrSliceP(name, "", def, usage)
}

// ListenAddrSliceVarP defines a []*net.TCPAddr slice flag with a name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrSliceVarP(ptr *[]*net.TCPAddr, name, short string, def []*net.TCPAddr, usage string) *SliceFlag[[]*net.TCPAddr] {
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (*net.TCPAddr, error) {
			return net.ResolveTCPAddr("tcp", s)
		},
		func(addr *net.TCPAddr) string {
			if addr == nil {
				return ""
			}
			return addr.String()
		},
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// ListenAddrSliceVar defines a []*net.TCPAddr slice flag without a shorthand.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrSliceVar(ptr *[]*net.TCPAddr, name string, def []*net.TCPAddr, usage string) *SliceFlag[[]*net.TCPAddr] {
	return f.ListenAddrSliceVarP(ptr, name, "", def, usage)
}
