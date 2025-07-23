package dynamic

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/containeroo/tinyflags/internal/utils"
)

// StringSlice
func (g *Group) StringSlice(field string, def []string, usage string) *SliceFlag[string] {
	return registerDynamicSlice(g, field, def, utils.ParseString, utils.FormatString)
}

// IntSlice
func (g *Group) IntSlice(field string, def []int, usage string) *SliceFlag[int] {
	return registerDynamicSlice(g, field, def, strconv.Atoi, strconv.Itoa)
}

// DurationSlice
func (g *Group) DurationSlice(field string, def []time.Duration, usage string) *SliceFlag[time.Duration] {
	return registerDynamicSlice(g, field, def, time.ParseDuration, time.Duration.String)
}

// Float64Slice
func (g *Group) Float64Slice(field string, def []float64, usage string) *SliceFlag[float64] {
	return registerDynamicSlice(g, field, def, utils.ParseFloat64, utils.FormatFloat64)
}

// Float32Slice
func (g *Group) Float32Slice(field string, def []float32, usage string) *SliceFlag[float32] {
	return registerDynamicSlice(g, field, def, utils.ParseFloat32, utils.FormatFloat32)
}

// TCPAddrSlice
func (g *Group) TCPAddrSlice(field string, def []*net.TCPAddr, usage string) *SliceFlag[*net.TCPAddr] {
	return registerDynamicSlice(g, field, def, utils.ParseTCPAddr, utils.FormatTCPAddr)
}

// URLSlice
func (g *Group) URLSlice(field string, def []*url.URL, usage string) *SliceFlag[*url.URL] {
	return registerDynamicSlice(g, field, def, url.Parse, func(u *url.URL) string { return u.String() })
}

// FileSlice
func (g *Group) FileSlice(field string, def []*os.File, usage string) *SliceFlag[*os.File] {
	return registerDynamicSlice(g, field, def, utils.ParseFile, utils.FormatFile)
}

// TimeSlice
func (g *Group) TimeSlice(field string, def []time.Time, usage string) *SliceFlag[time.Time] {
	return registerDynamicSlice(g, field, def, utils.ParseTime, utils.FormatTime)
}

// BytesSlice
func (g *Group) BytesSlice(field string, def []uint64, usage string) *SliceFlag[uint64] {
	return registerDynamicSlice(g, field, def, utils.ParseBytes, utils.FormatBytes)
}
