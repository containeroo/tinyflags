package tinyflags

type DynamicItem[T any] struct {
	field     string
	parse     func(string) (T, error)
	format    func(T) string
	validator func(T) error
	values    map[string]T
}

func NewDynamicItem[T any](field string, parse func(string) (T, error), format func(T) string) *DynamicItem[T] {
	return &DynamicItem[T]{
		field:  field,
		parse:  parse,
		format: format,
		values: make(map[string]T),
	}
}

func (d *DynamicItem[T]) Set(id string, input string) error {
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

func (d *DynamicItem[T]) Get(id string) (T, bool) {
	val, ok := d.values[id]
	return val, ok
}

func (d *DynamicItem[T]) Values() map[string]T {
	return d.values
}

func (d *DynamicItem[T]) ValuesAny() map[string]any {
	m := make(map[string]any, len(d.values))
	for k, v := range d.values {
		m[k] = v
	}
	return m
}

func (d *DynamicItem[T]) SetValidator(fn func(T) error) {
	d.validator = fn
}
