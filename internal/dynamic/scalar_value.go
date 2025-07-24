package dynamic

// DynamicScalarValue parses and stores one value per ID.
type DynamicScalarValue[T any] struct {
	field    string
	def      T
	parse    func(string) (T, error)
	format   func(T) string
	validate func(T) error
	values   map[string]T
}

// NewDynamicScalarValue creates a per-ID scalar parser.
func NewDynamicScalarValue[T any](
	field string,
	def T,
	parse func(string) (T, error),
	format func(T) string,
) *DynamicScalarValue[T] {
	return &DynamicScalarValue[T]{
		field:  field,
		def:    def,
		parse:  parse,
		format: format,
		values: map[string]T{},
	}
}

// Set parses and stores one entry.
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
	d.values[id] = val
	return nil
}

// setValidate sets a per-item validation function.
func (d *DynamicScalarValue[T]) setValidate(fn func(T) error) {
	d.validate = fn
}

func (d *DynamicScalarValue[T]) ValuesAny() map[string]any {
	m := make(map[string]any, len(d.values))
	for k, v := range d.values {
		m[k] = v
	}
	return m
}
