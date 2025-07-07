package tinyflags

import (
	"fmt"
	"net/url"
)

// SchemaHostPortSliceP defines a string flag that must represent a valid schema://host:port pair,
// with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) SchemaHostPortSliceP(name, short string, def []string, usage string) *SliceFlag[[]string] {
	ptr := new([]string)
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (string, error) {
			u, err := url.Parse(s)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return "", fmt.Errorf("invalid scheme://host:port format: %q", s)
			}
			return s, nil
		},
		func(s string) string { return s },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// SchemaHostPortSlice defines a string slice flag that must represent a valid schema://host:port pair,
// with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) SchemaHostPortSlice(name string, def []string, usage string) *SliceFlag[[]string] {
	return f.SchemaHostPortSliceP(name, "", def, usage)
}

// SchemaHostPortSliceVarP defines a string flag that must represent a valid schema://host:port pair,
// with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value
func (f *FlagSet) SchemaHostPortSliceVarP(ptr *[]string, name, short string, def []string, usage string) *SliceFlag[[]string] {
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (string, error) {
			u, err := url.Parse(s)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return "", fmt.Errorf("invalid scheme://host:port format: %q", s)
			}
			return s, nil
		},
		func(s string) string { return s },
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// SchemaHostPortSliceVar defines a string flag that must represent a valid schema://host:port pair,
// with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value
func (f *FlagSet) SchemaHostPortSliceVar(ptr *[]string, name string, def []string, usage string) *SliceFlag[[]string] {
	return f.SchemaHostPortSliceVarP(ptr, name, "", def, usage)
}
