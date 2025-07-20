package builder

import (
	"github.com/containeroo/tinyflags/internal/core"
)

// DefaultFlag provides common builder methods for scalar and slice flags.
type DefaultFlag[T any] struct {
	Registry core.Registry  // Reference to the flagset registry
	BF       *core.BaseFlag // Base coredata for the flag
	Ptr      *T             // Pointer to the destination variable
}

// Required marks this flag as required.
func (b *DefaultFlag[T]) Required() *DefaultFlag[T] {
	b.BF.Required = true
	return b
}

// Hidden marks this flag as hidden from help output.
func (b *DefaultFlag[T]) Hidden() *DefaultFlag[T] {
	b.BF.Hidden = true
	return b
}

// Deprecated marks this flag as deprecated with an optional reason.
func (b *DefaultFlag[T]) Deprecated(reason string) *DefaultFlag[T] {
	b.BF.Deprecated = reason
	return b
}

// Group adds this flag to a mutual exclusion group.
func (b *DefaultFlag[T]) Group(name string) *DefaultFlag[T] {
	if name == "" {
		return b
	}
	group := b.Registry.GetGroup(name)
	group.Flags = append(group.Flags, b.BF)
	b.BF.Group = group
	return b
}

// Env sets the environment variable override key for this flag.
func (b *DefaultFlag[T]) Env(key string) *DefaultFlag[T] {
	if b.BF.DisableEnv {
		panic("cannot call Env after DisableEnv")
	}
	b.BF.EnvKey = key
	return b
}

// DisableEnv disables environment variable lookup for this flag.
func (b *DefaultFlag[T]) DisableEnv() *DefaultFlag[T] {
	if b.BF.EnvKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	b.BF.DisableEnv = true
	return b
}

// Metavar sets the name to show in usage text for this flag’s value.
func (b *DefaultFlag[T]) Metavar(s string) *DefaultFlag[T] {
	b.BF.Metavar = s
	return b
}

// Value returns the destination pointer for this flag.
func (b *DefaultFlag[T]) Value() *T {
	return b.Ptr
}
