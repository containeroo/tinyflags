// internal/scalar/flag.go
package scalar

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/utils"
)

// ScalarFlag is the user-facing scalar flag builder.
type ScalarFlag[T any] struct {
	builder.StaticFlag[T]
	val *ScalarValue[T]
}

// Choices restricts allowed scalar values.
func (f *ScalarFlag[T]) Choices(allowed ...T) *ScalarFlag[T] {
	f.val.setValidate(func(v T) error {
		for _, a := range allowed {
			if f.val.format(a) == f.val.format(v) {
				return nil
			}
		}
		return fmt.Errorf("must be one of %s", utils.FormatAllowed(allowed, f.val.format))
	})

	formatted := make([]string, len(allowed))
	for i, a := range allowed {
		formatted[i] = f.val.format(a)
	}

	f.Allowed(formatted...)

	return f
}

// Validate lets you plug in arbitrary per‐element checks.
func (f *ScalarFlag[T]) Validate(fn func(T) error) *ScalarFlag[T] {
	f.val.setValidate(fn)
	return f
}
