package dynamic

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/core"
)

// ScalarFlag provides access to a dynamic scalar flag.
// For example: --http.alpha.port=8080
type ScalarFlag[T any] struct {
	registry core.Registry
	bf       *core.BaseFlag
	item     *DynamicScalarValue[T] // value storage and parsing logic
}

// Get returns the parsed value for a given instance ID.
func (f *ScalarFlag[T]) Get(id string) (T, bool) {
	return f.item.Get(id)
}

// MustGet returns the value or panics if not set.
func (f *ScalarFlag[T]) MustGet(id string) T {
	val, ok := f.item.Get(id)
	if !ok {
		panic(fmt.Sprintf("value for id %q not set", id))
	}
	return val
}

// Values returns all instance values.
func (f *ScalarFlag[T]) Values() map[string]T {
	return f.item.Values()
}

// ValuesAny returns instance values as a generic map.
func (f *ScalarFlag[T]) ValuesAny() map[string]any {
	return f.item.ValuesAny()
}

func (f *ScalarFlag[T]) Validate(fn func(T) error) *ScalarFlag[T] {
	f.item.setValidate(fn)
	return f
}

// Required marks this flag as required.
func (f *ScalarFlag[T]) Required() *ScalarFlag[T] {
	f.bf.Required = true
	return f
}

// Hidden marks this flag as hidden from help output.
func (f *ScalarFlag[T]) Hidden() *ScalarFlag[T] {
	f.bf.Hidden = true
	return f
}

// Deprecated marks this flag as deprecated with an optional reason.
func (f *ScalarFlag[T]) Deprecated(reason string) *ScalarFlag[T] {
	f.bf.Deprecated = reason
	return f
}

// Group adds this flag to a mutual exclusion group.
func (f *ScalarFlag[T]) Group(name string) *ScalarFlag[T] {
	if name == "" {
		return f
	}
	group := f.registry.GetGroup(name)
	group.Flags = append(group.Flags, f.bf)
	f.bf.Group = group
	return f
}

// Env sets the environment variable override key for this flag.
func (f *ScalarFlag[T]) Env(key string) *ScalarFlag[T] {
	if f.bf.DisableEnv {
		panic("cannot call Env after DisableEnv")
	}
	f.bf.EnvKey = key
	return f
}

// DisableEnv disables environment variable lookup for this flag.
func (f *ScalarFlag[T]) DisableEnv() *ScalarFlag[T] {
	if f.bf.EnvKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	f.bf.DisableEnv = true
	return f
}

// Metavar sets the name to show in usage text for this flag’s value.
func (f *ScalarFlag[T]) Metavar(s string) *ScalarFlag[T] {
	f.bf.Metavar = s
	return f
}
