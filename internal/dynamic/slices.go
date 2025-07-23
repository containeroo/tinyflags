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
func (g *Group) StringSlice(field, usage string) *SliceFlag[string] {
	return defineDynamicSlice(g, field, utils.ParseString, utils.FormatString)
}

// IntSlice
func (g *Group) IntSlice(field, usage string) *SliceFlag[int] {
	return defineDynamicSlice(g, field, strconv.Atoi, strconv.Itoa)
}

// DurationSlice
func (g *Group) DurationSlice(field, usage string) *SliceFlag[time.Duration] {
	return defineDynamicSlice(g, field, time.ParseDuration, time.Duration.String)
}

// Float64Slice
func (g *Group) Float64Slice(field, usage string) *SliceFlag[float64] {
	return defineDynamicSlice(g, field, utils.ParseFloat64, utils.FormatFloat64)
}

// Float32Slice
func (g *Group) Float32Slice(field, usage string) *SliceFlag[float32] {
	return defineDynamicSlice(g, field, utils.ParseFloat32, utils.FormatFloat32)
}

// TCPAddrSlice
func (g *Group) TCPAddrSlice(field, usage string) *SliceFlag[*net.TCPAddr] {
	return defineDynamicSlice(g, field, utils.ParseTCPAddr, utils.FormatTCPAddr)
}

// URLSlice
func (g *Group) URLSlice(field, usage string) *SliceFlag[*url.URL] {
	return defineDynamicSlice(g, field, url.Parse, func(u *url.URL) string { return u.String() })
}

// FileSlice
func (g *Group) FileSlice(field, usage string) *SliceFlag[*os.File] {
	return defineDynamicSlice(g, field,
		func(s string) (*os.File, error) { return os.Open(s) },
		func(f *os.File) string { return f.Name() },
	)
}

// TimeSlice
func (g *Group) TimeSlice(field, usage string) *SliceFlag[time.Time] {
	return defineDynamicSlice(g, field,
		func(s string) (time.Time, error) { return time.Parse(time.RFC3339, s) },
		func(t time.Time) string { return t.Format(time.RFC3339) },
	)
}

// BytesSlice
func (g *Group) BytesSlice(field, usage string) *SliceFlag[uint64] {
	return defineDynamicSlice(g, field, utils.ParseBytes, utils.FormatBytes)
}
