package tinyflags

import "fmt"

// DynamicFlag defines a builder for dynamic flags keyed by an instance ID (e.g., "alpha", "beta").
type DynamicFlag[T any] struct {
	builderImpl[T]                 // Inherits builder methods like .Required(), .Group(), etc.
	item           *DynamicItem[T] // Tracks all values by identifier
}

// Get retrieves the parsed value for a given instance ID, or false if not set.
func (d *DynamicFlag[T]) Get(id string) (T, bool) {
	return d.item.Get(id)
}

// MustGet retrieves the value or panics if not set. Useful in tests or trusted code paths.
func (d *DynamicFlag[T]) MustGet(id string) T {
	val, ok := d.item.Get(id)
	if !ok {
		panic(fmt.Sprintf("value for id %q not set", id))
	}
	return val
}

// Values returns a map of all parsed values keyed by instance ID.
func (d *DynamicFlag[T]) Values() map[string]T {
	return d.item.Values()
}

// ValuesAny returns parsed values as map[string]any.
func (d *DynamicFlag[T]) ValuesAny() map[string]any {
	return d.item.ValuesAny()
}
