package builder

import "github.com/containeroo/tinyflags/internal/core"

// DynamicFlag provides common builder methods for dynamic flags.
type DynamicFlag[T any] struct {
	registry core.Registry  // registry of this flag
	bf       *core.BaseFlag // core metadata for usage/errors
}

// NewDynamicFlag returns a DynamicFlag ready for embedding.
func NewDynamicFlag[T any](
	reg core.Registry,
	bf *core.BaseFlag,
) *DynamicFlag[T] {
	return &DynamicFlag[T]{registry: reg, bf: bf}
}

// Required marks the flag as mandatory.
func (d *DynamicFlag[T]) Required() *DynamicFlag[T] {
	d.bf.Required = true
	return d
}

// Hidden omits the flag from help output.
func (d *DynamicFlag[T]) Hidden() *DynamicFlag[T] {
	d.bf.Hidden = true
	return d
}

// Deprecated adds a deprecation notice.
func (d *DynamicFlag[T]) Deprecated(reason string) *DynamicFlag[T] {
	d.bf.Deprecated = reason
	return d
}

// OneOfGroup assigns this flag to an exclusive group.
func (d *DynamicFlag[T]) OneOfGroup(name string) *DynamicFlag[T] {
	if name == "" {
		return d
	}
	g := d.registry.GetOneOfGroup(name)
	g.Flags = append(g.Flags, d.bf)
	d.bf.OneOfGroup = g
	return d
}

// AllOrNone assigns this flag to a require-together group.
func (d *DynamicFlag[T]) AllOrNone(name string) *DynamicFlag[T] {
	if name == "" {
		return d
	}
	g := d.registry.GetAllOrNoneGroup(name)
	g.Flags = append(g.Flags, d.bf)
	d.bf.AllOrNone = g
	return d
}

// Env sets a custom environment‐variable key.
func (d *DynamicFlag[T]) Env(key string) *DynamicFlag[T] {
	if d.bf.DisableEnv {
		panic("cannot call Env after DisableEnv")
	}
	d.bf.EnvKey = key
	return d
}

// DisableEnv turns off environment‐variable lookup.
func (d *DynamicFlag[T]) DisableEnv() *DynamicFlag[T] {
	if d.bf.EnvKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	d.bf.DisableEnv = true
	return d
}

// Placeholder customizes the value placeholder in usage.
func (d *DynamicFlag[T]) Placeholder(s string) *DynamicFlag[T] {
	d.bf.Placeholder = s
	return d
}

// Allowed restricts help to show only these formatted values.
func (d *DynamicFlag[T]) Allowed(vals ...string) *DynamicFlag[T] {
	// copy to prevent external mutation
	d.bf.Allowed = append([]string(nil), vals...)
	return d
}
