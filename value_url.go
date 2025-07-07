package tinyflags

import (
	"fmt"
	"net/url"
)

// URLP defines a *url.URL flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) URLP(name, short string, def *url.URL, usage string) *Flag[*url.URL] {
	ptr := new(*url.URL)
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (*url.URL, error) {
			u, err := url.Parse(s)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return nil, fmt.Errorf("invalid URL: must include scheme and host")
			}
			return u, nil
		},
		func(u *url.URL) string {
			if u == nil {
				return ""
			}
			return u.String()
		},
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// URL defines a *url.URL flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) URL(name string, def *url.URL, usage string) *Flag[*url.URL] {
	return f.URLP(name, "", def, usage)
}

// URLVarP defines a *url.URL flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) URLVarP(ptr **url.URL, name, short string, def *url.URL, usage string) *Flag[*url.URL] {
	val := NewFlagItem(
		ptr,
		def,
		func(s string) (*url.URL, error) {
			u, err := url.Parse(s)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return nil, fmt.Errorf("invalid URL: must include scheme and host")
			}
			return u, nil
		},
		func(u *url.URL) string {
			if u == nil {
				return ""
			}
			return u.String()
		},
	)
	return addScalar(f, name, short, usage, val, ptr)
}

// URLVar defines a *url.URL flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) URLVar(ptr **url.URL, name string, def *url.URL, usage string) *Flag[*url.URL] {
	return f.URLVarP(ptr, name, "", def, usage)
}
