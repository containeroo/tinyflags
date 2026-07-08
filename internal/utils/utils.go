package utils

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// AllowOnly returns a validator function that permits only values from the given list.
// Each value is compared using the provided formatting function.
// The format function is used to normalize or stringify values for comparison,
// which is useful for case-insensitive or structured types.
func AllowOnly[T any](format func(T) string, allowed []T) func(T) error {
	return func(v T) error {
		got := format(v)
		for _, a := range allowed {
			if format(a) == got {
				return nil
			}
		}
		return fmt.Errorf("%q must be one of: %s", got, JoinFormatted(allowed, format))
	}
}

// FormatList applies the given format function to each element of the slice
// and returns a new slice of formatted strings.
// This is useful for preparing human-readable output from a slice of structured types.
func FormatList[T any](format func(T) string, values []T) []string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = format(v)
	}
	return out
}

// JoinFormatted applies the given format function to each value and joins them with commas.
// If format is nil, it falls back to using fmt.Sprintf("%v", values).
func JoinFormatted[T any](values []T, format func(T) string) string {
	if format == nil {
		return fmt.Sprintf("%v", values) // fallback
	}
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = format(v)
	}
	return strings.Join(parts, ", ")
}

// PluralSuffix returns "s" if the given number is not 1, otherwise it returns an empty string.
// It is useful for constructing basic pluralized words like "flag" or "flags".
func PluralSuffix(i int) string {
	if i != 1 {
		return "s"
	}
	return ""
}

// ParseString parses a string value.
func ParseString(s string) (string, error) { return s, nil }

// FormatString formats a string value.
func FormatString(s string) string { return s }

// ParseTCPAddr parses a TCP address.
func ParseTCPAddr(s string) (*net.TCPAddr, error) {
	addr, err := net.ResolveTCPAddr("tcp", s)
	if err != nil {
		return nil, fmt.Errorf("invalid TCP address %q: %w", s, err)
	}
	return addr, nil
}

// FormatTCPAddr formats a TCP address.
func FormatTCPAddr(addr *net.TCPAddr) string {
	if addr == nil {
		return ""
	}
	return addr.String()
}

// ParseFloat64 parses a float64 value.
func ParseFloat64(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

// FormatFloat64 formats a float64 value.
func FormatFloat64(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

// ParseFloat32 parses a float32 value.
func ParseFloat32(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

// FormatFloat32 formats a float32 value.
func FormatFloat32(f float32) string { return strconv.FormatFloat(float64(f), 'f', -1, 32) }

// ParseIP parses an IP address.
func ParseIP(s string) (net.IP, error) { return net.ParseIP(s), nil }

// FormatIP formats an IP address.
func FormatIP(ip net.IP) string { return ip.String() }

// ParseIPv4Mask parses an IPv4 mask.
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

// FormatIPv4Mask formats an IPv4 mask.
func FormatIPv4Mask(ip net.IPMask) string { return ip.String() }

// ParseBytes parses a byte count.
func ParseBytes(s string) (uint64, error) { return strconv.ParseUint(s, 10, 64) }

// FormatBytes formats a byte count.
func FormatBytes(b uint64) string { return strconv.FormatUint(b, 10) }

// ParseFile opens a file for reading.
func ParseFile(s string) (*os.File, error) { return os.Open(s) }

// FormatFile formats a file name.
func FormatFile(f *os.File) string { return f.Name() }

// ParseTime parses a time value.
func ParseTime(s string) (time.Time, error) { return time.Parse(time.RFC3339, s) }

// FormatTime formats a time value.
func FormatTime(t time.Time) string { return t.Format(time.RFC3339) }
