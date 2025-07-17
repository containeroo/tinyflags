package tinyflags

type FlagBase[T any] struct {
	ptr      *T
	def      T
	value    T
	changed  bool
	parse    func(string) (T, error)
	format   func(T) string
	validate func(T) error
}

func NewFlagBase[T any](ptr *T, def T, parse func(string) (T, error), format func(T) string) *FlagBase[T] {
	*ptr = def
	return &FlagBase[T]{
		ptr:     ptr,
		def:     def,
		parse:   parse,
		format:  format,
		changed: false,
	}
}

func (f *FlagBase[T]) Set(s string) error {
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

func (f *FlagBase[T]) Get() any {
	return *f.ptr
}

func (f *FlagBase[T]) Default() string {
	return f.format(f.def)
}

func (f *FlagBase[T]) IsChanged() bool {
	return f.changed
}

func (f *FlagBase[T]) SetValidator(v func(T) error) {
	f.validate = v
}
