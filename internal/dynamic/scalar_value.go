package dynamic

// DynamicScalarValue holds parsed scalar values per ID with parsing, formatting, and validation.
type DynamicScalarValue[T any] struct {
	field    string                  // Flag field name
	def      T                       // Default value
	parse    func(string) (T, error) // Parser from raw input
	format   func(T) string          // Formatter to string
	validate func(T) error           // Optional validation function
	finalize (func(T) T)             // Optional finalizer function
	values   map[string]T            // Parsed values per ID
}

// NewDynamicScalarValue creates a new dynamic scalar value.
func NewDynamicScalarValue[T any](field string, def T, parse func(string) (T, error), format func(T) string) *DynamicScalarValue[T] {
	return &DynamicScalarValue[T]{
		field:  field,
		def:    def,
		parse:  parse,
		format: format,
		values: make(map[string]T),
	}
}

// Set parses and stores a value for a specific ID.
func (d *DynamicScalarValue[T]) Set(id, raw string) error {
	val, err := d.parse(raw)
	if err != nil {
		return err
	}
	if d.validate != nil {
		if err := d.validate(val); err != nil {
			return err
		}
	}
	if d.finalize != nil {
		val = d.finalize(val)
	}
	d.values[id] = val
	return nil
}

// setValidate sets the optional validation function.
func (d *DynamicScalarValue[T]) setValidate(fn func(T) error) {
	d.validate = fn
}

// setFinalize sets the optional finalizer function.
func (d *DynamicScalarValue[T]) setFinalize(fn func(T) T) {
	d.finalize = fn
}

// Base returns itself for generic access.
func (d *DynamicScalarValue[T]) Base() *DynamicScalarValue[T] {
	return d
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
