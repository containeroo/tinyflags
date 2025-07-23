package dynamic

import (
	"fmt"
	"strings"
)

// DynamicSliceValue parses comma (or custom-sep) lists per ID.
// For example: --http.alpha.tags=a,b
type DynamicSliceValue[T any] struct {
	field     string                  // field name (e.g. "tags")
	def       []T                     // default value
	parse     func(string) (T, error) // parser for each item
	format    func(T) string          // formatter for output and help
	delimiter string                  // input separator (default: ",")
	validate  func(T) error           // per-item validation function
	values    map[string][]T          // instance ID → parsed slice values
}

// NewDynamicSliceValue builds a per-ID slice parser.
func NewDynamicSliceValue[T any](
	field string,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *DynamicSliceValue[T] {
	return &DynamicSliceValue[T]{
		def:       def,
		field:     field,
		parse:     parse,
		format:    format,
		delimiter: delimiter,
		values:    make(map[string][]T),
	}
}

// Set parses and stores the slice from a delimited string for a given ID.
func (d *DynamicSliceValue[T]) Set(id, input string) error {
	for _, raw := range strings.Split(input, d.delimiter) {
		val, err := d.parse(strings.TrimSpace(raw))
		if err != nil {
			return fmt.Errorf("invalid item %q: %w", raw, err)
		}
		if d.validate != nil {
			if err := d.validate(val); err != nil {
				return fmt.Errorf("invalid value %q: %w", raw, err)
			}
		}
		d.values[id] = append(d.values[id], val)
	}
	return nil
}

// SetValidator sets a per-item validation function.
func (d *DynamicSliceValue[T]) setValidate(fn func(T) error) {
	d.validate = fn
}

// SetDelimiter sets a custom delimiter for parsing the slice.
func (d *DynamicSliceValue[T]) setDelimiter(sep string) {
	d.delimiter = sep
}

func (d *DynamicSliceValue[T]) ValuesAny() map[string]any {
	m := make(map[string]any, len(d.values))
	for k, v := range d.values {
		m[k] = v
	}
	return m
}
