package tinyflags

import (
	"net"

	"github.com/containeroo/tinyflags/internal/scalar"
)

func (f *FlagSet) TCPAddr(name string, def *net.TCPAddr, usage string) *scalar.ScalarFlag[*net.TCPAddr] {
	ptr := new(*net.TCPAddr)
	return f.impl.TCPAddrVar(ptr, name, usage, def)
}

// If you want a “Var” that reuses an existing *net.TCPAddr variable:
func (f *FlagSet) TCPAddrVar(ptr *net.TCPAddr, name string, def *net.TCPAddr, usage string) *scalar.ScalarFlag[*net.TCPAddr] {
	wrapper := &ptr
	return f.impl.TCPAddrVar(wrapper, name, usage, def)
}
