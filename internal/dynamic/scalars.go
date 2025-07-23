package dynamic

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/containeroo/tinyflags/internal/utils"
)

// String
func (g *Group) String(field string, def string, usage string) *ScalarFlag[string] {
	return registerDynamicScalar(g, field, def, utils.ParseString, utils.FormatString)
}

// Int
func (g *Group) Int(field string, def int, usage string) *ScalarFlag[int] {
	return registerDynamicScalar(g, field, def, strconv.Atoi, strconv.Itoa)
}

// Duration
func (g *Group) Duration(field string, def time.Duration, usage string) *ScalarFlag[time.Duration] {
	return registerDynamicScalar(g, field, def, time.ParseDuration, time.Duration.String)
}

// Float64
func (g *Group) Float64(field string, def float64, usage string) *ScalarFlag[float64] {
	return registerDynamicScalar(g, field, def, utils.ParseFloat64, utils.FormatFloat64)
}

// Float32
func (g *Group) Float32(field string, def float32, usage string) *ScalarFlag[float32] {
	return registerDynamicScalar(g, field, def, utils.ParseFloat32, utils.FormatFloat32)
}

// TCPAddr
func (g *Group) TCPAddr(field string, def *net.TCPAddr, usage string) *ScalarFlag[*net.TCPAddr] {
	return registerDynamicScalar(g, field, def, utils.ParseTCPAddr, utils.FormatTCPAddr)
}

// URL
func (g *Group) URL(field string, def *url.URL, usage string) *ScalarFlag[*url.URL] {
	return registerDynamicScalar(g, field, def, url.Parse, func(u *url.URL) string { return u.String() })
}

// File
func (g *Group) File(field string, def *os.File, usage string) *ScalarFlag[*os.File] {
	return registerDynamicScalar(g, field, def, utils.ParseFile, utils.FormatFile)
}

// Time
func (g *Group) Time(field string, def time.Time, usage string) *ScalarFlag[time.Time] {
	return registerDynamicScalar(g, field, def, utils.ParseTime, utils.FormatTime)
}

// Bytes
func (g *Group) Bytes(field string, def uint64, usage string) *ScalarFlag[uint64] {
	return registerDynamicScalar(g, field, def, utils.ParseBytes, utils.FormatBytes)
}
