package scalar

// ScalarValue implements scalar flag parsing, formatting, and validation.
type ScalarValue[T any] struct {
	Ptr      *T
	def      T
	value    T
	changed  bool
	parse    func(string) (T, error)
	format   func(T) string
	validate func(T) error
}

func NewScalarValue[T any](ptr *T, def T, parse func(string) (T, error), format func(T) string) *ScalarValue[T] {
	*ptr = def
	return &ScalarValue[T]{
		Ptr:    ptr,
		def:    def,
		parse:  parse,
		format: format,
	}
}

func (f *ScalarValue[T]) Set(s string) error {
	val, err := f.parse(s)
	if err != nil {
		return err
	}
	if f.validate != nil {
		if err := f.validate(val); err != nil {
			return err
		}
	}
	*f.Ptr = val
	f.value = val
	f.changed = true
	return nil
}

func (f *ScalarValue[T]) Get() any                     { return *f.Ptr }
func (f *ScalarValue[T]) Default() string              { return f.format(f.def) }
func (f *ScalarValue[T]) Changed() bool                { return f.changed }
func (f *ScalarValue[T]) setValidate(fn func(T) error) { f.validate = fn }
func (f *ScalarValue[T]) Base() *ScalarValue[T]        { return f }
