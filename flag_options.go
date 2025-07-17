package tinyflags

import "fmt"

type builderBase[T any] struct {
	fs    *FlagSet
	bf    *baseFlag
	value *FlagBase[T]
	ptr   *T
}

func (b *builderBase[T]) Required() *builderBase[T] {
	b.bf.required = true
	return b
}

func (b *builderBase[T]) Hidden() *builderBase[T] {
	b.bf.hidden = true
	return b
}

func (b *builderBase[T]) Deprecated(reason string) *builderBase[T] {
	b.bf.deprecated = reason
	return b
}

func (b *builderBase[T]) Group(name string) *builderBase[T] {
	if name == "" {
		return b
	}
	var group *mutualGroup
	for _, g := range b.fs.groups {
		if g.name == name {
			group = g
			break
		}
	}
	if group == nil {
		group = &mutualGroup{name: name}
		b.fs.groups = append(b.fs.groups, group)
	}
	group.flags = append(group.flags, b.bf)
	b.bf.group = group
	return b
}

func (b *builderBase[T]) Env(key string) *builderBase[T] {
	if b.bf.disableEnv {
		panic("cannot call Env after DisableEnv")
	}
	b.bf.envKey = key
	return b
}

func (b *builderBase[T]) DisableEnv() *builderBase[T] {
	if b.bf.envKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	b.bf.disableEnv = true
	return b
}

func (b *builderBase[T]) Metavar(s string) *builderBase[T] {
	b.bf.metavar = s
	return b
}

func (b *builderBase[T]) Validator(fn func(T) error) *builderBase[T] {
	b.value.SetValidator(fn)
	return b
}

func (b *builderBase[T]) Value() *T {
	return b.ptr
}

// Choices restricts the allowed values for this flag to a predefined set.
func (b *builderBase[T]) Choices(allowed ...T) *builderBase[T] {
	bv, ok := b.bf.value.(*FlagBase[T])
	if !ok {
		return b
	}

	// Build validator from list
	bv.SetValidator(func(v T) error {
		for _, a := range allowed {
			if bv.format(a) == bv.format(v) {
				return nil
			}
		}
		return fmt.Errorf("must be one of %s", formatAllowed(allowed, bv.format))
	})

	// Convert allowed values to string for help text
	b.bf.allowed = make([]string, len(allowed))
	for i, x := range allowed {
		b.bf.allowed[i] = bv.format(x)
	}

	return b
}
