package tinyflags

import "fmt"

// DynamicFlag is a builder for a dynamic scalar flag.
// Each instance (e.g. `--http.alpha.port`) can hold one value.
type DynamicFlag[T any] struct {
	fs   *FlagSet            // parent flag set
	bf   *baseFlag           // metadata and registration
	item *DynamicItemImpl[T] // actual value store
}

// Get retrieves the value for the given instance ID.
func (d *DynamicFlag[T]) Get(id string) (T, bool) {
	return d.item.Get(id)
}

// MustGet retrieves the value or panics if not set.
func (d *DynamicFlag[T]) MustGet(id string) T {
	val, ok := d.item.Get(id)
	if !ok {
		panic(fmt.Sprintf("value for id %q not set", id))
	}
	return val
}

// Values returns a map of all instance IDs and their values.
func (d *DynamicFlag[T]) Values() map[string]T {
	return d.item.Values()
}

// ValuesAny returns a map[string]any for introspection/debug.
func (d *DynamicFlag[T]) ValuesAny() map[string]any {
	return d.item.ValuesAny()
}
