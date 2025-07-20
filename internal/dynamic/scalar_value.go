package dynamic

// DynamicScalarValue holds a scalar value per instance ID.
type DynamicScalarValue[T any] struct {
	field     string
	parse     func(string) (T, error)
	format    func(T) string
	validator func(T) error
	values    map[string]T
}

// NewDynamicScalarValue constructs a new dynamic scalar value handler.
func NewDynamicScalarValue[T any](field string, parse func(string) (T, error), format func(T) string) *DynamicScalarValue[T] {
	return &DynamicScalarValue[T]{
		field:  field,
		parse:  parse,
		format: format,
		values: make(map[string]T),
	}
}

// Set parses and stores a value for a given instance.
func (d *DynamicScalarValue[T]) Set(id, raw string) error {
	val, err := d.parse(raw)
	if err != nil {
		return err
	}
	if d.validator != nil {
		if err := d.validator(val); err != nil {
			return err
		}
	}
	d.values[id] = val
	return nil
}

// Get retrieves the parsed value.
func (d *DynamicScalarValue[T]) Get(id string) (T, bool) {
	val, ok := d.values[id]
	return val, ok
}

// Values returns all stored values.
func (d *DynamicScalarValue[T]) Values() map[string]T {
	return d.values
}

// ValuesAny returns values as a generic map.
func (d *DynamicScalarValue[T]) ValuesAny() map[string]any {
	m := make(map[string]any, len(d.values))
	for k, v := range d.values {
		m[k] = v
	}
	return m
}

// SetValidator sets a validation function for the values.
func (d *DynamicScalarValue[T]) setValidate(fn func(T) error) {
	d.validator = fn
}
