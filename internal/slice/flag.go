package slice

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// SliceFlag is the user-facing builder for slice flags.
type SliceFlag[T any] struct {
	builder.DefaultFlag[[]T]
	val *SliceValue[T]
}

// Delimiter sets the delimiter used to split input values.
func (f *SliceFlag[T]) Delimiter(sep string) *SliceFlag[T] {
	f.val.delimiter = sep
	return f
}

// Choices restricts allowed slice elements.
func (f *SliceFlag[T]) Choices(allowed ...T) *SliceFlag[T] {
	f.val.setValidate(func(v T) error {
		for _, a := range allowed {
			if f.val.format(a) == f.val.format(v) {
				return nil
			}
		}
		return fmt.Errorf("must be one of [%s]", utils.FormatAllowed(allowed, f.val.format))
	})
	f.BF.Allowed = make([]string, len(allowed))
	for i, a := range allowed {
		f.BF.Allowed[i] = f.val.format(a)
	}
	return f
}

func (f *SliceFlag[T]) Validate(fn func(T) error) *SliceFlag[T] {
	f.val.setValidate(fn)
	return f
}
