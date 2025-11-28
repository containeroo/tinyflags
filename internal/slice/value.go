package slice

import (
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/utils"
)

// SliceValue implements slice flag parsing and validation.
type SliceValue[T any] struct {
	ptr        *[]T
	def        []T
	changed    bool
	delimiter  string
	strictDel  bool
	allowEmpty bool
	parse      func(string) (T, error)
	format     func(T) string
	validate   func(T) error
	finalize   (func(T) T)
}

// NewSliceValue creates a new slice value.
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

// Set parses and stores the slice from a delimited string for a given ID.
func (v *SliceValue[T]) Set(s string) error {
	if !v.changed {
		*v.ptr = nil
	}
	parts := strings.Split(s, v.delimiter)
	if v.strictDel {
		for _, alt := range []string{",", ";", "|"} {
			if alt == v.delimiter {
				continue
			}
			if strings.Contains(s, alt) {
				return fmt.Errorf("mixed delimiters: found %q while using %q", alt, v.delimiter)
			}
		}
	}
	for _, raw := range parts {
		raw = strings.TrimSpace(raw)
		if raw == "" && !v.allowEmpty {
			return fmt.Errorf("invalid slice item %q: empty values are not allowed", raw)
		}
		val, err := v.parse(raw)
		if err != nil {
			return fmt.Errorf("invalid slice item %q: %w", raw, err)
		}
		val, err = utils.ApplyValueHooks(val, v.validate, v.finalize)
		if err != nil {
			return fmt.Errorf("invalid value %q: %w", raw, err)
		}
		*v.ptr = append(*v.ptr, val)
	}
	v.changed = true
	return nil
}

// Get returns the parsed slice for the given ID.
func (v *SliceValue[T]) Get() any {
	return *v.ptr
}

// Default returns the default value as string.
func (f *SliceValue[T]) Default() string {
	out := make([]string, 0, len(f.def))
	for _, v := range f.def {
		out = append(out, f.format(v))
	}
	return strings.Join(out, f.delimiter)
}

// Changed returns true if the value was changed.
func (v *SliceValue[T]) Changed() bool {
	return v.changed
}

// setValidate sets a per-item validation function.
func (v *SliceValue[T]) setValidate(fn func(T) error) { v.validate = fn }

// setFinalize sets a per-item finalizer function.
func (v *SliceValue[T]) setFinalize(fn func(T) T) { v.finalize = fn }

// setStrictDelimiter toggles mixed-delimiter rejection.
func (v *SliceValue[T]) setStrictDelimiter(strict bool) { v.strictDel = strict }

// setAllowEmpty toggles acceptance of empty items.
func (v *SliceValue[T]) setAllowEmpty(allow bool) { v.allowEmpty = allow }

// Base returns the underlying value.
func (v *SliceValue[T]) Base() *SliceValue[T] { return v }

// isSlice is a no-op marker method to implement core.SliceMarker.
func (v *SliceValue[T]) IsSlice() {}
