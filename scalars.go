package tinyflags

import (
	"net"
	"net/url"
	"os"
	"time"

	"github.com/containeroo/tinyflags/internal/scalar"
)

// StringVar defines a string flag and binds it to the given pointer.
func (f *FlagSet) StringVar(ptr *string, name string, def string, usage string) *scalar.ScalarFlag[string] {
	return f.impl.StringVar(ptr, name, def, usage)
}

// String defines a string flag and returns its handle.
func (f *FlagSet) String(name string, def string, usage string) *scalar.ScalarFlag[string] {
	return f.StringVar(new(string), name, def, usage)
}

// IntVar defines an int flag and binds it to the given pointer.
func (f *FlagSet) IntVar(ptr *int, name string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.impl.IntVar(ptr, name, def, usage)
}

// Int32Var defines an int flag and binds it to the given pointer.
func (f *FlagSet) Int32Var(ptr *int32, name string, def int32, usage string) *scalar.ScalarFlag[int32] {
	return f.impl.Int32Var(ptr, name, def, usage)
}

// Int64Var defines an int flag and binds it to the given pointer.
func (f *FlagSet) Int64Var(ptr *int64, name string, def int64, usage string) *scalar.ScalarFlag[int64] {
	return f.impl.Int64Var(ptr, name, def, usage)
}

// Int defines an int flag and returns its handle.
func (f *FlagSet) Int(name string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.IntVar(new(int), name, def, usage)
}

// Int32 defines an int flag and returns its handle.
func (f *FlagSet) Int32(name string, def int32, usage string) *scalar.ScalarFlag[int32] {
	return f.Int32Var(new(int32), name, def, usage)
}

// Int64 defines an int flag and returns its handle.
func (f *FlagSet) Int64(name string, def int64, usage string) *scalar.ScalarFlag[int64] {
	return f.Int64Var(new(int64), name, def, usage)
}

// Bool defines a bool flag.
// If Strict() is enabled, it must be set explicitly (--flag=true/false).
func (f *FlagSet) Bool(name string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVar(new(bool), name, def, usage)
}

// BoolVar defines a bool flag and binds it to the given pointer.
// If Strict() is enabled, it must be set explicitly (--flag=true/false).
func (f *FlagSet) BoolVar(ptr *bool, name string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVar(ptr, name, def, usage)
}

// Counter defines a counter flag that increments with each occurrence.
func (f *FlagSet) Counter(name string, def int, usage string) *scalar.CounterFlag {
	return f.impl.CounterVar(new(int), name, 0, usage)
}

// CounterVar defines a counter flag and binds it to the given pointer.
func (f *FlagSet) CounterVar(ptr *int, name string, def int, usage string) *scalar.CounterFlag {
	return f.impl.CounterVar(ptr, name, 0, usage)
}

// DurationVar defines a time.Duration flag and binds it to the given pointer.
func (f *FlagSet) DurationVar(ptr *time.Duration, name string, def time.Duration, usage string) *scalar.ScalarFlag[time.Duration] {
	return f.impl.DurationVar(ptr, name, def, usage)
}

// Duration defines a time.Duration flag and returns its handle.
func (f *FlagSet) Duration(name string, def time.Duration, usage string) *scalar.ScalarFlag[time.Duration] {
	return f.DurationVar(new(time.Duration), name, def, usage)
}

// Float32Var defines a float32 flag and binds it to the given pointer.
func (f *FlagSet) Float32Var(ptr *float32, name string, def float32, usage string) *scalar.ScalarFlag[float32] {
	return f.impl.Float32Var(ptr, name, def, usage)
}

// Float32 defines a float32 flag and returns its handle.
func (f *FlagSet) Float32(name string, def float32, usage string) *scalar.ScalarFlag[float32] {
	return f.Float32Var(new(float32), name, def, usage)
}

// Float64Var defines a float64 flag and binds it to the given pointer.
func (f *FlagSet) Float64Var(ptr *float64, name string, def float64, usage string) *scalar.ScalarFlag[float64] {
	return f.impl.Float64Var(ptr, name, def, usage)
}

// Float64 defines a float64 flag and returns its handle.
func (f *FlagSet) Float64(name string, def float64, usage string) *scalar.ScalarFlag[float64] {
	return f.Float64Var(new(float64), name, def, usage)
}

// IPVar defines a net.IP flag and binds it to the given pointer.
func (f *FlagSet) IPVar(ptr *net.IP, name string, def net.IP, usage string) *scalar.ScalarFlag[net.IP] {
	return f.impl.IPVar(ptr, name, def, usage)
}

// IP defines a net.IP flag and returns its handle.
func (f *FlagSet) IP(name string, def net.IP, usage string) *scalar.ScalarFlag[net.IP] {
	return f.IPVar(new(net.IP), name, def, usage)
}

// IPv4MaskVar defines a net.IPMask flag and binds it to the given pointer.
func (f *FlagSet) IPv4MaskVar(ptr *net.IPMask, name string, def net.IPMask, usage string) *scalar.ScalarFlag[net.IPMask] {
	return f.impl.IPv4MaskVar(ptr, name, def, usage)
}

// IPMask defines a net.IPMask flag and returns its handle.
func (f *FlagSet) IPMask(name string, def net.IPMask, usage string) *scalar.ScalarFlag[net.IPMask] {
	return f.IPv4MaskVar(new(net.IPMask), name, def, usage)
}

// TCPAddrVar defines a *net.TCPAddr flag and binds it to the given pointer.
func (f *FlagSet) TCPAddrVar(ptr **net.TCPAddr, name string, def *net.TCPAddr, usage string) *scalar.ScalarFlag[*net.TCPAddr] {
	return f.impl.TCPAddrVar(ptr, name, def, usage)
}

// TCPAddr defines a *net.TCPAddr flag and returns its handle.
func (f *FlagSet) TCPAddr(name string, def *net.TCPAddr, usage string) *scalar.ScalarFlag[*net.TCPAddr] {
	return f.TCPAddrVar(new(*net.TCPAddr), name, def, usage)
}

// URLVar defines a *url.URL flag and binds it to the given pointer.
func (f *FlagSet) URLVar(ptr **url.URL, name string, def *url.URL, usage string) *scalar.ScalarFlag[*url.URL] {
	return f.impl.URLVar(ptr, name, def, usage)
}

// URL defines a *url.URL flag and returns its handle.
func (f *FlagSet) URL(name string, def *url.URL, usage string) *scalar.ScalarFlag[*url.URL] {
	return f.URLVar(new(*url.URL), name, def, usage)
}

// FileVar defines an *os.File flag and binds it to the given pointer.
func (f *FlagSet) FileVar(ptr **os.File, name string, def *os.File, usage string) *scalar.ScalarFlag[*os.File] {
	return f.impl.FileVar(ptr, name, def, usage)
}

// File defines an *os.File flag and returns its handle.
func (f *FlagSet) File(name string, def *os.File, usage string) *scalar.ScalarFlag[*os.File] {
	return f.FileVar(new(*os.File), name, def, usage)
}

// TimeVar defines a time.Time flag and binds it to the given pointer.
func (f *FlagSet) TimeVar(ptr *time.Time, name string, def time.Time, usage string) *scalar.ScalarFlag[time.Time] {
	return f.impl.TimeVar(ptr, name, def, usage)
}

// Time defines a time.Time flag and returns its handle.
func (f *FlagSet) Time(name string, def time.Time, usage string) *scalar.ScalarFlag[time.Time] {
	return f.TimeVar(new(time.Time), name, def, usage)
}

// BytesVar defines a uint64 flag with byte parsing and binds it to the given pointer.
func (f *FlagSet) BytesVar(ptr *uint64, name string, def uint64, usage string) *scalar.ScalarFlag[uint64] {
	return f.impl.BytesVar(ptr, name, def, usage)
}

// Bytes defines a uint64 flag that parses sizes (e.g., "10MB") and returns its handle.
func (f *FlagSet) Bytes(name string, def uint64, usage string) *scalar.ScalarFlag[uint64] {
	return f.BytesVar(new(uint64), name, def, usage)
}
