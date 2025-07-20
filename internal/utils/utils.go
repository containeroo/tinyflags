package utils

import (
	"fmt"
	"strings"
)

// FormatAllowed returns a comma-separated list of allowed values.
func FormatAllowed[T any](allowed []T, format func(T) string) string {
	if format == nil {
		return fmt.Sprintf("%v", allowed) // fallback
	}
	formatted := make([]string, len(allowed))
	for i, a := range allowed {
		formatted[i] = format(a) // format each value
	}
	return strings.Join(formatted, ", ")
}

// PluralSuffix returns "s" if the given number is not 1, otherwise it returns an empty string.
// It is useful for constructing basic pluralized words like "flag" or "flags".
func PluralSuffix(i int) string {
	if i != 1 {
		return "s"
	}
	return ""
}
