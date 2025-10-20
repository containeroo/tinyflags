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

// BoolVar defines a bool flag.
func (f *FlagSet) BoolVar(ptr *bool, name string, def bool, usage string) *scalar.BoolFlag {
	return scalar.NewBool(f, ptr, name, def, usage)
}

// Bool defines a bool flag.
func (f *FlagSet) Bool(name string, def bool, usage string) *scalar.BoolFlag {
	return scalar.NewBool(f, new(bool), name, def, usage)
}

// CounterVar defines a counter flag.
func (f *FlagSet) CounterVar(ptr *int, name string, def int, usage string) *scalar.CounterFlag {
	return scalar.NewCounter(f, ptr, name, def, usage)
}

// StringVar defines a string flag.
func (f *FlagSet) StringVar(ptr *string, name string, def string, usage string) *scalar.ScalarFlag[string] {
	return registerScalar(f, ptr, name, usage, def,
		utils.ParseString, utils.FormatString,
	)
}

// IntVar defines an int flag.
func (f *FlagSet) IntVar(ptr *int, name string, def int, usage string) *scalar.ScalarFlag[int] {
	return registerScalar(f, ptr, name, usage, def,
		strconv.Atoi, strconv.Itoa,
	)
}

// Int32Var defines an int flag.
func (f *FlagSet) Int32Var(ptr *int32, name string, def int32, usage string) *scalar.ScalarFlag[int32] {
	return registerScalar(f, ptr, name, usage, def,
		func(s string) (int32, error) {
			v, err := strconv.ParseInt(s, 10, 32)
			return int32(v), err
		},
		func(v int32) string {
			return strconv.FormatInt(int64(v), 10)
		},
	)
}

// Int64Var defines an int flag.
func (f *FlagSet) Int64Var(ptr *int64, name string, def int64, usage string) *scalar.ScalarFlag[int64] {
	return registerScalar(f, ptr, name, usage, def,
		func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 64)
		},
		func(v int64) string {
			return strconv.FormatInt(v, 10)
		},
	)
}

// DurationVar defines a time.Duration flag.
func (f *FlagSet) DurationVar(ptr *time.Duration, name string, def time.Duration, usage string) *scalar.ScalarFlag[time.Duration] {
	return registerScalar(f, ptr, name, usage, def,
		time.ParseDuration, time.Duration.String,
	)
}

// Float32Var defines a float32 flag.
func (f *FlagSet) Float32Var(ptr *float32, name string, def float32, usage string) *scalar.ScalarFlag[float32] {
	return registerScalar(f, ptr, name, usage, def,
		utils.ParseFloat32, utils.FormatFloat32,
	)
}

// Float64Var defines a float64 flag.
func (f *FlagSet) Float64Var(ptr *float64, name string, def float64, usage string) *scalar.ScalarFlag[float64] {
	return registerScalar(f, ptr, name, usage, def,
		utils.ParseFloat64, utils.FormatFloat64,
	)
}

// IPVar defines a net.IP flag.
func (f *FlagSet) IPVar(ptr *net.IP, name string, def net.IP, usage string) *scalar.ScalarFlag[net.IP] {
	return registerScalar(f, ptr, name, usage, def,
		utils.ParseIP, utils.FormatIP,
	)
}

// IPv4MaskVar defines a net.IPMask flag.
func (f *FlagSet) IPv4MaskVar(ptr *net.IPMask, name string, def net.IPMask, usage string) *scalar.ScalarFlag[net.IPMask] {
	return registerScalar(f, ptr, name, usage, def,
		utils.ParseIPv4Mask, utils.FormatIPv4Mask,
	)
}

// TCPAddrVar defines a *net.TCPAddr flag.
func (f *FlagSet) TCPAddrVar(ptr **net.TCPAddr, name string, def *net.TCPAddr, usage string) *scalar.ScalarFlag[*net.TCPAddr] {
	return registerScalar(f, ptr, name, usage, def,
		utils.ParseTCPAddr, utils.FormatTCPAddr,
	)
}

// URLVar defines a url.URL flag.
func (f *FlagSet) URLVar(ptr **url.URL, name string, def *url.URL, usage string) *scalar.ScalarFlag[*url.URL] {
	return registerScalar(f, ptr, name, usage, def,
		url.Parse, (*url.URL).String,
	)
}

// FileVar defines an *os.File flag (opened for reading).
func (f *FlagSet) FileVar(ptr **os.File, name string, def *os.File, usage string) *scalar.ScalarFlag[*os.File] {
	return registerScalar(f, ptr, name, usage, def,
		utils.ParseFile, utils.FormatFile,
	)
}

// TimeVar defines a time.Time flag (parsed as RFC3339).
func (f *FlagSet) TimeVar(ptr *time.Time, name string, def time.Time, usage string) *scalar.ScalarFlag[time.Time] {
	return registerScalar(f, ptr, name, usage, def,
		utils.ParseTime, utils.FormatTime,
	)
}

// BytesVar defines a uint64 “bytes” flag (e.g. "1GB", "512M").
func (f *FlagSet) BytesVar(ptr *uint64, name string, def uint64, usage string) *scalar.ScalarFlag[uint64] {
	return registerScalar(f, ptr, name, usage, def,
		utils.ParseBytes, utils.FormatBytes,
	)
}
