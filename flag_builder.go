package tinyflags

type builderBase[T any] struct {
	fs    *FlagSet
	bf    *baseFlag
	value *FlagValue[T]
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
