package dynamic

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/builder"
)

// BoolFlag represents a dynamic boolean flag with per-ID values.
type BoolFlag struct {
	*builder.DynamicFlag[bool]                           // Embedded base flag metadata
	item                       *DynamicScalarValue[bool] // Parsed values and defaults
	strictMode                 *bool                     // Pointer to shared strict mode flag
}

// Strict enables strict mode on this flag and returns itself.
func (b *BoolFlag) Strict() *BoolFlag {
	*b.strictMode = true
	return b
}

// Get returns the value for a given ID and whether it exists.
func (f *BoolFlag) Get(id string) (bool, bool) {
	val, ok := f.item.values[id]
	if !ok {
		return f.item.def, false
	}
	return val, true
}

// MustGet returns the value for a given ID or panics if missing.
func (f *BoolFlag) MustGet(id string) bool {
	val, ok := f.Get(id)
	if !ok {
		panic(fmt.Sprintf("missing value for bool flag: %s (%s)", f.item.field, id))
	}
	return val
}

// Has reports whether a value was set for the given ID.
func (f *BoolFlag) Has(id string) bool {
	_, ok := f.item.values[id]
	return ok
}

// Values returns all parsed values keyed by ID.
func (f *BoolFlag) Values() map[string]bool {
	return f.item.values
}

// ValuesAny returns all parsed values as a map of any.
func (f *BoolFlag) ValuesAny() map[string]any {
	out := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		out[k] = v
	}
	return out
}
