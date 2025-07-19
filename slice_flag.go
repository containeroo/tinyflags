package tinyflags

// SliceFlag is the user-facing builder for slice flags.
type SliceFlag[T any] struct {
	Flag[[]T] // inherits scalar builder methods
}

// Delimiter configures the input separator for slice values.
func (s *SliceFlag[T]) Delimiter(sep string) *SliceFlag[T] {
	if d, ok := s.bf.value.(HasDelimiter); ok {
		d.SetDelimiter(sep)
	}
	return s
}

// Choices restricts allowed slice elements.
func (s *SliceFlag[T]) Choices(allowed ...T) *SliceFlag[T] {
	sv, ok := any(s.bf.value).(*SliceValueImpl[T])
	if !ok {
		return s
	}
	sv.SetValidator(func(val T) bool {
		for _, a := range allowed {
			if sv.format(a) == sv.format(val) {
				return true
			}
		}
		return false
	}, allowed)

	s.bf.allowed = make([]string, len(allowed))
	for i, a := range allowed {
		s.bf.allowed[i] = sv.format(a)
	}
	return s
}
