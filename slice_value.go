package tinyflags

import (
	"fmt"
	"strings"
)

// SliceValueImpl implements slice value parsing and validation.
type SliceValueImpl[T any] struct {
	*ScalarValueImpl[[]T]
	delimiter string
	parse     func(string) (T, error)
	format    func(T) string
	validator func(T) bool
	allowed   []T
}

func NewSliceValueImpl[T any](
	ptr *[]T,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *SliceValueImpl[T] {
	*ptr = append([]T{}, def...)
	return &SliceValueImpl[T]{
		ScalarValueImpl: NewScalarValueImpl(ptr, def, nil, nil),
		delimiter:       delimiter,
		parse:           parse,
		format:          format,
	}
}

func (s *SliceValueImpl[T]) Set(input string) error {
	if !s.changed {
		*s.ptr = nil
	}
	for _, p := range strings.Split(input, s.delimiter) {
		val, err := s.parse(strings.TrimSpace(p))
		if err != nil {
			return fmt.Errorf("invalid slice item %q: %w", p, err)
		}
		if s.validator != nil && !s.validator(val) {
			return fmt.Errorf("invalid value %q: must be one of [%s]",
				s.format(val), formatAllowed(s.allowed, s.format))
		}
		*s.ptr = append(*s.ptr, val)
	}
	s.value = *s.ptr
	s.changed = true
	return nil
}

func (s *SliceValueImpl[T]) Default() string {
	var out []string
	for _, item := range s.def {
		out = append(out, s.format(item))
	}
	return strings.Join(out, s.delimiter)
}

func (s *SliceValueImpl[T]) SetDelimiter(d string) {
	s.delimiter = d
}

func (s *SliceValueImpl[T]) SetValidator(fn func(T) bool, allowed []T) {
	s.validator = fn
	s.allowed = allowed
}

func (s *SliceValueImpl[T]) isSlice()                    {}
func (s *SliceValueImpl[T]) Base() *ScalarValueImpl[[]T] { return s.ScalarValueImpl }
