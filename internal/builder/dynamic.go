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
func (b *DynamicFlag[T]) Required() *DynamicFlag[T] {
	b.bf.Required = true
	return b
}

// Hidden omits the flag from help output.
func (b *DynamicFlag[T]) Hidden() *DynamicFlag[T] {
	b.bf.Hidden = true
	return b
}

// Deprecated adds a deprecation notice.
func (b *DynamicFlag[T]) Deprecated(reason string) *DynamicFlag[T] {
	b.bf.Deprecated = reason
	return b
}

// MutualExlusive assigns this flag to an exclusive group.
func (b *DynamicFlag[T]) MutualExlusive(name string) *DynamicFlag[T] {
	if name == "" {
		return b
	}
	g := b.registry.GetMutualGroup(name)
	g.Flags = append(g.Flags, b.bf)
	b.bf.MutualGroup = g
	return b
}

// RequireTogether assigns this flag to a require-together group.
func (b *DynamicFlag[T]) RequireTogether(name string) *DynamicFlag[T] {
	if name == "" {
		return b
	}
	g := b.registry.GetRequireTogetherGroup(name)
	g.Flags = append(g.Flags, b.bf)
	return b
}

// Env sets a custom environment‐variable key.
func (b *DynamicFlag[T]) Env(key string) *DynamicFlag[T] {
	if b.bf.DisableEnv {
		panic("cannot call Env after DisableEnv")
	}
	b.bf.EnvKey = key
	return b
}

// DisableEnv turns off environment‐variable lookup.
func (b *DynamicFlag[T]) DisableEnv() *DynamicFlag[T] {
	if b.bf.EnvKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	b.bf.DisableEnv = true
	return b
}

// Placeholder customizes the value placeholder in usage.
func (b *DynamicFlag[T]) Placeholder(s string) *DynamicFlag[T] {
	b.bf.Placeholder = s
	return b
}

// Allowed restricts help to show only these formatted values.
func (b *DynamicFlag[T]) Allowed(vals ...string) *DynamicFlag[T] {
	// copy to prevent external mutation
	b.bf.Allowed = append([]string(nil), vals...)
	return b
}
