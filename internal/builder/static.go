package builder

import "github.com/containeroo/tinyflags/internal/core"

// StaticFlag provides common builder methods for scalar and slice flags.
type StaticFlag[T any] struct {
	meta flagMeta // shared metadata helpers
	ptr  *T       // destination for parsed values
}

// NewStaticFlag returns a DefaultFlag ready for embedding.
func NewStaticFlag[T any](reg core.Registry, bf *core.BaseFlag, ptr *T) StaticFlag[T] {
	return StaticFlag[T]{meta: flagMeta{registry: reg, bf: bf}, ptr: ptr}
}

// Short sets the one‐letter alias for this flag.
// Panics if you pass an empty or multi-rune string.
func (b *StaticFlag[T]) Short(s string) *StaticFlag[T] {
	// count code points, not bytes
	if len([]rune(s)) != 1 {
		panic("Short: alias must be exactly one character")
	}
	b.meta.bf.Short = s
	return b
}

// Required marks the flag as mandatory.
func (s *StaticFlag[T]) Required() *StaticFlag[T] {
	s.meta.required()
	return s
}

// HideRequired hides the “(Required)” suffix from help.
func (s *StaticFlag[T]) HideRequired() *StaticFlag[T] {
	s.meta.hideRequired()
	return s
}

// Hidden omits the flag from help output.
func (s *StaticFlag[T]) Hidden() *StaticFlag[T] {
	s.meta.hidden()
	return s
}

// Deprecated adds a deprecation notice.
func (s *StaticFlag[T]) Deprecated(reason string) *StaticFlag[T] {
	s.meta.deprecated(reason)
	return s
}

// OneOfGroup assigns this flag to an exclusive group.
func (s *StaticFlag[T]) OneOfGroup(name string) *StaticFlag[T] {
	s.meta.oneOfGroup(name)
	return s
}

// AllOrNone assigns this flag to a require-together group.
func (s *StaticFlag[T]) AllOrNone(name string) *StaticFlag[T] {
	s.meta.allOrNoneGroup(name)
	return s
}

// Env sets a custom environment‐variable key.
func (s *StaticFlag[T]) Env(key string) *StaticFlag[T] {
	s.meta.env(key)
	return s
}

// DisableEnv turns off environment‐variable lookup.
func (s *StaticFlag[T]) DisableEnv() *StaticFlag[T] {
	s.meta.disableEnv()
	return s
}

// Placeholder customizes the value placeholder in usage.
func (b *StaticFlag[T]) Placeholder(s string) *StaticFlag[T] {
	b.meta.placeholder(s)
	return b
}

// Section assigns this flag to a help section.
func (b *StaticFlag[T]) Section(name string) *StaticFlag[T] {
	b.meta.section(name)
	return b
}

// Value exposes the underlying pointer for reading.
func (s *StaticFlag[T]) Value() *T {
	return s.ptr
}

// Allowed restricts help to show only these formatted values.
func (s *StaticFlag[T]) Allowed(vals ...string) *StaticFlag[T] {
	// copy to prevent external mutation
	s.meta.allowed(vals...)
	return s
}

// HideAllowed hides the allowed values from help.
func (s *StaticFlag[T]) HideAllowed() *StaticFlag[T] {
	s.meta.hideAllowed()
	return s
}

// HideDefault hides the default value from help output.
func (s *StaticFlag[T]) HideDefault() *StaticFlag[T] {
	s.meta.hideDefault()
	return s
}

// Requires marks this flag as required by the given flag.
func (s *StaticFlag[T]) Requires(names ...string) *StaticFlag[T] {
	// copy to prevent external mutation
	s.meta.requires(names...)
	return s
}

// HideRequires hides the “(Requires)” suffix from help.
func (s *StaticFlag[T]) HideRequires() *StaticFlag[T] {
	s.meta.hideRequires()
	return s
}
