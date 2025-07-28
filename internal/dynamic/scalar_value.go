package dynamic

type DynamicScalarValue[T any] struct {
	field    string
	def      T
	parse    func(string) (T, error)
	format   func(T) string
	validate func(T) error
	values   map[string]T
}

func NewDynamicScalarValue[T any](field string, def T, parse func(string) (T, error), format func(T) string) *DynamicScalarValue[T] {
	return &DynamicScalarValue[T]{field: field, def: def, parse: parse, format: format, values: make(map[string]T)}
}

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

func (d *DynamicScalarValue[T]) FieldName() string { return d.field }
func (d *DynamicScalarValue[T]) GetAny(id string) (any, bool) {
	val, ok := d.values[id]
	if ok {
		return val, true
	}
	return d.def, false
}

func (d *DynamicScalarValue[T]) ValuesAny() map[string]any {
	out := make(map[string]any)
	for k, v := range d.values {
		out[k] = v
	}
	return out
}
func (d *DynamicScalarValue[T]) setValidate(fn func(T) error) { d.validate = fn }
func (d *DynamicScalarValue[T]) Base() *DynamicScalarValue[T] { return d }
