package tinyflags

import (
	"fmt"
	"net/url"
)

// URLSliceP defines a *url.URL slice flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) URLSliceP(name, short string, def []*url.URL, usage string) *SliceFlag[[]*url.URL] {
	ptr := new([]*url.URL)
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (*url.URL, error) {
			u, err := url.Parse(s)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return nil, fmt.Errorf("invalid URL: must include scheme and host: %q", s)
			}
			return u, nil
		},
		func(u *url.URL) string {
			if u == nil {
				return ""
			}
			return u.String()
		},
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// URLSlice defines a *url.URL slice flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) URLSlice(name string, def []*url.URL, usage string) *SliceFlag[[]*url.URL] {
	return f.URLSliceP(name, "", def, usage)
}

// URLSliceP defines a *url.URL flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) URLSliceVarP(ptr *[]*url.URL, name, short string, def []*url.URL, usage string) *SliceFlag[[]*url.URL] {
	val := NewSliceItem(
		ptr,
		def,
		func(s string) (*url.URL, error) {
			u, err := url.Parse(s)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return nil, fmt.Errorf("invalid URL: must include scheme and host: %q", s)
			}
			return u, nil
		},
		func(u *url.URL) string {
			if u == nil {
				return ""
			}
			return u.String()
		},
		f.defaultDelimiter,
	)
	return addSlice(f, name, short, usage, val, ptr)
}

// URLSlice defines a *url.URL slice flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) URLSliceVar(ptr *[]*url.URL, name string, def []*url.URL, usage string) *SliceFlag[[]*url.URL] {
	return f.URLSliceVarP(ptr, name, "", def, usage)
}
