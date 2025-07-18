package tinyflags

// ValueImpl implements scalar flag parsing, formatting, and validation.
type ValueImpl[T any] struct {
	ptr      *T
	def      T
	value    T
	changed  bool
	parse    func(string) (T, error)
	format   func(T) string
	validate func(T) error
}

func NewValueImpl[T any](ptr *T, def T, parse func(string) (T, error), format func(T) string) *ValueImpl[T] {
	*ptr = def
	return &ValueImpl[T]{
		ptr:    ptr,
		def:    def,
		parse:  parse,
		format: format,
	}
}

func (f *ValueImpl[T]) Set(s string) error {
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

func (f *ValueImpl[T]) Get() any                      { return *f.ptr }
func (f *ValueImpl[T]) Default() string               { return f.format(f.def) }
func (f *ValueImpl[T]) IsChanged() bool               { return f.changed }
func (f *ValueImpl[T]) SetValidator(fn func(T) error) { f.validate = fn }
func (f *ValueImpl[T]) Base() *ValueImpl[T]           { return f }
