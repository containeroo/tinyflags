package tinyflags

import (
	"fmt"
	"strings"
)

type SliceFlagValue[T any] struct {
	*FlagValue[[]T]
	delimiter string
	parse     func(string) (T, error)
	format    func(T) string
	validator func(T) bool
	allowed   []T
}

func NewSliceFlagValue[T any](
	ptr *[]T,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
) *SliceFlagValue[T] {
	*ptr = append([]T{}, def...)
	return &SliceFlagValue[T]{
		FlagValue: NewFlagValue(ptr, def, nil, nil),
		delimiter: delimiter,
		parse:     parse,
		format:    format,
	}
}

func (s *SliceFlagValue[T]) Set(input string) error {
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

func (s *SliceFlagValue[T]) Default() string {
	var out []string
	for _, item := range s.def {
		out = append(out, s.format(item))
	}
	return strings.Join(out, s.delimiter)
}

func (s *SliceFlagValue[T]) SetDelimiter(d string) {
	s.delimiter = d
}

func (s *SliceFlagValue[T]) SetValidator(fn func(T) bool, allowed []T) {
	s.validator = fn
	s.allowed = allowed
}

func (s *SliceFlagValue[T]) isSlice()              {}
func (s *SliceFlagValue[T]) Base() *FlagValue[[]T] { return s.FlagValue }
