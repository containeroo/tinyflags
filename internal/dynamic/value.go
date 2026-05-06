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

// Set ignores placeholder input.
func (p *placeholderValue) Set(string) error { return nil }

// Get returns no concrete placeholder value.
func (p *placeholderValue) Get() any { return nil }

// Changed reports that placeholders are never user-set.
func (p *placeholderValue) Changed() bool { return false }

// Default returns the placeholder default string.
func (p *placeholderValue) Default() string { return p.def }

// slicePlaceholder is a dummy Value used for slice flags.
type slicePlaceholder struct {
	def string // String form of default value (e.g. comma-separated)
}

// Set ignores placeholder input.
func (v *slicePlaceholder) Set(string) error { return nil }

// Get returns no concrete placeholder value.
func (v *slicePlaceholder) Get() any { return nil }

// Changed reports that placeholders are never user-set.
func (v *slicePlaceholder) Changed() bool { return false }

// Default returns the placeholder default string.
func (v *slicePlaceholder) Default() string { return v.def }

// IsSlice marks the placeholder as slice-backed.
func (v *slicePlaceholder) IsSlice() {} // Marker method

// boolPlaceholder is a dummy Value used for bool flags with strict support.
type boolPlaceholder struct {
	def        string // String form of default value ("true" or "false")
	strictMode *bool  // Shared strict-mode marker
}

// Set ignores placeholder input.
func (v *boolPlaceholder) Set(string) error { return nil }

// Get returns no concrete placeholder value.
func (v *boolPlaceholder) Get() any { return nil }

// Changed reports that placeholders are never user-set.
func (v *boolPlaceholder) Changed() bool { return false }

// Default returns the placeholder default string.
func (v *boolPlaceholder) Default() string { return v.def }

// IsStrictBool reports whether strict bool parsing is enabled.
func (v *boolPlaceholder) IsStrictBool() bool { return v.strictMode != nil && *v.strictMode }
