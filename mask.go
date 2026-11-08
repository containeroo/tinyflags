package tinyflags

import "strings"

// MaskFirstLast masks a string by showing only the first and last characters.
// Non-string values are returned as-is. String slices are masked per element.
func MaskFirstLast(v any) any {
	return maskFirstLast(v, 0)
}

// MaskFirstLastLimit returns a mask function that preserves the first and last
// characters while limiting the masked output to at most limit characters.
func MaskFirstLastLimit(limit int) func(any) any {
	if limit < 3 {
		panic("MaskFirstLastLimit: limit must be at least 3")
	}

	return func(v any) any {
		return maskFirstLast(v, limit)
	}
}

// maskFirstLast masks supported values while optionally limiting the output length.
// A limit of zero means no limit.
func maskFirstLast(v any, limit int) any {
	switch val := v.(type) {
	case string:
		return maskFirstLastString(val, limit)
	case []string:
		out := make([]string, len(val))
		for i, s := range val {
			out[i] = maskFirstLastString(s, limit)
		}
		return out
	default:
		return v
	}
}

// maskFirstLastString masks a string while preserving its first and last characters.
// A limit of zero means no limit.
func maskFirstLastString(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= 2 {
		return s
	}

	length := len(runes)
	if limit > 0 {
		length = min(length, limit)
	}

	return string(runes[0]) +
		strings.Repeat("*", length-2) +
		string(runes[len(runes)-1])
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

// maskPostgresURLString masks credentials in a Postgres connection string.
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

			creds := rest[:at]
			masked := maskCredentialPair(creds)

			return prefix + masked + rest[at:]
		}
	}

	return s
}

// maskCredentialPair masks a username/password credential pair.
func maskCredentialPair(creds string) string {
	parts := strings.SplitN(creds, ":", 2)
	if len(parts) == 1 {
		if parts[0] == "" {
			return creds
		}

		return strings.Repeat("*", len(parts[0]))
	}

	user := parts[0]
	pass := parts[1]

	if user == "" && pass == "" {
		return creds
	}
	if user == "" {
		return ":" + strings.Repeat("*", len(pass))
	}
	if pass == "" {
		return strings.Repeat("*", len(user)) + ":"
	}

	return strings.Repeat("*", len(user)) + ":" + strings.Repeat("*", len(pass))
}
