package tinyflags

import "strings"

// MaskFirstLast masks a string by showing only the first and last characters.
// Non-string values are returned as-is. String slices are masked per element.
func MaskFirstLast(v any) any {
	switch val := v.(type) {
	case string:
		return maskFirstLastString(val)
	case []string:
		out := make([]string, len(val))
		for i, s := range val {
			out[i] = maskFirstLastString(s)
		}
		return out
	default:
		return v
	}
}

func maskFirstLastString(s string) string {
	runes := []rune(s)
	if len(runes) <= 2 {
		return s
	}
	return string(runes[0]) + strings.Repeat("*", len(runes)-2) + string(runes[len(runes)-1])
}

// MaskPostgresURL masks credentials in postgres URLs (e.g. postgres://user:pass@host/db).
// Non-string values are returned as-is. String slices are masked per element.
func MaskPostgresURL(v any) any {
	switch val := v.(type) {
	case string:
		return maskPostgresURLString(val)
	case []string:
		out := make([]string, len(val))
		for i, s := range val {
			out[i] = maskPostgresURLString(s)
		}
		return out
	default:
		return v
	}
}

func maskPostgresURLString(s string) string {
	for _, prefix := range []string{"postgres://", "postgresql://"} {
		if strings.HasPrefix(s, prefix) {
			rest := s[len(prefix):]
			at := strings.Index(rest, "@")
			if at == -1 {
				return s
			}
			if at == 0 {
				return s
			}
			return prefix + strings.Repeat("*", at) + rest[at:]
		}
	}
	return s
}
