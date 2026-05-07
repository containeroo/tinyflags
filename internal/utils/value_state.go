package utils

// ApplyDefaultValueFinalize applies the default-only finalizer to a scalar value when eligible.
func ApplyDefaultValueFinalize[T any](current *T, changed bool, defaultFinalized *bool, finalizeDefault bool, finalize func(T) T) {
	if changed || *defaultFinalized || !finalizeDefault || finalize == nil {
		return
	}
	*current = finalize(*current)
	*defaultFinalized = true
}

// ApplyDefaultSliceFinalize applies the default-only finalizer to each item in a slice.
func ApplyDefaultSliceFinalize[T any](items []T, changed bool, defaultFinalized *bool, finalizeDefault bool, finalize func(T) T) {
	if changed || *defaultFinalized || !finalizeDefault || finalize == nil {
		return
	}
	for i, item := range items {
		items[i] = finalize(item)
	}
	*defaultFinalized = true
}

// ResetScalarState restores a scalar to its default value and clears parse lifecycle markers.
func ResetScalarState[T any](current *T, def T, changed *bool, defaultFinalized *bool) {
	*current = def
	*changed = false
	*defaultFinalized = false
}

// ResetSliceState restores a slice to its default values and clears parse lifecycle markers.
func ResetSliceState[T any](current *[]T, def []T, changed *bool, defaultFinalized *bool) {
	*current = append((*current)[:0], def...)
	*changed = false
	*defaultFinalized = false
}
