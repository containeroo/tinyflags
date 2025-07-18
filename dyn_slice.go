package tinyflags

import (
	"fmt"
	"strings"
)

// DynamicSliceItem stores multiple values per instance ID, parsed from delimited strings.
type DynamicSliceItem[T any] struct {
	field     string
	parse     func(string) (T, error)
	format    func(T) string
	delimiter string
	validator func(T) error
	values    map[string][]T
}

// NewDynamicSliceItem constructs a new dynamic slice item.
func NewDynamicSliceItem[T any](
	field string,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *DynamicSliceItem[T] {
	return &DynamicSliceItem[T]{
		field:     field,
		parse:     parse,
		format:    format,
		delimiter: delimiter,
		values:    make(map[string][]T),
	}
}

// Set parses and appends to the slice for a given instance ID.
func (d *DynamicSliceItem[T]) Set(id, input string) error {
	parts := strings.Split(input, d.delimiter)
	for _, raw := range parts {
		val, err := d.parse(strings.TrimSpace(raw))
		if err != nil {
			return fmt.Errorf("invalid item %q for --%s.%s.%s: %w", raw, d.field, id, d.field, err)
		}
		if d.validator != nil {
			if err := d.validator(val); err != nil {
				return fmt.Errorf("invalid item %q: %w", raw, err)
			}
		}
		d.values[id] = append(d.values[id], val)
	}
	return nil
}

// Get returns the slice for the given instance ID.
func (d *DynamicSliceItem[T]) Get(id string) ([]T, bool) {
	val, ok := d.values[id]
	return val, ok
}

// Values returns the full map of instanceID → slice.
func (d *DynamicSliceItem[T]) Values() map[string][]T {
	return d.values
}

// ValuesAny returns values as map[string]any.
func (d *DynamicSliceItem[T]) ValuesAny() map[string]any {
	m := make(map[string]any, len(d.values))
	for id, vals := range d.values {
		m[id] = vals
	}
	return m
}

// SetValidator sets a validation function for slice items.
func (d *DynamicSliceItem[T]) SetValidator(fn func(T) error) {
	d.validator = fn
}

// DynamicSliceFlag is a builder for dynamic slice flags.
type DynamicSliceFlag[T any] struct {
	builderBase[[]T]
	item *DynamicSliceItem[T]
}

// Get retrieves the slice for the given instance ID.
func (d *DynamicSliceFlag[T]) Get(id string) ([]T, bool) {
	return d.item.Get(id)
}

// MustGet retrieves the slice or panics if not found.
func (d *DynamicSliceFlag[T]) MustGet(id string) []T {
	val, ok := d.item.Get(id)
	if !ok {
		panic(fmt.Sprintf("value for id %q not set", id))
	}
	return val
}

// Values returns all parsed slices keyed by instance ID.
func (d *DynamicSliceFlag[T]) Values() map[string][]T {
	return d.item.Values()
}

// ValuesAny returns parsed values as map[string]any.
func (d *DynamicSliceFlag[T]) ValuesAny() map[string]any {
	return d.item.ValuesAny()
}
