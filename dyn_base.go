package tinyflags

// DynamicItemImpl stores a single value per instance ID.
type DynamicItemImpl[T any] struct {
	field     string                  // field name (e.g. "port")
	parse     func(string) (T, error) // parses string input into T
	format    func(T) string          // formats T as string
	validator func(T) error           // optional per-value validation
	values    map[string]T            // instanceID → value
}

// NewDynamicItemImpl constructs a new dynamic scalar item.
func NewDynamicItemImpl[T any](field string, parse func(string) (T, error), format func(T) string) *DynamicItemImpl[T] {
	return &DynamicItemImpl[T]{
		field:  field,
		parse:  parse,
		format: format,
		values: make(map[string]T),
	}
}

// Set parses and stores a value for a given instance ID.
func (d *DynamicItemImpl[T]) Set(id string, input string) error {
	val, err := d.parse(input)
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

// Get returns the parsed value for a given instance ID.
func (d *DynamicItemImpl[T]) Get(id string) (T, bool) {
	val, ok := d.values[id]
	return val, ok
}

// Values returns all parsed values as map[instanceID]T.
func (d *DynamicItemImpl[T]) Values() map[string]T {
	return d.values
}

// ValuesAny returns all parsed values as map[instanceID]any.
func (d *DynamicItemImpl[T]) ValuesAny() map[string]any {
	m := make(map[string]any, len(d.values))
	for k, v := range d.values {
		m[k] = v
	}
	return m
}

// SetValidator assigns a validation function for parsed values.
func (d *DynamicItemImpl[T]) SetValidator(fn func(T) error) {
	d.validator = fn
}
