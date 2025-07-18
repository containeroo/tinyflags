package tinyflags

func (d *DynamicFlag[T]) Required() *DynamicFlag[T] {
	d.bf.required = true
	return d
}

func (d *DynamicFlag[T]) Hidden() *DynamicFlag[T] {
	d.bf.hidden = true
	return d
}

func (d *DynamicFlag[T]) Group(name string) *DynamicFlag[T] {
	if name == "" {
		return d
	}
	for _, g := range d.fs.groups {
		if g.name == name {
			g.flags = append(g.flags, d.bf)
			d.bf.group = g
			return d
		}
	}
	group := &mutualGroup{name: name, flags: []*baseFlag{d.bf}}
	d.fs.groups = append(d.fs.groups, group)
	d.bf.group = group
	return d
}

func (d *DynamicFlag[T]) Env(key string) *DynamicFlag[T] {
	if d.bf.disableEnv {
		panic("cannot call Env after DisableEnv")
	}
	d.bf.envKey = key
	return d
}

func (d *DynamicFlag[T]) DisableEnv() *DynamicFlag[T] {
	if d.bf.envKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	d.bf.disableEnv = true
	return d
}

func (d *DynamicFlag[T]) Metavar(s string) *DynamicFlag[T] {
	d.bf.metavar = s
	return d
}

func (d *DynamicFlag[T]) Validator(fn func(T) error) *DynamicFlag[T] {
	d.item.SetValidator(fn)
	return d
}
