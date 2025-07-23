package engine

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/containeroo/tinyflags/internal/scalar"
	"github.com/containeroo/tinyflags/internal/utils"
)

// StringVar defines a string flag.
func (f *FlagSet) StringVar(ptr *string, name string, def string, usage string) *scalar.ScalarFlag[string] {
	return defineScalar(f, ptr, name, usage, def,
		utils.ParseString, utils.FormatString,
	)
}

// IntVar defines an int flag.
func (f *FlagSet) IntVar(ptr *int, name string, def int, usage string) *scalar.ScalarFlag[int] {
	return defineScalar(f, ptr, name, usage, def,
		strconv.Atoi, strconv.Itoa,
	)
}

// DurationVar defines a time.Duration flag.
func (f *FlagSet) DurationVar(ptr *time.Duration, name string, def time.Duration, usage string) *scalar.ScalarFlag[time.Duration] {
	return defineScalar(f, ptr, name, usage, def,
		time.ParseDuration, time.Duration.String,
	)
}

// Float32Var defines a float32 flag.
func (f *FlagSet) Float32Var(ptr *float32, name string, def float32, usage string) *scalar.ScalarFlag[float32] {
	return defineScalar(f, ptr, name, usage, def,
		utils.ParseFloat32, utils.FormatFloat32,
	)
}

// Float64Var defines a float64 flag.
func (f *FlagSet) Float64Var(ptr *float64, name string, def float64, usage string) *scalar.ScalarFlag[float64] {
	return defineScalar(f, ptr, name, usage, def,
		utils.ParseFloat64, utils.FormatFloat64,
	)
}

// IPVar defines a net.IP flag.
func (f *FlagSet) IPVar(ptr *net.IP, name string, def net.IP, usage string) *scalar.ScalarFlag[net.IP] {
	return defineScalar(f, ptr, name, usage, def,
		utils.ParseIP, utils.FormatIP,
	)
}

// IPv4MaskVar defines a net.IPMask flag.
func (f *FlagSet) IPv4MaskVar(ptr *net.IPMask, name string, def net.IPMask, usage string) *scalar.ScalarFlag[net.IPMask] {
	return defineScalar(f, ptr, name, usage, def,
		utils.ParseIPv4Mask, utils.FormatIPv4Mask,
	)
}

// TCPAddrVar defines a *net.TCPAddr flag.
func (f *FlagSet) TCPAddrVar(ptr **net.TCPAddr, name string, def *net.TCPAddr, usage string) *scalar.ScalarFlag[*net.TCPAddr] {
	return defineScalar(f, ptr, name, usage, def,
		utils.ParseTCPAddr, utils.FormatTCPAddr,
	)
}

// URLVar defines a url.URL flag.
func (f *FlagSet) URLVar(ptr **url.URL, name string, def *url.URL, usage string) *scalar.ScalarFlag[*url.URL] {
	return defineScalar(f, ptr, name, usage, def,
		url.Parse, (*url.URL).String,
	)
}

// FileVar defines an *os.File flag (opened for reading).
func (f *FlagSet) FileVar(ptr **os.File, name string, def *os.File, usage string) *scalar.ScalarFlag[*os.File] {
	return defineScalar(f, ptr, name, usage, def,
		utils.ParseFile, utils.FormatFile,
	)
}

// TimeVar defines a time.Time flag (parsed as RFC3339).
func (f *FlagSet) TimeVar(ptr *time.Time, name string, def time.Time, usage string) *scalar.ScalarFlag[time.Time] {
	return defineScalar(f, ptr, name, usage, def,
		utils.ParseTime, utils.FormatTime,
	)
}

// BytesVar defines a uint64 “bytes” flag (e.g. "1GB", "512M").
func (f *FlagSet) BytesVar(ptr *uint64, name string, def uint64, usage string) *scalar.ScalarFlag[uint64] {
	return defineScalar(f, ptr, name, usage, def,
		utils.ParseBytes, utils.FormatBytes,
	)
}
