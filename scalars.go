package tinyflags

import (
	"net"
	"net/url"
	"os"
	"time"

	"github.com/containeroo/tinyflags/internal/scalar"
)

// StringVar declares a string flag, binding the flag to a pointer (ptr).
func (f *FlagSet) StringVar(ptr *string, name string, def string, usage string) *scalar.ScalarFlag[string] {
	return f.impl.StringVar(ptr, name, def, usage)
}

// String declares a string flag.
func (f *FlagSet) String(name string, def string, usage string) *scalar.ScalarFlag[string] {
	return f.StringVar(new(string), name, def, usage)
}

// IntVar declares an int flag, binding the flag to a pointer (ptr).
func (f *FlagSet) IntVar(ptr *int, name string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.impl.IntVar(ptr, name, def, usage)
}

// Int declares an int flag.
func (f *FlagSet) Int(name string, def int, usage string) *scalar.ScalarFlag[int] {
	return f.IntVar(new(int), name, def, usage)
}

// Bool defines a boolean flag.
// If Strict() is true, the flag requires an explicit value (--flag=true/false).
func (f *FlagSet) Bool(name string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVar(new(bool), name, def, usage)
}

// BoolVar defines a boolean flag, binding the flag to a pointer (ptr).
// If Strict() is true, the flag requires an explicit value (--flag=true/false).
func (f *FlagSet) BoolVar(ptr *bool, name, short string, def bool, usage string) *scalar.BoolFlag {
	return f.impl.BoolVar(ptr, name, def, usage)
}

// Counter  defines a counter flag.
// A counter is a flag that increments on each occurrence.
func (f *FlagSet) Counter(name string, def int, usage string) *scalar.CounterFlag {
	return f.impl.CounterVar(new(int), name, 0, usage)
}

// CounterVar defines a counter flag, binding the flag to a pointer (ptr).
// A counter is a flag that increments on each occurrence.
func (f *FlagSet) CounterVar(ptr *int, name string, def int, usage string) *scalar.CounterFlag {
	return f.impl.CounterVar(ptr, name, 0, usage)
}

// DurationVar declares a time.Duration flag, binding the flag to a pointer (ptr).
func (f *FlagSet) DurationVar(ptr *time.Duration, name string, def time.Duration, usage string) *scalar.ScalarFlag[time.Duration] {
	return f.impl.DurationVar(ptr, name, def, usage)
}

// Duration declares a time.Duration flag.
func (f *FlagSet) Duration(name string, def time.Duration, usage string) *scalar.ScalarFlag[time.Duration] {
	return f.DurationVar(new(time.Duration), name, def, usage)
}

// Float32Var declares a float32 flag, binding the flag to a pointer (ptr).
func (f *FlagSet) Float32Var(ptr *float32, name string, def float32, usage string) *scalar.ScalarFlag[float32] {
	return f.impl.Float32Var(ptr, name, def, usage)
}

// Float32 declares a float32 flag.
func (f *FlagSet) Float32(name string, def float32, usage string) *scalar.ScalarFlag[float32] {
	return f.Float32Var(new(float32), name, def, usage)
}

// Float64Var declares a float64 flag, binding the flag to a pointer (ptr).
func (f *FlagSet) Float64Var(ptr *float64, name string, def float64, usage string) *scalar.ScalarFlag[float64] {
	return f.impl.Float64Var(ptr, name, def, usage)
}

// Float64 declares a float64 flag.
func (f *FlagSet) Float64(name string, def float64, usage string) *scalar.ScalarFlag[float64] {
	return f.Float64Var(new(float64), name, def, usage)
}

// IPVar declares a net.IP flag, binding the flag to a pointer (ptr).
func (f *FlagSet) IPVar(ptr *net.IP, name string, def net.IP, usage string) *scalar.ScalarFlag[net.IP] {
	return f.impl.IPVar(ptr, name, def, usage)
}

// IP declares a net.IP flag.
func (f *FlagSet) IP(name string, def net.IP, usage string) *scalar.ScalarFlag[net.IP] {
	return f.IPVar(new(net.IP), name, def, usage)
}

// IPv4MaskVar declares a net.IPMask flag, binding the flag to a pointer (ptr).
func (f *FlagSet) IPv4MaskVar(ptr *net.IPMask, name string, def net.IPMask, usage string) *scalar.ScalarFlag[net.IPMask] {
	return f.impl.IPv4MaskVar(ptr, name, def, usage)
}

// IPMask declares a net.IPMask flag.
func (f *FlagSet) IPMask(name string, def net.IPMask, usage string) *scalar.ScalarFlag[net.IPMask] {
	return f.IPv4MaskVar(new(net.IPMask), name, def, usage)
}

// TCPAddrVar declares a *net.TCPAddr flag, binding the flag to a pointer (ptr).
func (f *FlagSet) TCPAddrVar(ptr **net.TCPAddr, name string, def *net.TCPAddr, usage string) *scalar.ScalarFlag[*net.TCPAddr] {
	return f.impl.TCPAddrVar(ptr, name, def, usage)
}

// TCPAddr declares a *net.TCPAddr flag.
func (f *FlagSet) TCPAddr(name string, def *net.TCPAddr, usage string) *scalar.ScalarFlag[*net.TCPAddr] {
	return f.TCPAddrVar(new(*net.TCPAddr), name, def, usage)
}

// URLVar declares a url.URL flag, binding the flag to a pointer (ptr).
func (f *FlagSet) URLVar(ptr **url.URL, name string, def *url.URL, usage string) *scalar.ScalarFlag[*url.URL] {
	return f.impl.URLVar(ptr, name, def, usage)
}

// URL declares a url.URL flag.
func (f *FlagSet) URL(name string, def *url.URL, usage string) *scalar.ScalarFlag[*url.URL] {
	return f.URLVar(new(*url.URL), name, def, usage)
}

// FileVar declares an *os.File flag, binding the flag to a pointer (ptr).
func (f *FlagSet) FileVar(ptr **os.File, name string, def *os.File, usage string) *scalar.ScalarFlag[*os.File] {
	return f.impl.FileVar(ptr, name, def, usage)
}

// File declares an *os.File flag.
func (f *FlagSet) File(name string, def *os.File, usage string) *scalar.ScalarFlag[*os.File] {
	return f.FileVar(new(*os.File), name, def, usage)
}

// TimeVar declares a time.Time flag, binding the flag to a pointer (ptr).
func (f *FlagSet) TimeVar(ptr *time.Time, name string, def time.Time, usage string) *scalar.ScalarFlag[time.Time] {
	return f.impl.TimeVar(ptr, name, def, usage)
}

// Time declares a time.Time flag.
func (f *FlagSet) Time(name string, def time.Time, usage string) *scalar.ScalarFlag[time.Time] {
	return f.TimeVar(new(time.Time), name, def, usage)
}

// BytesVar declares a uint64 “bytes” flag, binding the flag to a pointer (ptr).
func (f *FlagSet) BytesVar(ptr *uint64, name string, def uint64, usage string) *scalar.ScalarFlag[uint64] {
	return f.impl.BytesVar(ptr, name, def, usage)
}

// Bytes declares a uint64 “bytes” flag.
func (f *FlagSet) Bytes(name string, def uint64, usage string) *scalar.ScalarFlag[uint64] {
	return f.BytesVar(new(uint64), name, def, usage)
}
