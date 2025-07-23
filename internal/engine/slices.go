package engine

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/containeroo/tinyflags/internal/slice"
	"github.com/containeroo/tinyflags/internal/utils"
)

// StringVar defines a string flag.
func (f *FlagSet) StringSliceVar(ptr *[]string, name string, def []string, usage string) *slice.SliceFlag[string] {
	return defineSlice(f, ptr, name, usage, def,
		utils.ParseString, utils.FormatString,
	)
}

// IntVar defines an int flag.
func (f *FlagSet) IntSliceVar(ptr *[]int, name string, def []int, usage string) *slice.SliceFlag[int] {
	return defineSlice(f, ptr, name, usage, def,
		strconv.Atoi, strconv.Itoa,
	)
}

// DurationVar defines a time.Duration flag.
func (f *FlagSet) DurationSliceVar(ptr *[]time.Duration, name string, def []time.Duration, usage string) *slice.SliceFlag[time.Duration] {
	return defineSlice(f, ptr, name, usage, def,
		time.ParseDuration, time.Duration.String,
	)
}

// Float32Var defines a float32 flag.
func (f *FlagSet) Float32SliceVar(ptr *[]float32, name string, def []float32, usage string) *slice.SliceFlag[float32] {
	return defineSlice(f, ptr, name, usage, def,
		utils.ParseFloat32, utils.FormatFloat32,
	)
}

// Float64Var defines a float64 flag.
func (f *FlagSet) Float64SliceVar(ptr *[]float64, name string, def []float64, usage string) *slice.SliceFlag[float64] {
	return defineSlice(f, ptr, name, usage, def,
		utils.ParseFloat64, utils.FormatFloat64,
	)
}

// IPVar defines a net.IP flag.
func (f *FlagSet) IPSliceVar(ptr *[]net.IP, name string, def []net.IP, usage string) *slice.SliceFlag[net.IP] {
	return defineSlice(f, ptr, name, usage, def,
		utils.ParseIP, utils.FormatIP,
	)
}

// IPv4MaskVar defines a net.IPMask flag.
func (f *FlagSet) IPv4MaskSliceVar(ptr *[]net.IPMask, name string, def []net.IPMask, usage string) *slice.SliceFlag[net.IPMask] {
	return defineSlice(f, ptr, name, usage, def,
		utils.ParseIPv4Mask, utils.FormatIPv4Mask,
	)
}

// TCPAddrVar defines a *net.TCPAddr flag.
func (f *FlagSet) TCPAddrSliceVar(ptr *[]*net.TCPAddr, name string, def []*net.TCPAddr, usage string) *slice.SliceFlag[*net.TCPAddr] {
	return defineSlice(f, ptr, name, usage, def,
		utils.ParseTCPAddr, utils.FormatTCPAddr,
	)
}

// URLVar defines a url.URL flag.
func (f *FlagSet) URLSliceVar(ptr *[]*url.URL, name string, def []*url.URL, usage string) *slice.SliceFlag[*url.URL] {
	return defineSlice(f, ptr, name, usage, def,
		url.Parse, (*url.URL).String,
	)
}

// FileVar defines an *os.File flag (opened for reading).
func (f *FlagSet) FileSliceVar(ptr *[]*os.File, name string, def []*os.File, usage string) *slice.SliceFlag[*os.File] {
	return defineSlice(f, ptr, name, usage, def,
		utils.ParseFile, utils.FormatFile,
	)
}

// TimeVar defines a time.Time flag (parsed as RFC3339).
func (f *FlagSet) TimeSliceVar(ptr *[]time.Time, name string, def []time.Time, usage string) *slice.SliceFlag[time.Time] {
	return defineSlice(f, ptr, name, usage, def,
		utils.ParseTime, utils.FormatTime,
	)
}

// BytesVar defines a uint64 “bytes” flag (e.g. "1GB", "512M").
func (f *FlagSet) BytesSliceVar(ptr *[]uint64, name string, def []uint64, usage string) *slice.SliceFlag[uint64] {
	return defineSlice(f, ptr, name, usage, def,
		utils.ParseBytes, utils.FormatBytes,
	)
}
