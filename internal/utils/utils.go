package utils

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func AllowOnly[T any](format func(T) string, allowed []T) func(T) error {
	return func(v T) error {
		for _, a := range allowed {
			if format(a) == format(v) {
				return nil
			}
		}
		return fmt.Errorf("must be one of [%v]", allowed)
	}
}

func FormatList[T any](format func(T) string, values []T) []string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = format(v)
	}
	return out
}

func JoinFormatted[T any](values []T, format func(T) string) string {
	if format == nil {
		return fmt.Sprintf("%v", values) // fallback
	}
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = format(v)
	}
	return strings.Join(parts, ",")
}

// PluralSuffix returns "s" if the given number is not 1, otherwise it returns an empty string.
// It is useful for constructing basic pluralized words like "flag" or "flags".
func PluralSuffix(i int) string {
	if i != 1 {
		return "s"
	}
	return ""
}

// ParseString string → string
func ParseString(s string) (string, error) { return s, nil }

// FormatString string → string
func FormatString(s string) string { return s }

// ParseTCPAddr string → *net.TCPAddr
func ParseTCPAddr(s string) (*net.TCPAddr, error) {
	addr, err := net.ResolveTCPAddr("tcp", s)
	if err != nil {
		return nil, fmt.Errorf("invalid TCP address %q: %w", s, err)
	}
	return addr, nil
}

// FormatTCPAddr *net.TCPAddr → string
func FormatTCPAddr(addr *net.TCPAddr) string {
	if addr == nil {
		return ""
	}
	return addr.String()
}

// ParseDuration string → time.Duration
func ParseFloat64(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

// FormatFloat64 time.Duration → string
func FormatFloat64(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

// ParseFloat32 string → float32
func ParseFloat32(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

// FormatFloat32 float32 → string
func FormatFloat32(f float32) string { return strconv.FormatFloat(float64(f), 'f', -1, 32) }

// ParseIP net.IP → string
func ParseIP(s string) (net.IP, error) { return net.ParseIP(s), nil }

// FormatIP net.IP → string
func FormatIP(ip net.IP) string { return ip.String() }

// ParseIPv4Mask string → net.IPMask
func ParseIPv4Mask(s string) (net.IPMask, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid IP mask: %s", s)
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP format: %s", s)
	}
	return net.IPMask(ip.To4()), nil
}

// FormatIPv4Mask net.IPMask → string
func FormatIPv4Mask(ip net.IPMask) string { return ip.String() }

// ParseBytes string → uint64
func ParseBytes(s string) (uint64, error) { return strconv.ParseUint(s, 10, 64) }

// FormatBytes uint64 → string
func FormatBytes(b uint64) string { return strconv.FormatUint(b, 10) }

// ParseFile string → *os.File
func ParseFile(s string) (*os.File, error) { return os.Open(s) }

// FormatFile *os.File → string
func FormatFile(f *os.File) string { return f.Name() }

// ParseTime string → time.Time
func ParseTime(s string) (time.Time, error) { return time.Parse(time.RFC3339, s) }

// FormatTime time.Time → string
func FormatTime(t time.Time) string { return t.Format(time.RFC3339) }
