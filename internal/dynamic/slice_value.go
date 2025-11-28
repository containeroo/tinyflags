package dynamic

import (
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/utils"
)

// DynamicSliceValue holds parsed slice values per ID with parsing, formatting, and validation.
type DynamicSliceValue[T any] struct {
	field      string                  // Flag field name
	def        []T                     // Default slice value
	changed    bool                    // Whether the value was changed
	parse      func(string) (T, error) // Function to parse a single element
	format     func(T) string          // Function to format a single element
	delimiter  string                  // Separator used to split input
	validate   func(T) error           // Optional validation function
	finalize   (func(T) T)             // Optional finalizer function
	finalizeID func(string, T) T       // Optional finalizer function with ID
	strictDel  bool                    // Reject mixed delimiters when true
	allowEmpty bool                    // Allow empty items when true
	values     map[string][]T          // Parsed values per ID
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
		def:       def,
		parse:     parse,
		format:    format,
		delimiter: delimiter,
		values:    make(map[string][]T),
	}
}

// Set parses and stores one or more values for a given ID.
func (d *DynamicSliceValue[T]) Set(id, raw string) error {
	if d.strictDel {
		for _, alt := range []string{",", ";", "|"} {
			if alt == d.delimiter {
				continue
			}
			if strings.Contains(raw, alt) {
				return fmt.Errorf("mixed delimiters: found %q while using %q", alt, d.delimiter)
			}
		}
	}

	for _, chunk := range strings.Split(raw, d.delimiter) {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" && !d.allowEmpty {
			return fmt.Errorf("invalid value %q: empty values are not allowed", chunk)
		}

		val, err := d.parse(chunk)
		if err != nil {
			return fmt.Errorf("invalid %q: %w", chunk, err)
		}
		val, err = utils.ApplyValueHooks(val, d.validate, d.finalize)
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
	d.validate = fn
}

// setFinalize sets a per-item finalizer function.
func (d *DynamicSliceValue[T]) setFinalize(fn func(T) T) {
	d.finalize = fn
}

// setFinalizeWithID sets a per-item finalizer with access to the ID.
func (d *DynamicSliceValue[T]) setFinalizeWithID(fn func(string, T) T) {
	d.finalizeID = fn
}

// setDelimiter sets the delimiter used to split input values.
func (d *DynamicSliceValue[T]) setDelimiter(sep string) {
	d.delimiter = sep
}

// setStrictDelimiter toggles mixed-delimiter rejection.
func (d *DynamicSliceValue[T]) setStrictDelimiter(strict bool) {
	d.strictDel = strict
}

// setAllowEmpty toggles acceptance of empty items.
func (d *DynamicSliceValue[T]) setAllowEmpty(allow bool) {
	d.allowEmpty = allow
}

// FieldName returns the field name of the flag.
func (d *DynamicSliceValue[T]) FieldName() string {
	return d.field
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
