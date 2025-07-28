// Package dynamic implements dynamic flag types that support multiple
// per-identifier values (e.g. --http.a.timeout, --http.b.timeout).
// This file defines dummy CLI-facing flag.Value implementations that
// act as placeholders to satisfy the core.BaseFlag.Value interface
// without interfering with internal parsing logic.
package dynamic

// placeholderValue is a dummy Value used for non-slice scalar flags.
type placeholderValue struct {
	def string // String form of default value (for help)
}

func (p *placeholderValue) Set(string) error { return nil }
func (p *placeholderValue) Get() any         { return nil }
func (p *placeholderValue) Changed() bool    { return false }
func (p *placeholderValue) Default() string  { return p.def }

// slicePlaceholder is a dummy Value used for slice flags.
type slicePlaceholder struct {
	def string // String form of default value (e.g. comma-separated)
}

func (v *slicePlaceholder) Set(string) error { return nil }
func (v *slicePlaceholder) Get() any         { return nil }
func (v *slicePlaceholder) Changed() bool    { return false }
func (v *slicePlaceholder) Default() string  { return v.def }
func (v *slicePlaceholder) IsSlice()         {} // Marker method

// boolPlaceholder is a dummy Value used for bool flags with strict support.
type boolPlaceholder struct {
	def        string // String form of default value ("true" or "false")
	strictMode *bool  // Shared strict-mode marker
}

func (v *boolPlaceholder) Set(string) error   { return nil }
func (v *boolPlaceholder) Get() any           { return nil }
func (v *boolPlaceholder) Changed() bool      { return false }
func (v *boolPlaceholder) Default() string    { return v.def }
func (v *boolPlaceholder) IsStrictBool() bool { return v.strictMode != nil && *v.strictMode }
