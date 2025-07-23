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
func (g *Group) String(field, usage string) *ScalarFlag[string] {
	return registerDynamicScalar(g, field, utils.ParseString, utils.FormatString)
}

// Int
func (g *Group) Int(field, usage string) *ScalarFlag[int] {
	return registerDynamicScalar(g, field, strconv.Atoi, strconv.Itoa)
}

// Duration
func (g *Group) Duration(field, usage string) *ScalarFlag[time.Duration] {
	return registerDynamicScalar(g, field, time.ParseDuration, time.Duration.String)
}

// Float64
func (g *Group) Float64(field, usage string) *ScalarFlag[float64] {
	return registerDynamicScalar(g, field, utils.ParseFloat64, utils.FormatFloat64)
}

// Float32
func (g *Group) Float32(field, usage string) *ScalarFlag[float32] {
	return registerDynamicScalar(g, field, utils.ParseFloat32, utils.FormatFloat32)
}

// TCPAddr
func (g *Group) TCPAddr(field, usage string) *ScalarFlag[*net.TCPAddr] {
	return registerDynamicScalar(g, field, utils.ParseTCPAddr, utils.FormatTCPAddr)
}

// URL
func (g *Group) URL(field, usage string) *ScalarFlag[*url.URL] {
	return registerDynamicScalar(g, field, url.Parse, func(u *url.URL) string { return u.String() })
}

// File
func (g *Group) File(field, usage string) *ScalarFlag[*os.File] {
	return registerDynamicScalar(g, field, utils.ParseFile, utils.FormatFile)
}

// Time
func (g *Group) Time(field, usage string) *ScalarFlag[time.Time] {
	return registerDynamicScalar(g, field, utils.ParseTime, utils.FormatTime)
}

// Bytes
func (g *Group) Bytes(field, usage string) *ScalarFlag[uint64] {
	return registerDynamicScalar(g, field, utils.ParseBytes, utils.FormatBytes)
}
