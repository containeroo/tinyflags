package tinyflags

import (
	"fmt"
	"strings"
)

// formatAllowed returns a comma-separated list of allowed values.
func formatAllowed[T any](allowed []T, format func(T) string) string {
	if format == nil {
		return fmt.Sprintf("%v", allowed) // fallback
	}
	formatted := make([]string, len(allowed))
	for i, a := range allowed {
		formatted[i] = format(a) // format each value
	}
	return strings.Join(formatted, ", ")
}

// pluralSuffix returns "s" if the given number is not 1, otherwise it returns an empty string.
// It is useful for constructing basic pluralized words like "flag" or "flags".
func pluralSuffix(i int) string {
	if i != 1 {
		return "s"
	}
	return ""
}
