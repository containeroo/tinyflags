package dynamic

import (
	"fmt"
	"strings"
)

type DynamicSliceValue[T any] struct {
	field     string
	def       []T
	parse     func(string) (T, error)
	format    func(T) string
	delimiter string
	validate  func(T) error
	values    map[string][]T
}

func NewDynamicSliceValue[T any](field string, def []T, parse func(string) (T, error), format func(T) string, delimiter string) *DynamicSliceValue[T] {
	return &DynamicSliceValue[T]{field: field, def: def, parse: parse, format: format, delimiter: delimiter, values: make(map[string][]T)}
}

func (d *DynamicSliceValue[T]) Set(id, raw string) error {
	for _, chunk := range strings.Split(raw, d.delimiter) {
		val, err := d.parse(strings.TrimSpace(chunk))
		if err != nil {
			return fmt.Errorf("invalid %q: %w", chunk, err)
		}
		if d.validate != nil {
			if err := d.validate(val); err != nil {
				return fmt.Errorf("invalid value %q: %w", chunk, err)
			}
		}
		d.values[id] = append(d.values[id], val)
	}
	return nil
}

func (d *DynamicSliceValue[T]) FieldName() string { return d.field }
func (d *DynamicSliceValue[T]) GetAny(id string) (any, bool) {
	val, ok := d.values[id]
	if ok {
		return val, true
	}
	return d.def, false
}
func (d *DynamicSliceValue[T]) ValuesAny() map[string]any {
	out := make(map[string]any)
	for k, v := range d.values {
		out[k] = v
	}
	return out
}
func (d *DynamicSliceValue[T]) setValidate(fn func(T) error) { d.validate = fn }
func (d *DynamicSliceValue[T]) setDelimiter(sep string)      { d.delimiter = sep }
