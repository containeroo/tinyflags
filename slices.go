package tinyflags

import (
	"net"
	"net/url"
	"os"
	"time"

	"github.com/containeroo/tinyflags/internal/slice"
)

// StringSliceVar declares a []string flag, binding the flag to a pointer (ptr).
func (f *FlagSet) StringSliceVar(ptr *[]string, name string, def []string, usage string) *slice.SliceFlag[string] {
	return f.impl.StringSliceVar(ptr, name, def, usage)
}

// StringSlice declares a []string flag.
func (f *FlagSet) StringSlice(name string, def []string, usage string) *slice.SliceFlag[string] {
	return f.StringSliceVar(new([]string), name, def, usage)
}

// IntSliceVar declares a []int flag, binding the flag to a pointer (ptr).
func (f *FlagSet) IntSliceVar(ptr *[]int, name string, def []int, usage string) *slice.SliceFlag[int] {
	return f.impl.IntSliceVar(ptr, name, def, usage)
}

// IntSlice declares a []int flag.
func (f *FlagSet) IntSlice(name string, def []int, usage string) *slice.SliceFlag[int] {
	return f.IntSliceVar(new([]int), name, def, usage)
}

// DurationSliceVar declares a []time.Duration flag, binding the flag to a pointer (ptr).
func (f *FlagSet) DurationSliceVar(ptr *[]time.Duration, name string, def []time.Duration, usage string) *slice.SliceFlag[time.Duration] {
	return f.impl.DurationSliceVar(ptr, name, def, usage)
}

// DurationSlice declares a []time.Duration flag.
func (f *FlagSet) DurationSlice(name string, def []time.Duration, usage string) *slice.SliceFlag[time.Duration] {
	return f.DurationSliceVar(new([]time.Duration), name, def, usage)
}

// Float32SliceVar declares a []float32 flag, binding the flag to a pointer (ptr).
func (f *FlagSet) Float32SliceVar(ptr *[]float32, name string, def []float32, usage string) *slice.SliceFlag[float32] {
	return f.impl.Float32SliceVar(ptr, name, def, usage)
}

// Float32Slice declares a []float32 flag.
func (f *FlagSet) Float32Slice(name string, def []float32, usage string) *slice.SliceFlag[float32] {
	return f.Float32SliceVar(new([]float32), name, def, usage)
}

// Float64SliceVar declares a []float64 flag, binding the flag to a pointer (ptr).
func (f *FlagSet) Float64SliceVar(ptr *[]float64, name string, def []float64, usage string) *slice.SliceFlag[float64] {
	return f.impl.Float64SliceVar(ptr, name, def, usage)
}

// Float64Slice declares a []float64 flag.
func (f *FlagSet) Float64Slice(name string, def []float64, usage string) *slice.SliceFlag[float64] {
	return f.Float64SliceVar(new([]float64), name, def, usage)
}

// IPSliceVar declares a []net.IP flag, binding the flag to a pointer (ptr).
func (f *FlagSet) IPSliceVar(ptr *[]net.IP, name string, def []net.IP, usage string) *slice.SliceFlag[net.IP] {
	return f.impl.IPSliceVar(ptr, name, def, usage)
}

// IPSlice declares a []net.IP flag.
func (f *FlagSet) IPSlice(name string, def []net.IP, usage string) *slice.SliceFlag[net.IP] {
	return f.IPSliceVar(new([]net.IP), name, def, usage)
}

// IPMaskSliceVar declares a []net.IPMask flag, binding the flag to a pointer (ptr).
func (f *FlagSet) IPMaskSliceVar(ptr *[]net.IPMask, name string, def []net.IPMask, usage string) *slice.SliceFlag[net.IPMask] {
	return f.impl.IPv4MaskSliceVar(ptr, name, def, usage)
}

// IPMaskSlice declares a []net.IPMask flag.
func (f *FlagSet) IPMaskSlice(name string, def []net.IPMask, usage string) *slice.SliceFlag[net.IPMask] {
	return f.IPMaskSliceVar(new([]net.IPMask), name, def, usage)
}

// TCPAddrSliceVar declares a []*net.TCPAddr flag, binding the flag to a pointer (ptr).
func (f *FlagSet) TCPAddrSliceVar(ptr *[]*net.TCPAddr, name string, def []*net.TCPAddr, usage string) *slice.SliceFlag[*net.TCPAddr] {
	return f.impl.TCPAddrSliceVar(ptr, name, def, usage)
}

// TCPAddrSlice declares a []*net.TCPAddr flag.
func (f *FlagSet) TCPAddrSlice(name string, def []*net.TCPAddr, usage string) *slice.SliceFlag[*net.TCPAddr] {
	return f.TCPAddrSliceVar(new([]*net.TCPAddr), name, def, usage)
}

// URLSliceVar declares a []url.URL flag, binding the flag to a pointer (ptr).
func (f *FlagSet) URLSliceVar(ptr *[]*url.URL, name string, def []*url.URL, usage string) *slice.SliceFlag[*url.URL] {
	return f.impl.URLSliceVar(ptr, name, def, usage)
}

// URLSlice declares a []url.URL flag.
func (f *FlagSet) URLSlice(name string, def []*url.URL, usage string) *slice.SliceFlag[*url.URL] {
	return f.URLSliceVar(new([]*url.URL), name, def, usage)
}

// FileSliceVar declares a []*os.File flag, binding the flag to a pointer (ptr).
func (f *FlagSet) FileSliceVar(ptr *[]*os.File, name string, def []*os.File, usage string) *slice.SliceFlag[*os.File] {
	return f.impl.FileSliceVar(ptr, name, def, usage)
}

// FileSlice declares a []*os.File flag.
func (f *FlagSet) FileSlice(name string, def []*os.File, usage string) *slice.SliceFlag[*os.File] {
	return f.FileSliceVar(new([]*os.File), name, def, usage)
}

// TimeSliceVar declares a []time.Time flag, binding the flag to a pointer (ptr).
func (f *FlagSet) TimeSliceVar(ptr *[]time.Time, name string, def []time.Time, usage string) *slice.SliceFlag[time.Time] {
	return f.impl.TimeSliceVar(ptr, name, def, usage)
}

// TimeSlice declares a []time.Time flag.
func (f *FlagSet) TimeSlice(name string, def []time.Time, usage string) *slice.SliceFlag[time.Time] {
	return f.TimeSliceVar(new([]time.Time), name, def, usage)
}

// BytesSliceVar declares a []uint64 “bytes” flag, binding the flag to a pointer (ptr).
func (f *FlagSet) BytesSliceVar(ptr *[]uint64, name string, def []uint64, usage string) *slice.SliceFlag[uint64] {
	return f.impl.BytesSliceVar(ptr, name, def, usage)
}

// BytesSlice declares a []uint64 “bytes” flag.
func (f *FlagSet) BytesSlice(name string, def []uint64, usage string) *slice.SliceFlag[uint64] {
	return f.BytesSliceVar(new([]uint64), name, def, usage)
}
