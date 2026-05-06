package core

import (
	"github.com/containeroo/tinyflags/internal/utils"
)

// ValueHooks centralizes parse/format/validate/finalize behavior for typed values.
type ValueHooks[T any] struct {
	Parse            func(string) (T, error)
	Format           func(T) string
	Validate         func(T) error
	Finalize         func(T) T
	FinalizeDefault  bool
	DefaultFinalized bool
}

// NewValueHooks returns a new hook container for a typed value.
func NewValueHooks[T any](parse func(string) (T, error), format func(T) string) ValueHooks[T] {
	return ValueHooks[T]{
		Parse:  parse,
		Format: format,
	}
}

// ParseValue parses and applies validation/finalization hooks.
func (h *ValueHooks[T]) ParseValue(raw string) (T, error) {
	val, err := h.Parse(raw)
	if err != nil {
		return val, err
	}
	return utils.ApplyValueHooks(val, h.Validate, h.Finalize)
}

// DefaultString returns the formatted default string.
func (h *ValueHooks[T]) DefaultString(def T) string {
	return h.Format(def)
}

// SetValidate installs a validation hook.
func (h *ValueHooks[T]) SetValidate(fn func(T) error) {
	h.Validate = fn
}

// SetFinalize installs a finalization hook.
func (h *ValueHooks[T]) SetFinalize(fn func(T) T) {
	h.Finalize = fn
}

// EnableFinalizeDefault enables running the finalizer on defaults when unset.
func (h *ValueHooks[T]) EnableFinalizeDefault() {
	h.FinalizeDefault = true
}

// ApplyDefaultScalar finalizes the scalar default when eligible.
func (h *ValueHooks[T]) ApplyDefaultScalar(current *T, changed bool) {
	utils.ApplyDefaultValueFinalize(current, changed, &h.DefaultFinalized, h.FinalizeDefault, h.Finalize)
}

// ApplyDefaultSlice finalizes each slice default item when eligible.
func (h *ValueHooks[T]) ApplyDefaultSlice(items []T, changed bool) {
	utils.ApplyDefaultSliceFinalize(items, changed, &h.DefaultFinalized, h.FinalizeDefault, h.Finalize)
}

// ResetDefaultFinalize clears the default-finalized lifecycle bit.
func (h *ValueHooks[T]) ResetDefaultFinalize() {
	h.DefaultFinalized = false
}
