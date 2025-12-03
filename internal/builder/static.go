package builder

import "github.com/containeroo/tinyflags/internal/core"

// StaticFlag provides common builder methods for scalar and slice flags.
// The Self type parameter allows fluent methods to return the concrete flag type.
type StaticFlag[T any, Self any] struct {
	meta flagMeta // shared metadata helpers
	ptr  *T       // destination for parsed values
	self Self     // concrete flag for fluent returns
}

// NewStaticFlag returns a DefaultFlag ready for embedding.
// Pass the concrete flag as self to keep fluent chaining on that type.
func NewStaticFlag[T any, Self any](reg core.Registry, bf *core.BaseFlag, ptr *T, self Self) StaticFlag[T, Self] {
	return StaticFlag[T, Self]{meta: flagMeta{registry: reg, bf: bf}, ptr: ptr, self: self}
}

// Short sets the one‐letter alias for this flag.
// Panics if you pass an empty or multi-rune string.
func (b *StaticFlag[T, Self]) Short(s string) Self {
	// count code points, not bytes
	if len([]rune(s)) != 1 {
		panic("Short: alias must be exactly one character")
	}
	b.meta.bf.Short = s
	return b.self
}

// Required marks the flag as mandatory.
func (s *StaticFlag[T, Self]) Required() Self {
	s.meta.required()
	return s.self
}

// HideRequired hides the "(Required)" suffix from help.
func (s *StaticFlag[T, Self]) HideRequired() Self {
	s.meta.hideRequired()
	return s.self
}

// Hidden omits the flag from help output.
func (s *StaticFlag[T, Self]) Hidden() Self {
	s.meta.hidden()
	return s.self
}

// Deprecated adds a deprecation notice.
func (s *StaticFlag[T, Self]) Deprecated(reason string) Self {
	s.meta.deprecated(reason)
	return s.self
}

// OneOfGroup assigns this flag to an exclusive group.
func (s *StaticFlag[T, Self]) OneOfGroup(name string) Self {
	s.meta.oneOfGroup(name)
	return s.self
}

// AllOrNone assigns this flag to a require-together group.
func (s *StaticFlag[T, Self]) AllOrNone(name string) Self {
	s.meta.allOrNoneGroup(name)
	return s.self
}

// Env sets a custom environment‐variable key.
func (s *StaticFlag[T, Self]) Env(key string) Self {
	s.meta.env(key)
	return s.self
}

// DisableEnv turns off environment‐variable lookup.
func (s *StaticFlag[T, Self]) DisableEnv() Self {
	s.meta.disableEnv()
	return s.self
}

// Placeholder customizes the value placeholder in usage.
func (b *StaticFlag[T, Self]) Placeholder(s string) Self {
	b.meta.placeholder(s)
	return b.self
}

// Section assigns this flag to a help section.
func (b *StaticFlag[T, Self]) Section(name string) Self {
	b.meta.section(name)
	return b.self
}

// Value exposes the underlying pointer for reading.
func (s *StaticFlag[T, Self]) Value() *T {
	return s.ptr
}

// Allowed restricts help to show only these formatted values.
func (s *StaticFlag[T, Self]) Allowed(vals ...string) Self {
	// copy to prevent external mutation
	s.meta.allowed(vals...)
	return s.self
}

// HideAllowed hides the allowed values from help.
func (s *StaticFlag[T, Self]) HideAllowed() Self {
	s.meta.hideAllowed()
	return s.self
}

// HideDefault hides the default value from help output.
func (s *StaticFlag[T, Self]) HideDefault() Self {
	s.meta.hideDefault()
	return s.self
}

// Requires marks this flag as required by the given flag.
func (s *StaticFlag[T, Self]) Requires(names ...string) Self {
	// copy to prevent external mutation
	s.meta.requires(names...)
	return s.self
}

// HideRequires hides the "(Requires)" suffix from help.
func (s *StaticFlag[T, Self]) HideRequires() Self {
	s.meta.hideRequires()
	return s.self
}
