package tinyflags

import "fmt"

// Flag is the user-facing scalar flag builder.
type Flag[T any] struct {
	builderImpl[T] // scalar flag builder logic
}

// Choices restricts allowed scalar values.
func (f *Flag[T]) Choices(allowed ...T) *Flag[T] {
	f.value.SetValidator(func(v T) error {
		for _, a := range allowed {
			if f.value.format(a) == f.value.format(v) {
				return nil
			}
		}
		return fmt.Errorf("must be one of %s", formatAllowed(allowed, f.value.format))
	})
	f.bf.allowed = make([]string, len(allowed))
	for i, a := range allowed {
		f.bf.allowed[i] = f.value.format(a)
	}
	return f
}
