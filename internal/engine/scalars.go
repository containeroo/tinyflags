package engine

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/containeroo/tinyflags/internal/scalar"
)

func (f *FlagSet) StringVar(ptr *string, name string, def string, usage string) *scalar.ScalarFlag[string] {
	return defineScalar(f, ptr, name, usage, def,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)
}

func (f *FlagSet) DurationVar(ptr *time.Duration, name string, def time.Duration, usage string) *scalar.ScalarFlag[time.Duration] {
	return defineScalar(f, ptr, name, usage, def,
		time.ParseDuration,
		time.Duration.String,
	)
}

func (f *FlagSet) IntVar(ptr *int, name string, def int, usage string) *scalar.ScalarFlag[int] {
	return defineScalar(f, ptr, name, usage, def,
		strconv.Atoi,
		strconv.Itoa,
	)
}

func (f *FlagSet) IPVar(ptr *net.IP, name string, def net.IP, usage string) *scalar.ScalarFlag[net.IP] {
	return defineScalar(f, ptr, name, usage, def,
		func(s string) (net.IP, error) { return net.ParseIP(s), nil },
		func(ip net.IP) string { return ip.String() },
	)
}

func (f *FlagSet) IPv4MaskVar(ptr *net.IPMask, name string, def net.IPMask, usage string) *scalar.ScalarFlag[net.IPMask] {
	return defineScalar(f, ptr, name, usage, def,
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
		func(ip net.IPMask) string { return ip.String() },
	)
}

// TCPAddrVar defines a flag whose type is *net.TCPAddr.  You
// *net.TCPAddr, and the usual name/short/usage.
// engine/scalars.go
func (f *FlagSet) TCPAddrVar(
	ptr **net.TCPAddr,
	name, usage string,
	def *net.TCPAddr,
) *scalar.ScalarFlag[*net.TCPAddr] {
	return defineScalar(f, ptr, name, usage, def,
		// parse string → *net.TCPAddr
		func(s string) (*net.TCPAddr, error) {
			addr, err := net.ResolveTCPAddr("tcp", s)
			if err != nil {
				return nil, fmt.Errorf("invalid TCP address %q: %w", s, err)
			}
			return addr, nil
		},
		// format *net.TCPAddr → string
		func(addr *net.TCPAddr) string {
			if addr == nil {
				return ""
			}
			return addr.String()
		},
	)
}

func (f *FlagSet) Float32Var(ptr *float32, name string, def float32, usage string) *scalar.ScalarFlag[float32] {
	return defineScalar(f, ptr, name, usage, def,
		func(s string) (float32, error) {
			v, err := strconv.ParseFloat(s, 32)
			return float32(v), err
		},
		func(f float32) string { return strconv.FormatFloat(float64(f), 'f', -1, 32) },
	)
}

func (f *FlagSet) Float64Var(ptr *float64, name string, def float64, usage string) *scalar.ScalarFlag[float64] {
	return defineScalar(f, ptr, name, usage, def,
		func(s string) (float64, error) { return strconv.ParseFloat(s, 64) },
		func(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) },
	)
}
