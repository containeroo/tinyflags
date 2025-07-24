package dynamic

import "github.com/containeroo/tinyflags/internal/builder"

// BoolValue wraps a DynamicScalarValue[bool] and lets the FSM see IsStrictBool.
type BoolValue struct {
	item       *DynamicScalarValue[bool]
	strictMode *bool
}

// Set parses and stores one entry.
func (b *BoolValue) Set(id, raw string) error {
	return b.item.Set(id, raw)
}

func (b *BoolValue) FieldName() string {
	return b.item.field
}

func (b *BoolValue) GetAny(id string) (any, bool) {
	val, ok := b.item.values[id]
	if ok {
		return val, true
	}
	return b.item.def, false
}

// IsStrictBool reports whether the flag requires an explicit value (--flag=true/false).
func (b *BoolValue) IsStrictBool() bool {
	return *b.strictMode
}

// BoolFlag is the builder for per-ID boolean flags.
type BoolFlag struct {
	*builder.DynamicFlag[bool]
	item       *DynamicScalarValue[bool]
	strictMode bool
}

// Strict requires explicit `=true|false`.
func (b *BoolFlag) Strict() *BoolFlag {
	b.strictMode = true
	return b
}

// Get retrieves the parsed value.
func (f *BoolFlag) Get(id string) (bool, bool) {
	val, ok := f.item.values[id]
	return val, ok
}

// MustGet returns the parsed value, panicking if not set.
func (f *BoolFlag) MustGet(id string) bool {
	val, ok := f.Get(id)
	if !ok {
		panic("value not set")
	}
	return val
}

// Values returns all stored values.
func (f *BoolFlag) Values() map[string]bool {
	return f.item.values
}

// ValuesAny returns values as a generic map.
func (f *BoolFlag) ValuesAny() map[string]any {
	m := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		m[k] = v
	}
	return m
}
