package dynamic

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/containeroo/tinyflags/internal/utils"
)

// Bool
func (g *Group) StrictBool(field string, def bool, usage string) *BoolFlag {
	return registerDynamicBool(g, field, def, usage, strconv.ParseBool, strconv.FormatBool)
}

// Bool registers a dynamic bool flag.
func (g *Group) Bool(field string, def bool, usage string) *BoolFlag {
	return registerDynamicBool(g, field, def, usage, strconv.ParseBool, strconv.FormatBool)
}

// String
func (g *Group) String(field string, def string, usage string) *ScalarFlag[string] {
	return registerDynamicScalar(g, field, def, usage, utils.ParseString, utils.FormatString)
}

// Enum registers a dynamic string enum flag.
func (g *Group) Enum(field string, def string, usage string, allowed ...string) *ScalarFlag[string] {
	return g.String(field, def, usage).Choices(allowed...)
}

// Enum registers a typed dynamic string enum flag.
func Enum[T enumValue](g *Group, field string, def T, usage string, allowed ...T) *ScalarFlag[T] {
	return registerDynamicScalar(
		g,
		field,
		def,
		usage,
		parseEnumValue[T],
		formatEnumValue[T],
	).Choices(allowed...)
}

type enumValue interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func parseEnumValue[T enumValue](raw string) (T, error) {
	var zero T
	typ := reflect.TypeOf(zero)
	out := reflect.New(typ).Elem()

	switch typ.Kind() {
	case reflect.String:
		out.SetString(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(raw, 10, typ.Bits())
		if err != nil {
			return zero, err
		}
		out.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(raw, 10, typ.Bits())
		if err != nil {
			return zero, err
		}
		out.SetUint(v)
	default:
		return zero, fmt.Errorf("unsupported enum kind %s", typ.Kind())
	}

	return out.Interface().(T), nil
}

func formatEnumValue[T enumValue](v T) string {
	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Int
func (g *Group) Int(field string, def int, usage string) *ScalarFlag[int] {
	return registerDynamicScalar(g, field, def, usage, strconv.Atoi, strconv.Itoa)
}

// Int32
func (g *Group) Int32(field string, def int32, usage string) *ScalarFlag[int32] {
	return registerDynamicScalar(g, field, def, usage,
		func(s string) (int32, error) {
			v, err := strconv.ParseInt(s, 10, 32)
			return int32(v), err
		},
		func(v int32) string {
			return strconv.FormatInt(int64(v), 10)
		},
	)
}

// Int64
func (g *Group) Int64(field string, def int64, usage string) *ScalarFlag[int64] {
	return registerDynamicScalar(g, field, def, usage,
		func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 64)
		},
		func(v int64) string {
			return strconv.FormatInt(v, 10)
		},
	)
}

// Duration
func (g *Group) Duration(field string, def time.Duration, usage string) *ScalarFlag[time.Duration] {
	return registerDynamicScalar(g, field, def, usage, time.ParseDuration, time.Duration.String)
}

// Float64
func (g *Group) Float64(field string, def float64, usage string) *ScalarFlag[float64] {
	return registerDynamicScalar(g, field, def, usage, utils.ParseFloat64, utils.FormatFloat64)
}

// Float32
func (g *Group) Float32(field string, def float32, usage string) *ScalarFlag[float32] {
	return registerDynamicScalar(g, field, def, usage, utils.ParseFloat32, utils.FormatFloat32)
}

// TCPAddr
func (g *Group) TCPAddr(field string, def *net.TCPAddr, usage string) *ScalarFlag[*net.TCPAddr] {
	return registerDynamicScalar(g, field, def, usage, utils.ParseTCPAddr, utils.FormatTCPAddr)
}

// URL
func (g *Group) URL(field string, def *url.URL, usage string) *ScalarFlag[*url.URL] {
	return registerDynamicScalar(g, field, def, usage, url.Parse, func(u *url.URL) string { return u.String() })
}

// File
func (g *Group) File(field string, def *os.File, usage string) *ScalarFlag[*os.File] {
	return registerDynamicScalar(g, field, def, usage, utils.ParseFile, utils.FormatFile)
}

// Time
func (g *Group) Time(field string, def time.Time, usage string) *ScalarFlag[time.Time] {
	return registerDynamicScalar(g, field, def, usage, utils.ParseTime, utils.FormatTime)
}

// Bytes
func (g *Group) Bytes(field string, def uint64, usage string) *ScalarFlag[uint64] {
	return registerDynamicScalar(g, field, def, usage, utils.ParseBytes, utils.FormatBytes)
}
