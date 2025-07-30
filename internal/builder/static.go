package builder

import "github.com/containeroo/tinyflags/internal/core"

// StaticFlag provides common builder methods for scalar and slice flags.
type StaticFlag[T any] struct {
	registry core.Registry  // registry of this flag
	bf       *core.BaseFlag // core metadata for usage/errors
	ptr      *T             // destination for parsed values
}

// NewStaticFlag returns a DefaultFlag ready for embedding.
func NewStaticFlag[T any](reg core.Registry, bf *core.BaseFlag, ptr *T) StaticFlag[T] {
	return StaticFlag[T]{registry: reg, bf: bf, ptr: ptr}
}

// Short sets the one‐letter alias for this flag.
// Panics if you pass an empty or multi-rune string.
func (b *StaticFlag[T]) Short(s string) *StaticFlag[T] {
	// count code points, not bytes
	if len([]rune(s)) != 1 {
		panic("Short: alias must be exactly one character")
	}
	b.bf.Short = s
	return b
}

// Required marks the flag as mandatory.
func (b *StaticFlag[T]) Required() *StaticFlag[T] {
	b.bf.Required = true
	return b
}

// Hidden omits the flag from help output.
func (b *StaticFlag[T]) Hidden() *StaticFlag[T] {
	b.bf.Hidden = true
	return b
}

// Deprecated adds a deprecation notice.
func (b *StaticFlag[T]) Deprecated(reason string) *StaticFlag[T] {
	b.bf.Deprecated = reason
	return b
}

// MutualExlusive assigns this flag to an exclusive group.
func (b *StaticFlag[T]) MutualExlusive(name string) *StaticFlag[T] {
	if name == "" {
		return b
	}
	g := b.registry.GetMutualGroup(name)
	g.Flags = append(g.Flags, b.bf)
	b.bf.MutualGroup = g
	return b
}

// RequireTogether assigns this flag to a require-together group.
func (b *StaticFlag[T]) RequireTogether(name string) *StaticFlag[T] {
	if name == "" {
		return b
	}
	g := b.registry.GetRequireTogetherGroup(name)
	g.Flags = append(g.Flags, b.bf)
	return b
}

// Env sets a custom environment‐variable key.
func (b *StaticFlag[T]) Env(key string) *StaticFlag[T] {
	if b.bf.DisableEnv {
		panic("cannot call Env after DisableEnv")
	}
	b.bf.EnvKey = key
	return b
}

// DisableEnv turns off environment‐variable lookup.
func (b *StaticFlag[T]) DisableEnv() *StaticFlag[T] {
	if b.bf.EnvKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	b.bf.DisableEnv = true
	return b
}

// Placeholder customizes the value placeholder in usage.
func (b *StaticFlag[T]) Placeholder(s string) *StaticFlag[T] {
	b.bf.Placeholder = s
	return b
}

// Value exposes the underlying pointer for reading.
func (b *StaticFlag[T]) Value() *T {
	return b.ptr
}

// Allowed restricts help to show only these formatted values.
func (b *StaticFlag[T]) Allowed(vals ...string) *StaticFlag[T] {
	// copy to prevent external mutation
	b.bf.Allowed = append([]string(nil), vals...)
	return b
}
