package tinyflags

import (
	"net"
	"net/url"
	"os"
	"time"

	"github.com/containeroo/tinyflags/internal/slice"
)

// StringSliceVar defines a []string flag and binds it to the given pointer.
func (f *FlagSet) StringSliceVar(ptr *[]string, name string, def []string, usage string) *slice.SliceFlag[string] {
	return f.impl.StringSliceVar(ptr, name, def, usage)
}

// StringSlice defines a []string flag and returns its handle.
func (f *FlagSet) StringSlice(name string, def []string, usage string) *slice.SliceFlag[string] {
	return f.StringSliceVar(new([]string), name, def, usage)
}

// IntSliceVar defines a []int flag and binds it to the given pointer.
func (f *FlagSet) IntSliceVar(ptr *[]int, name string, def []int, usage string) *slice.SliceFlag[int] {
	return f.impl.IntSliceVar(ptr, name, def, usage)
}

// Int32SliceVar defines a []int32 flag and binds it to the given pointer.
func (f *FlagSet) Int32SliceVar(ptr *[]int32, name string, def []int32, usage string) *slice.SliceFlag[int32] {
	return f.impl.Int32SliceVar(ptr, name, def, usage)
}

// Int64SliceVar defines a []int64 flag and binds it to the given pointer.
func (f *FlagSet) Int64SliceVar(ptr *[]int64, name string, def []int64, usage string) *slice.SliceFlag[int64] {
	return f.impl.Int64SliceVar(ptr, name, def, usage)
}

// IntSlice defines a []int flag and returns its handle.
func (f *FlagSet) IntSlice(name string, def []int, usage string) *slice.SliceFlag[int] {
	return f.IntSliceVar(new([]int), name, def, usage)
}

// Int32Slice defines a []int32 flag and returns its handle.
func (f *FlagSet) Int32Slice(name string, def []int32, usage string) *slice.SliceFlag[int32] {
	return f.Int32SliceVar(new([]int32), name, def, usage)
}

// Int64Slice defines a []int64 flag and returns its handle.
func (f *FlagSet) Int64Slice(name string, def []int64, usage string) *slice.SliceFlag[int64] {
	return f.Int64SliceVar(new([]int64), name, def, usage)
}

// DurationSliceVar defines a []time.Duration flag and binds it to the given pointer.
func (f *FlagSet) DurationSliceVar(ptr *[]time.Duration, name string, def []time.Duration, usage string) *slice.SliceFlag[time.Duration] {
	return f.impl.DurationSliceVar(ptr, name, def, usage)
}

// DurationSlice defines a []time.Duration flag and returns its handle.
func (f *FlagSet) DurationSlice(name string, def []time.Duration, usage string) *slice.SliceFlag[time.Duration] {
	return f.DurationSliceVar(new([]time.Duration), name, def, usage)
}

// Float32SliceVar defines a []float32 flag and binds it to the given pointer.
func (f *FlagSet) Float32SliceVar(ptr *[]float32, name string, def []float32, usage string) *slice.SliceFlag[float32] {
	return f.impl.Float32SliceVar(ptr, name, def, usage)
}

// Float32Slice defines a []float32 flag and returns its handle.
func (f *FlagSet) Float32Slice(name string, def []float32, usage string) *slice.SliceFlag[float32] {
	return f.Float32SliceVar(new([]float32), name, def, usage)
}

// Float64SliceVar defines a []float64 flag and binds it to the given pointer.
func (f *FlagSet) Float64SliceVar(ptr *[]float64, name string, def []float64, usage string) *slice.SliceFlag[float64] {
	return f.impl.Float64SliceVar(ptr, name, def, usage)
}

// Float64Slice defines a []float64 flag and returns its handle.
func (f *FlagSet) Float64Slice(name string, def []float64, usage string) *slice.SliceFlag[float64] {
	return f.Float64SliceVar(new([]float64), name, def, usage)
}

// IPSliceVar defines a []net.IP flag and binds it to the given pointer.
func (f *FlagSet) IPSliceVar(ptr *[]net.IP, name string, def []net.IP, usage string) *slice.SliceFlag[net.IP] {
	return f.impl.IPSliceVar(ptr, name, def, usage)
}

// IPSlice defines a []net.IP flag and returns its handle.
func (f *FlagSet) IPSlice(name string, def []net.IP, usage string) *slice.SliceFlag[net.IP] {
	return f.IPSliceVar(new([]net.IP), name, def, usage)
}

// IPMaskSliceVar defines a []net.IPMask flag and binds it to the given pointer.
func (f *FlagSet) IPMaskSliceVar(ptr *[]net.IPMask, name string, def []net.IPMask, usage string) *slice.SliceFlag[net.IPMask] {
	return f.impl.IPv4MaskSliceVar(ptr, name, def, usage)
}

// IPMaskSlice defines a []net.IPMask flag and returns its handle.
func (f *FlagSet) IPMaskSlice(name string, def []net.IPMask, usage string) *slice.SliceFlag[net.IPMask] {
	return f.IPMaskSliceVar(new([]net.IPMask), name, def, usage)
}

// TCPAddrSliceVar defines a []*net.TCPAddr flag and binds it to the given pointer.
func (f *FlagSet) TCPAddrSliceVar(ptr *[]*net.TCPAddr, name string, def []*net.TCPAddr, usage string) *slice.SliceFlag[*net.TCPAddr] {
	return f.impl.TCPAddrSliceVar(ptr, name, def, usage)
}

// TCPAddrSlice defines a []*net.TCPAddr flag and returns its handle.
func (f *FlagSet) TCPAddrSlice(name string, def []*net.TCPAddr, usage string) *slice.SliceFlag[*net.TCPAddr] {
	return f.TCPAddrSliceVar(new([]*net.TCPAddr), name, def, usage)
}

// URLSliceVar defines a []*url.URL flag and binds it to the given pointer.
func (f *FlagSet) URLSliceVar(ptr *[]*url.URL, name string, def []*url.URL, usage string) *slice.SliceFlag[*url.URL] {
	return f.impl.URLSliceVar(ptr, name, def, usage)
}

// URLSlice defines a []*url.URL flag and returns its handle.
func (f *FlagSet) URLSlice(name string, def []*url.URL, usage string) *slice.SliceFlag[*url.URL] {
	return f.URLSliceVar(new([]*url.URL), name, def, usage)
}

// FileSliceVar defines a []*os.File flag and binds it to the given pointer.
func (f *FlagSet) FileSliceVar(ptr *[]*os.File, name string, def []*os.File, usage string) *slice.SliceFlag[*os.File] {
	return f.impl.FileSliceVar(ptr, name, def, usage)
}

// FileSlice defines a []*os.File flag and returns its handle.
func (f *FlagSet) FileSlice(name string, def []*os.File, usage string) *slice.SliceFlag[*os.File] {
	return f.FileSliceVar(new([]*os.File), name, def, usage)
}

// TimeSliceVar defines a []time.Time flag and binds it to the given pointer.
func (f *FlagSet) TimeSliceVar(ptr *[]time.Time, name string, def []time.Time, usage string) *slice.SliceFlag[time.Time] {
	return f.impl.TimeSliceVar(ptr, name, def, usage)
}

// TimeSlice defines a []time.Time flag and returns its handle.
func (f *FlagSet) TimeSlice(name string, def []time.Time, usage string) *slice.SliceFlag[time.Time] {
	return f.TimeSliceVar(new([]time.Time), name, def, usage)
}

// BytesSliceVar defines a []uint64 flag that parses human-readable sizes and binds it to the given pointer.
func (f *FlagSet) BytesSliceVar(ptr *[]uint64, name string, def []uint64, usage string) *slice.SliceFlag[uint64] {
	return f.impl.BytesSliceVar(ptr, name, def, usage)
}

// BytesSlice defines a []uint64 “bytes” flag and returns its handle.
func (f *FlagSet) BytesSlice(name string, def []uint64, usage string) *slice.SliceFlag[uint64] {
	return f.BytesSliceVar(new([]uint64), name, def, usage)
}
