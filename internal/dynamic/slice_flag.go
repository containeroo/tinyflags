package dynamic

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/utils"
)

// SliceFlag provides the user-facing API for a dynamic slice flag,
// e.g., --http.alpha.tags=a,b or --node.node1.labels=env,prod.
type SliceFlag[T any] struct {
	registry core.Registry
	bf       *core.BaseFlag        // optional pointer to the base flag for help and metadata
	item     *DynamicSliceValue[T] // value storage and parsing logic
}

// Get returns the parsed slice for the given instance ID, or false if not set.
func (f *SliceFlag[T]) Get(id string) ([]T, bool) {
	return f.item.Get(id)
}

// MustGet returns the slice for the given instance ID or panics if not set.
func (f *SliceFlag[T]) MustGet(id string) []T {
	val, ok := f.item.Get(id)
	if !ok {
		panic(fmt.Sprintf("value for id %q not set", id))
	}
	return val
}

// Values returns all parsed slice values keyed by instance ID.
func (f *SliceFlag[T]) Values() map[string][]T {
	return f.item.Values()
}

// ValuesAny returns all values in a map[string]any form.
func (f *SliceFlag[T]) ValuesAny() map[string]any {
	return f.item.ValuesAny()
}

// Delimiter sets a custom delimiter (e.g., ";" or "|") for the slice flag.
func (f *SliceFlag[T]) Delimiter(sep string) *SliceFlag[T] {
	f.item.setDelimiter(sep)
	return f
}

// Choices restricts each element in the slice to one of the allowed values.
func (f *SliceFlag[T]) Choices(allowed ...T) *SliceFlag[T] {
	f.item.setValidate(func(val T) error {
		for _, a := range allowed {
			if f.item.format(a) == f.item.format(val) {
				return nil
			}
		}
		return fmt.Errorf("must be one of [%s]", utils.FormatAllowed(allowed, f.item.format))
	})
	// Populate for help output
	f.bf.Allowed = make([]string, len(allowed))
	for i, a := range allowed {
		f.bf.Allowed[i] = f.item.format(a)
	}
	return f
}

// Validate sets a custom validation function for each slice item.
func (f *SliceFlag[T]) Validate(fn func(T) error) *SliceFlag[T] {
	f.item.setValidate(fn)
	return f
}

// Required marks this flag as required.
func (f *SliceFlag[T]) Required() *SliceFlag[T] {
	f.bf.Required = true
	return f
}

// Hidden marks this flag as hidden from help output.
func (f *SliceFlag[T]) Hidden() *SliceFlag[T] {
	f.bf.Hidden = true
	return f
}

// Deprecated marks this flag as deprecated with an optional reason.
func (f *SliceFlag[T]) Deprecated(reason string) *SliceFlag[T] {
	f.bf.Deprecated = reason
	return f
}

// Group adds this flag to a mutual exclusion group.
func (f *SliceFlag[T]) Group(name string) *SliceFlag[T] {
	if name == "" {
		return f
	}
	group := f.registry.GetGroup(name)
	group.Flags = append(group.Flags, f.bf)
	f.bf.Group = group
	return f
}

// Env sets the environment variable override key for this flag.
func (f *SliceFlag[T]) Env(key string) *SliceFlag[T] {
	if f.bf.DisableEnv {
		panic("cannot call Env after DisableEnv")
	}
	f.bf.EnvKey = key
	return f
}

// DisableEnv disables environment variable lookup for this flag.
func (f *SliceFlag[T]) DisableEnv() *SliceFlag[T] {
	if f.bf.EnvKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	f.bf.DisableEnv = true
	return f
}

// Metavar sets the name to show in usage text for this flag’s value.
func (f *SliceFlag[T]) Metavar(s string) *SliceFlag[T] {
	f.bf.Metavar = s
	return f
}
