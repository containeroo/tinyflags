package tinyflags

import (
	"fmt"
	"strings"
)

// DynamicSliceFlag is a builder for dynamic slice flags.
// Each instance (e.g. `--http.alpha.tags=a,b,c`) maps to a []T.
type DynamicSliceFlag[T any] struct {
	fs   *FlagSet
	bf   *baseFlag
	item *DynamicSliceItemImpl[T]
}

// Get retrieves the slice for a given instance ID.
func (d *DynamicSliceFlag[T]) Get(id string) ([]T, bool) {
	return d.item.Get(id)
}

// MustGet retrieves the slice or panics if not set.
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

// Delimiter sets the delimiter used to split values.
func (d *DynamicSliceFlag[T]) Delimiter(sep string) *DynamicSliceFlag[T] {
	d.item.SetDelimiter(sep)
	return d
}

// Choices restricts allowed values for each slice element.
func (d *DynamicSliceFlag[T]) Choices(allowed ...T) *DynamicSliceFlag[T] {
	d.item.SetValidator(func(val T) error {
		for _, a := range allowed {
			if d.item.format(a) == d.item.format(val) {
				return nil
			}
		}
		return fmt.Errorf("must be one of [%s]", formatAllowed(allowed, d.item.format))
	})
	d.bf.allowed = make([]string, len(allowed))
	for i, a := range allowed {
		d.bf.allowed[i] = d.item.format(a)
	}
	return d
}

// DynamicSliceItemImpl stores a slice of values per instance ID.
type DynamicSliceItemImpl[T any] struct {
	field     string
	parse     func(string) (T, error)
	format    func(T) string
	delimiter string
	validator func(T) error
	values    map[string][]T
}

// NewDynamicSliceItemImpl creates a new dynamic slice item.
func NewDynamicSliceItemImpl[T any](
	field string,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *DynamicSliceItemImpl[T] {
	return &DynamicSliceItemImpl[T]{
		field:     field,
		parse:     parse,
		format:    format,
		delimiter: delimiter,
		values:    make(map[string][]T),
	}
}

// Set parses and appends values for a given instance ID.
func (d *DynamicSliceItemImpl[T]) Set(id, input string) error {
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

// Get returns the slice for the given instance ID.
func (d *DynamicSliceItemImpl[T]) Get(id string) ([]T, bool) {
	val, ok := d.values[id]
	return val, ok
}

// Values returns the full map of instanceID → []T.
func (d *DynamicSliceItemImpl[T]) Values() map[string][]T {
	return d.values
}

// ValuesAny returns parsed values as map[string]any.
func (d *DynamicSliceItemImpl[T]) ValuesAny() map[string]any {
	m := make(map[string]any, len(d.values))
	for id, vals := range d.values {
		m[id] = vals
	}
	return m
}

// SetValidator sets an optional validator for slice elements.
func (d *DynamicSliceItemImpl[T]) SetValidator(fn func(T) error) {
	d.validator = fn
}

// SetDelimiter changes the delimiter used when splitting input strings.
func (d *DynamicSliceItemImpl[T]) SetDelimiter(delim string) {
	d.delimiter = delim
}
