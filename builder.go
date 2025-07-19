package tinyflags

// builderImpl provides common builder methods for scalar and slice flags.
type builderImpl[T any] struct {
	fs    *FlagSet
	bf    *baseFlag
	value *ScalarValueImpl[T]
	ptr   *T
}

func (b *builderImpl[T]) Required() *builderImpl[T] {
	b.bf.required = true
	return b
}

func (b *builderImpl[T]) Hidden() *builderImpl[T] {
	b.bf.hidden = true
	return b
}

func (b *builderImpl[T]) Deprecated(reason string) *builderImpl[T] {
	b.bf.deprecated = reason
	return b
}

func (b *builderImpl[T]) Group(name string) *builderImpl[T] {
	if name == "" {
		return b
	}
	for _, g := range b.fs.groups {
		if g.name == name {
			g.flags = append(g.flags, b.bf)
			b.bf.group = g
			return b
		}
	}
	group := &mutualGroup{name: name, flags: []*baseFlag{b.bf}}
	b.fs.groups = append(b.fs.groups, group)
	b.bf.group = group
	return b
}

func (b *builderImpl[T]) Env(key string) *builderImpl[T] {
	if b.bf.disableEnv {
		panic("cannot call Env after DisableEnv")
	}
	b.bf.envKey = key
	return b
}

func (b *builderImpl[T]) DisableEnv() *builderImpl[T] {
	if b.bf.envKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	b.bf.disableEnv = true
	return b
}

func (b *builderImpl[T]) Metavar(s string) *builderImpl[T] {
	b.bf.metavar = s
	return b
}

func (b *builderImpl[T]) Validator(fn func(T) error) *builderImpl[T] {
	b.value.SetValidator(fn)
	return b
}

func (b *builderImpl[T]) Value() *T {
	return b.ptr
}
