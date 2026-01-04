package builder

import "github.com/containeroo/tinyflags/internal/core"

// DynamicFlag provides common builder methods for dynamic flags.
type DynamicFlag[T any] struct {
	meta flagMeta // shared metadata helpers
}

// NewDynamicFlag returns a DynamicFlag ready for embedding.
func NewDynamicFlag[T any](
	reg core.Registry,
	bf *core.BaseFlag,
) *DynamicFlag[T] {
	return &DynamicFlag[T]{meta: flagMeta{registry: reg, bf: bf}}
}

// Required marks the flag as mandatory.
func (d *DynamicFlag[T]) Required() *DynamicFlag[T] {
	d.meta.required()
	return d
}

// HideRequired hides the “(Required)” suffix from help.
func (d *DynamicFlag[T]) HideRequired() *DynamicFlag[T] {
	d.meta.hideRequired()
	return d
}

// Hidden omits the flag from help output.
func (d *DynamicFlag[T]) Hidden() *DynamicFlag[T] {
	d.meta.hidden()
	return d
}

// Deprecated adds a deprecation notice.
func (d *DynamicFlag[T]) Deprecated(reason string) *DynamicFlag[T] {
	d.meta.deprecated(reason)
	return d
}

// OneOfGroup assigns this flag to an exclusive group.
func (d *DynamicFlag[T]) OneOfGroup(name string) *DynamicFlag[T] {
	d.meta.oneOfGroup(name)
	return d
}

// AllOrNone assigns this flag to a require-together group.
func (d *DynamicFlag[T]) AllOrNone(name string) *DynamicFlag[T] {
	d.meta.allOrNoneGroup(name)
	return d
}

// Env sets a custom environment‐variable key.
func (d *DynamicFlag[T]) Env(key string) *DynamicFlag[T] {
	d.meta.env(key)
	return d
}

// DisableEnv turns off environment‐variable lookup.
func (d *DynamicFlag[T]) DisableEnv() *DynamicFlag[T] {
	d.meta.disableEnv()
	return d
}

// Placeholder customizes the value placeholder in usage.
func (d *DynamicFlag[T]) Placeholder(s string) *DynamicFlag[T] {
	d.meta.placeholder(s)
	return d
}

// Section assigns this flag to a help section.
func (d *DynamicFlag[T]) Section(name string) *DynamicFlag[T] {
	d.meta.section(name)
	return d
}

// Allowed restricts help to show only these formatted values.
func (d *DynamicFlag[T]) Allowed(vals ...string) *DynamicFlag[T] {
	// copy to prevent external mutation
	d.meta.allowed(vals...)
	return d
}

// HideAllowed hides the allowed values from help.
func (d *DynamicFlag[T]) HideAllowed() *DynamicFlag[T] {
	d.meta.hideAllowed()
	return d
}

// OverriddenValueMaskFn sets a mask function used by OverriddenValues().
func (d *DynamicFlag[T]) OverriddenValueMaskFn(fn func(any) any) *DynamicFlag[T] {
	d.meta.maskFn(fn)
	return d
}

// HideDefault hides the default value from help output.
func (d *DynamicFlag[T]) HideDefault() *DynamicFlag[T] {
	d.meta.hideDefault()
	return d
}

// Requires marks this flag as required by the given flag.
func (d *DynamicFlag[T]) Requires(names ...string) *DynamicFlag[T] {
	// copy to prevent external mutation
	d.meta.requires(names...)
	return d
}

// HideRequires hides the “(Requires)” suffix from help.
func (d *DynamicFlag[T]) HideRequires() *DynamicFlag[T] {
	d.meta.hideRequires()
	return d
}
