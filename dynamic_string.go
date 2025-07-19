package tinyflags

func (g *DynamicGroup) String(field, usage string) *DynamicFlag[string] {
	item := NewDynamicItemImpl(
		field,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)

	g.items[field] = item

	return addDynamic(g.fs, g.prefix, field, usage, item)
}
