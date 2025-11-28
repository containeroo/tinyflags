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
func (s *StaticFlag[T]) Required() *StaticFlag[T] {
	s.bf.Required = true
	return s
}

// HideRequired hides the “(Required)” suffix from help.
func (s *StaticFlag[T]) HideRequired() *StaticFlag[T] {
	s.bf.HideRequired = true
	return s
}

// Hidden omits the flag from help output.
func (s *StaticFlag[T]) Hidden() *StaticFlag[T] {
	s.bf.Hidden = true
	return s
}

// Deprecated adds a deprecation notice.
func (s *StaticFlag[T]) Deprecated(reason string) *StaticFlag[T] {
	s.bf.Deprecated = reason
	return s
}

// OneOfGroup assigns this flag to an exclusive group.
func (s *StaticFlag[T]) OneOfGroup(name string) *StaticFlag[T] {
	if name == "" {
		return s
	}
	g := s.registry.GetOneOfGroup(name)
	g.Flags = append(g.Flags, s.bf)
	s.bf.OneOfGroup = g
	return s
}

// AllOrNone assigns this flag to a require-together group.
func (s *StaticFlag[T]) AllOrNone(name string) *StaticFlag[T] {
	if name == "" {
		return s
	}
	g := s.registry.GetAllOrNoneGroup(name)
	g.Flags = append(g.Flags, s.bf)
	s.bf.AllOrNone = g
	return s
}

// Env sets a custom environment‐variable key.
func (s *StaticFlag[T]) Env(key string) *StaticFlag[T] {
	if s.bf.DisableEnv {
		panic("cannot call Env after DisableEnv")
	}
	s.bf.EnvKey = key
	return s
}

// DisableEnv turns off environment‐variable lookup.
func (s *StaticFlag[T]) DisableEnv() *StaticFlag[T] {
	if s.bf.EnvKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	s.bf.DisableEnv = true
	return s
}

// Placeholder customizes the value placeholder in usage.
func (b *StaticFlag[T]) Placeholder(s string) *StaticFlag[T] {
	b.bf.Placeholder = s
	return b
}

// Value exposes the underlying pointer for reading.
func (s *StaticFlag[T]) Value() *T {
	return s.ptr
}

// Allowed restricts help to show only these formatted values.
func (s *StaticFlag[T]) Allowed(vals ...string) *StaticFlag[T] {
	// copy to prevent external mutation
	s.bf.Allowed = append([]string(nil), vals...)
	return s
}

// HideAllowed hides the allowed values from help.
func (s *StaticFlag[T]) HideAllowed() *StaticFlag[T] {
	s.bf.HideAllowed = true
	return s
}

// Requires marks this flag as required by the given flag.
func (s *StaticFlag[T]) Requires(names ...string) *StaticFlag[T] {
	// copy to prevent external mutation
	s.bf.Requires = append([]string(nil), names...)
	return s
}

// HideRequires hides the “(Requires)” suffix from help.
func (s *StaticFlag[T]) HideRequires() *StaticFlag[T] {
	s.bf.HideRequires = true
	return s
}
