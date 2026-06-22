package slice

import (
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/utils"
)

// SliceValue implements slice flag parsing and validation.
type SliceValue[T any] struct {
	ptr     *[]T
	def     []T
	changed bool
	input   core.SliceInputConfig
	hooks   core.ValueHooks[T]
}

// NewSliceValue creates a new slice value.
func NewSliceValue[T any](
	ptr *[]T,
	def []T,
	parse func(string) (T, error),
	format func(T) string,
	delimiter string,
	trimSpace bool,
) *SliceValue[T] {
	*ptr = append([]T{}, def...)
	return &SliceValue[T]{
		ptr:   ptr,
		def:   def,
		input: core.SliceInputConfig{Delimiter: delimiter, TrimSpace: trimSpace},
		hooks: core.NewValueHooks(parse, format),
	}
}

// Set parses and stores the slice from a delimited string for a given ID.
func (v *SliceValue[T]) Set(s string) error {
	if !v.changed {
		*v.ptr = nil
	}
	parts, err := v.input.Split(s)
	if err != nil {
		return err
	}
	for _, raw := range parts {
		raw = v.input.Normalize(raw)
		if raw == "" && !v.input.AllowEmpty {
			return fmt.Errorf("invalid slice item %q: empty values are not allowed", raw)
		}
		val, err := v.hooks.ParseValue(raw)
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
		out = append(out, f.hooks.Format(v))
	}
	return strings.Join(out, f.input.Delimiter)
}

// Changed returns true if the value was changed.
func (v *SliceValue[T]) Changed() bool {
	return v.changed
}

// setValidate sets a per-item validation function.
func (v *SliceValue[T]) setValidate(fn func(T) error) { v.hooks.SetValidate(fn) }

// setFinalize sets a per-item finalizer function.
func (v *SliceValue[T]) setFinalize(fn func(T) T) { v.hooks.SetFinalize(fn) }

// setFinalizeDefaultValue enables running the finalizer on defaults when unset.
func (v *SliceValue[T]) setFinalizeDefaultValue() { v.hooks.EnableFinalizeDefault() }

// setAllowEmpty toggles acceptance of empty items.
func (v *SliceValue[T]) setAllowEmpty(allow bool) { v.input.AllowEmpty = allow }

// setTrimSpace toggles trimming leading and trailing space from each item.
func (v *SliceValue[T]) setTrimSpace(trim bool) { v.input.TrimSpace = trim }

// Base returns the underlying value.
func (v *SliceValue[T]) Base() *SliceValue[T] { return v }

// ApplyDefaultFinalize applies the default-only finalizer when unset.
func (v *SliceValue[T]) ApplyDefaultFinalize() {
	v.hooks.ApplyDefaultSlice(*v.ptr, v.changed)
}

// isSlice is a no-op marker method to implement core.SliceMarker.
func (v *SliceValue[T]) IsSlice() {}

// ResetParseState restores the default slice and clears changed/finalized state.
func (v *SliceValue[T]) ResetParseState() {
	utils.ResetSliceState(v.ptr, v.def, &v.changed, &v.hooks.DefaultFinalized)
}
