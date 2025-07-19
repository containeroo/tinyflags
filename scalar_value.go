package tinyflags

// ScalarValueImpl implements scalar flag parsing, formatting, and validation.
type ScalarValueImpl[T any] struct {
	ptr      *T
	def      T
	value    T
	changed  bool
	parse    func(string) (T, error)
	format   func(T) string
	validate func(T) error
}

func NewScalarValueImpl[T any](ptr *T, def T, parse func(string) (T, error), format func(T) string) *ScalarValueImpl[T] {
	*ptr = def
	return &ScalarValueImpl[T]{
		ptr:    ptr,
		def:    def,
		parse:  parse,
		format: format,
	}
}

func (f *ScalarValueImpl[T]) Set(s string) error {
	val, err := f.parse(s)
	if err != nil {
		return err
	}
	if f.validate != nil {
		if err := f.validate(val); err != nil {
			return err
		}
	}
	*f.ptr = val
	f.value = val
	f.changed = true
	return nil
}

func (f *ScalarValueImpl[T]) Get() any                      { return *f.ptr }
func (f *ScalarValueImpl[T]) Default() string               { return f.format(f.def) }
func (f *ScalarValueImpl[T]) IsChanged() bool               { return f.changed }
func (f *ScalarValueImpl[T]) SetValidator(fn func(T) error) { f.validate = fn }
func (f *ScalarValueImpl[T]) Base() *ScalarValueImpl[T]     { return f }
