package tinyflags

import "net"

// ListenAddrP defines a *net.TCPAddr flag with the specified name, shorthand, default value, and usage string.
// The value is parsed using net.ResolveTCPAddr with "tcp" as the network.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddrP(name, short string, def *net.TCPAddr, usage string) *Flag[*net.TCPAddr] {
	ptr := new(*net.TCPAddr)
	val := NewFlagItem(
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
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// ListenAddr defines a *net.TCPAddr flag with the specified name, default value, and usage string.
// The value is parsed using net.ResolveTCPAddr with "tcp" as the network.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) ListenAddr(name string, def *net.TCPAddr, usage string) *Flag[*net.TCPAddr] {
	return f.ListenAddrP(name, "", def, usage)
}

// ListenAddrVarP defines a *net.TCPAddr flag with the specified name, shorthand, default value, and usage string.
// The value is parsed using net.ResolveTCPAddr with "tcp" as the network.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value
func (f *FlagSet) ListenAddrVarP(ptr **net.TCPAddr, name, short string, def *net.TCPAddr, usage string) *Flag[*net.TCPAddr] {
	val := NewFlagItem(
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
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// ListenAddrVar defines a *net.TCPAddr flag with the specified name, default value, and usage string.
// The value is parsed using net.ResolveTCPAddr with "tcp" as the network.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value
func (f *FlagSet) ListenAddrVar(ptr **net.TCPAddr, name string, def *net.TCPAddr, usage string) *Flag[*net.TCPAddr] {
	return f.ListenAddrVarP(ptr, name, "", def, usage)
}
