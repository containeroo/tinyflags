package dynamic

import (
	"github.com/containeroo/tinyflags/internal/builder"
)

// BoolValue wraps a DynamicScalarValue[bool] and exposes strict-mode information.
type BoolValue struct {
	*DynamicScalarValue[bool]
	strictMode bool
}

// IsStrictBool reports whether the flag requires an explicit value (--flag=true/false).
func (b *BoolValue) IsStrictBool() bool {
	return b.strictMode
}

// BoolFlag provides fluent builder methods for dynamic boolean flags.
type BoolFlag struct {
	*builder.DynamicFlag[bool]
	item       *DynamicScalarValue[bool]
	strictMode *bool
}

// Strict requires the flag to be passed as --flag=true|false.
func (b *BoolFlag) Strict() *BoolFlag {
	*b.strictMode = true
	return b
}

// Get retrieves the parsed value for the given ID, falling back to default.
func (f *BoolFlag) Get(id string) (bool, bool) {
	val, ok := f.item.values[id]
	if !ok {
		return f.item.def, false
	}
	return val, true
}

// MustGet retrieves the parsed value or panics if not set.
func (f *BoolFlag) MustGet(id string) bool {
	val, ok := f.Get(id)
	if !ok {
		panic("value not set for dynamic bool flag: " + f.item.field + " (" + id + ")")
	}
	return val
}

// Has reports whether the value was explicitly set for this ID.
func (f *BoolFlag) Has(id string) bool {
	_, ok := f.item.values[id]
	return ok
}

// Values returns all stored values for all IDs.
func (f *BoolFlag) Values() map[string]bool {
	return f.item.values
}

// ValuesAny returns all stored values as map[string]any for interface use.
func (f *BoolFlag) ValuesAny() map[string]any {
	m := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		m[k] = v
	}
	return m
}
