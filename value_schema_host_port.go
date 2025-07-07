package tinyflags

import (
	"fmt"
	"net/url"
)

// SchemaHostPortP defines a string flag that must match scheme://host:port format,
// with a shorthand, default value, and usage string.
// The value is parsed with url.Parse and must include a scheme and host.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) SchemaHostPortP(name, short, def string, usage string) *Flag[string] {
	ptr := new(string)
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (string, error) {
			u, err := url.Parse(s)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return "", fmt.Errorf("invalid scheme://host:port format")
			}
			return s, nil
		},
		func(s string) string { return s },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// SchemaHostPort defines a string flag that must match scheme://host:port format,
// with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) SchemaHostPort(name, def string, usage string) *Flag[string] {
	return f.SchemaHostPortP(name, "", def, usage)
}

// SchemaHostPortVarP defines a string flag that must match scheme://host:port format,
// with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) SchemaHostPortVarP(ptr *string, name, short string, def string, usage string) *Flag[string] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (string, error) {
			u, err := url.Parse(s)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return "", fmt.Errorf("invalid scheme://host:port format")
			}
			return s, nil
		},
		func(s string) string { return s },
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// SchemaHostPortVar defines a string flag that must match scheme://host:port format,
// with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) SchemaHostPortVar(ptr *string, name string, def string, usage string) *Flag[string] {
	return f.SchemaHostPortVarP(ptr, name, "", def, usage)
}
