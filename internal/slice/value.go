package slice

import (
	"fmt"
	"strings"
)

// SliceValue implements slice flag parsing and validation.
type SliceValue[T any] struct {
	ptr       *[]T
	def       []T
	value     []T
	changed   bool
	delimiter string
	parse     func(string) (T, error)
	format    func(T) string
	validate  func(T) error
}

func NewSliceValue[T any](
	ptr *[]T,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *SliceValue[T] {
	*ptr = append([]T{}, def...)
	return &SliceValue[T]{
		ptr:       ptr,
		def:       def,
		delimiter: delimiter,
		parse:     parse,
		format:    format,
	}
}

func (f *SliceValue[T]) Set(s string) error {
	if !f.changed {
		*f.ptr = nil
	}
	parts := strings.Split(s, f.delimiter)
	for _, raw := range parts {
		val, err := f.parse(strings.TrimSpace(raw))
		if err != nil {
			return fmt.Errorf("invalid slice item %q: %w", raw, err)
		}
		if f.validate != nil {
			if err := f.validate(val); err != nil {
				return fmt.Errorf("invalid value %q: %w", raw, err)
			}
		}
		*f.ptr = append(*f.ptr, val)
	}
	f.value = *f.ptr
	f.changed = true
	return nil
}

func (f *SliceValue[T]) Get() any {
	return *f.ptr
}

func (f *SliceValue[T]) DefaultString() string {
	out := make([]string, 0, len(f.def))
	for _, v := range f.def {
		out = append(out, f.format(v))
	}
	return strings.Join(out, f.delimiter)
}

func (f *SliceValue[T]) Changed() bool {
	return f.changed
}

func (f *SliceValue[T]) setValidate(fn func(T) error) { f.validate = fn }

func (s *SliceValue[T]) isSlice()             {}
func (f *SliceValue[T]) Base() *SliceValue[T] { return f }

func (f *SliceValue[T]) Default() string {
	out := make([]string, 0, len(f.def))
	for _, v := range f.def {
		out = append(out, f.format(v))
	}
	return strings.Join(out, f.delimiter)
}
