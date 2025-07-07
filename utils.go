package tinyflags

import (
	"errors"
	"fmt"
	"strings"
)

// addScalar registers a scalar flag and returns its builder.
func addScalar[T any](f *FlagSet, name, short, usage string, val Value, ptr *T) *Flag[T] {
	bf := &baseFlag{
		name:  name,  // long name: --flag
		short: short, // short name: -f
		usage: usage, // help text
		value: val,   // parsed value
	}
	f.flags[name] = bf                      // register in lookup map
	f.registered = append(f.registered, bf) // preserve order

	return &Flag[T]{fs: f, bf: bf, ptr: ptr}
}

// addSlice registers a slice flag and returns its slice builder.
func addSlice[T any](f *FlagSet, name, short, usage string, val Value, ptr *T) *SliceFlag[T] {
	bf := &baseFlag{
		name:  name,
		short: short,
		usage: usage,
		value: val,
	}
	f.flags[name] = bf
	f.registered = append(f.registered, bf)

	return &SliceFlag[T]{Flag: Flag[T]{fs: f, bf: bf, ptr: ptr}}
}

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

// IsHelpRequested checks if the error is a HelpRequested sentinel.
func IsHelpRequested(err error) bool {
	var helpErr *HelpRequested
	return errors.As(err, &helpErr)
}

// IsVersionRequested checks if the error is a VersionRequested sentinel
func IsVersionRequested(err error) bool {
	var versionErr *VersionRequested
	return errors.As(err, &versionErr)
}

// RequestHelp returns an error with type HelpRequested and the given message.
func RequestHelp(msg string) error {
	return &HelpRequested{Message: msg}
}

// RequestVersion returns an error with type VersionRequested and the given message.
func RequestVersion(msg string) error {
	return &VersionRequested{Version: msg}
}
