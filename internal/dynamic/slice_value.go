package dynamic

import (
	"fmt"
	"strings"
)

// DynamicSliceValue stores a slice of values per instance ID.
// For example: --http.alpha.tags=a,b
type DynamicSliceValue[T any] struct {
	field     string                  // field name (e.g. "tags")
	parse     func(string) (T, error) // parser for each item
	format    func(T) string          // formatter for output and help
	delimiter string                  // input separator (default: ",")
	validator func(T) error           // per-item validation function
	values    map[string][]T          // instance ID → parsed slice values
}

// NewDynamicSliceValue creates a new dynamic slice value parser.
func NewDynamicSliceValue[T any](
	field string,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *DynamicSliceValue[T] {
	return &DynamicSliceValue[T]{
		field:     field,
		parse:     parse,
		format:    format,
		delimiter: delimiter,
		values:    make(map[string][]T),
	}
}

// Set parses and stores the slice from a delimited string for a given ID.
func (d *DynamicSliceValue[T]) Set(id, input string) error {
	parts := strings.Split(input, d.delimiter)
	for _, raw := range parts {
		val, err := d.parse(strings.TrimSpace(raw))
		if err != nil {
			return fmt.Errorf("invalid item %q: %w", raw, err)
		}
		if d.validator != nil {
			if err := d.validator(val); err != nil {
				return fmt.Errorf("invalid value %q: %w", raw, err)
			}
		}
		d.values[id] = append(d.values[id], val)
	}
	return nil
}

// Get returns the parsed slice for the given ID.
func (d *DynamicSliceValue[T]) Get(id string) ([]T, bool) {
	val, ok := d.values[id]
	return val, ok
}

// Values returns all instance values.
func (d *DynamicSliceValue[T]) Values() map[string][]T {
	return d.values
}

// ValuesAny returns values as a map[string]any for interface compatibility.
func (d *DynamicSliceValue[T]) ValuesAny() map[string]any {
	out := make(map[string]any, len(d.values))
	for k, v := range d.values {
		out[k] = v
	}
	return out
}

// SetValidator sets a per-item validation function.
func (d *DynamicSliceValue[T]) setValidate(fn func(T) error) {
	d.validator = fn
}

// SetDelimiter sets a custom delimiter for parsing the slice.
func (d *DynamicSliceValue[T]) setDelimiter(sep string) {
	d.delimiter = sep
}
