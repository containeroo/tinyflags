package tinyflags

import (
	"net"

	"github.com/containeroo/tinyflags/internal/scalar"
)

// IP defines a scalar string flag with default value.
func (f *FlagSet) IP(name string, def net.IP, usage string) *scalar.ScalarFlag[net.IP] {
	return f.impl.IPVar(new(net.IP), name, def, usage)
}

func (f *FlagSet) IPVar(ptr *net.IP, name string, def net.IP, usage string) *scalar.ScalarFlag[net.IP] {
	return f.impl.IPVar(ptr, name, def, usage)
}
