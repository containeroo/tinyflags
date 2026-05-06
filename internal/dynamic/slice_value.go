package dynamic

import (
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/utils"
)

// DynamicSliceValue holds parsed slice values per ID with parsing, formatting, and validation.
type DynamicSliceValue[T any] struct {
	field            string                  // Flag field name
	def              []T                     // Default slice value
	baseDef          []T                     // Original default slice value
	changed          bool                    // Whether the value was changed
	input            core.SliceInputConfig   // Shared slice-input behavior
	hooks            core.ValueHooks[T]      // Shared parse/format/validate/finalize behavior
	finalizeID       func(string, T) T       // Optional finalizer function with ID
	values           map[string][]T          // Parsed values per ID
}

// NewDynamicSliceValue creates a new dynamic slice value.
func NewDynamicSliceValue[T any](
	field string,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *DynamicSliceValue[T] {
	return &DynamicSliceValue[T]{
		field:     field,
		def:       append([]T(nil), def...),
		baseDef:   append([]T(nil), def...),
		input:     core.SliceInputConfig{Delimiter: delimiter},
		hooks:     core.NewValueHooks(parse, format),
		values:    make(map[string][]T),
	}
}

// Set parses and stores one or more values for a given ID.
func (d *DynamicSliceValue[T]) Set(id, raw string) error {
	chunks, err := d.input.Split(raw)
	if err != nil {
		return err
	}

	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" && !d.input.AllowEmpty {
			return fmt.Errorf("invalid value %q: empty values are not allowed", chunk)
		}

		val, err := d.hooks.ParseValue(chunk)
		if err != nil {
			return fmt.Errorf("invalid value %q: %w", chunk, err)
		}
		if d.finalizeID != nil {
			val = d.finalizeID(id, val)
		}
		d.values[id] = append(d.values[id], val)
	}
	d.changed = true
	return nil
}

// setValidate sets a validation function for individual elements.
func (d *DynamicSliceValue[T]) setValidate(fn func(T) error) {
	d.hooks.SetValidate(fn)
}

// setFinalize sets a per-item finalizer function.
func (d *DynamicSliceValue[T]) setFinalize(fn func(T) T) {
	d.hooks.SetFinalize(fn)
}

// setFinalizeDefaultValue enables running the finalizer on defaults when unset.
func (d *DynamicSliceValue[T]) setFinalizeDefaultValue() {
	d.hooks.EnableFinalizeDefault()
}

// setFinalizeWithID sets a per-item finalizer with access to the ID.
func (d *DynamicSliceValue[T]) setFinalizeWithID(fn func(string, T) T) {
	d.finalizeID = fn
}

// setDelimiter sets the delimiter used to split input values.
func (d *DynamicSliceValue[T]) setDelimiter(sep string) {
	d.input.Delimiter = sep
}

// setStrictDelimiter toggles mixed-delimiter rejection.
func (d *DynamicSliceValue[T]) setStrictDelimiter(strict bool) {
	d.input.StrictDel = strict
}

// setAllowEmpty toggles acceptance of empty items.
func (d *DynamicSliceValue[T]) setAllowEmpty(allow bool) {
	d.input.AllowEmpty = allow
}

// FieldName returns the field name of the flag.
func (d *DynamicSliceValue[T]) FieldName() string {
	return d.field
}

// ApplyDefaultFinalize applies the default-only finalizer for unset IDs.
func (d *DynamicSliceValue[T]) ApplyDefaultFinalize() {
	d.hooks.ApplyDefaultSlice(d.def, false)
}

// GetAny returns the slice as any for a given ID, falling back to default.
func (d *DynamicSliceValue[T]) GetAny(id string) (any, bool) {
	val, ok := d.values[id]
	if ok {
		return val, true
	}
	return d.def, false
}

// ValuesAny returns all values as a map of any.
func (d *DynamicSliceValue[T]) ValuesAny() map[string]any {
	out := make(map[string]any, len(d.values))
	for k, v := range d.values {
		out[k] = v
	}
	return out
}

// ResetParseState clears all parsed IDs and restores the original defaults.
func (d *DynamicSliceValue[T]) ResetParseState() {
	clear(d.values)
	utils.ResetSliceState(&d.def, d.baseDef, &d.changed, &d.hooks.DefaultFinalized)
}
