package dynamic

import (
	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/utils"
)

// DynamicScalarValue holds parsed scalar values per ID with parsing, formatting, and validation.
type DynamicScalarValue[T any] struct {
	field         string             // Flag field name
	def           T                  // Default value
	baseDef       T                  // Original default value
	changed       bool               // Whether the value was changed
	hooks         core.ValueHooks[T] // Shared parse/format/validate/finalize behavior
	finalizeID    func(string, T) T  // Optional finalizer with ID
	values        map[string]T       // Parsed values per ID
	allowOverride bool               // Allow re-assignment per ID when true.
}

// NewDynamicScalarValue creates a new dynamic scalar value.
func NewDynamicScalarValue[T any](field string, def T, parse func(string) (T, error), format func(T) string) *DynamicScalarValue[T] {
	return &DynamicScalarValue[T]{
		field:   field,
		def:     def,
		baseDef: def,
		hooks:   core.NewValueHooks(parse, format),
		values:  make(map[string]T),
	}
}

// Set parses and stores a value for a specific ID.
func (d *DynamicScalarValue[T]) Set(id, raw string) error {
	// Enforce before any parsing/validation to catch duplicates early.
	if !d.allowOverride {
		if _, exists := d.values[id]; exists {
			return &core.DuplicatePerIDError{Field: d.field, ID: id}
		}
	}
	val, err := d.hooks.ParseValue(raw)
	if err != nil {
		return err
	}
	if d.finalizeID != nil {
		val = d.finalizeID(id, val)
	}
	d.values[id] = val
	d.changed = true
	return nil
}

// setValidate sets the optional validation function.
func (d *DynamicScalarValue[T]) setValidate(fn func(T) error) {
	d.hooks.SetValidate(fn)
}

// setFinalize sets the optional finalizer function.
func (d *DynamicScalarValue[T]) setFinalize(fn func(T) T) {
	d.hooks.SetFinalize(fn)
}

// setFinalizeDefaultValue enables running the finalizer on defaults when unset.
func (d *DynamicScalarValue[T]) setFinalizeDefaultValue() {
	d.hooks.EnableFinalizeDefault()
}

// setFinalizeWithID sets the optional finalizer function with ID context.
func (d *DynamicScalarValue[T]) setFinalizeWithID(fn func(string, T) T) {
	d.finalizeID = fn
}

// Base returns itself for generic access.
func (d *DynamicScalarValue[T]) Base() *DynamicScalarValue[T] {
	return d
}

// ApplyDefaultFinalize applies the default-only finalizer for unset IDs.
func (d *DynamicScalarValue[T]) ApplyDefaultFinalize() {
	d.hooks.ApplyDefaultScalar(&d.def, false)
}

// FieldName returns the field name of the flag.
func (d *DynamicScalarValue[T]) FieldName() string {
	return d.field
}

// GetAny returns the value as any for a given ID, falling back to default.
func (d *DynamicScalarValue[T]) GetAny(id string) (any, bool) {
	val, ok := d.values[id]
	if ok {
		return val, true
	}
	return d.def, false
}

// ValuesAny returns all stored values as a map of any.
func (d *DynamicScalarValue[T]) ValuesAny() map[string]any {
	out := make(map[string]any, len(d.values))
	for k, v := range d.values {
		out[k] = v
	}
	return out
}

// ResetParseState clears all parsed IDs and restores default-finalizer state.
func (d *DynamicScalarValue[T]) ResetParseState() {
	clear(d.values)
	utils.ResetScalarState(&d.def, d.baseDef, &d.changed, &d.hooks.DefaultFinalized)
}
