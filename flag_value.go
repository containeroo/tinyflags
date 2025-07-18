package tinyflags

type FlagValue[T any] struct {
	ptr      *T
	def      T
	value    T
	changed  bool
	parse    func(string) (T, error)
	format   func(T) string
	validate func(T) error
}

func NewFlagValue[T any](ptr *T, def T, parse func(string) (T, error), format func(T) string) *FlagValue[T] {
	*ptr = def
	return &FlagValue[T]{
		ptr:    ptr,
		def:    def,
		parse:  parse,
		format: format,
	}
}

func (f *FlagValue[T]) Set(s string) error {
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

func (f *FlagValue[T]) Get() any                      { return *f.ptr }
func (f *FlagValue[T]) Default() string               { return f.format(f.def) }
func (f *FlagValue[T]) IsChanged() bool               { return f.changed }
func (f *FlagValue[T]) SetValidator(fn func(T) error) { f.validate = fn }
func (f *FlagValue[T]) Base() *FlagValue[T]           { return f }
