package tinyflags

// StringSlice defines a dynamic string slice flag under the group (e.g. --http.alpha.tags=one,two).
func (g *DynamicGroup) StringSlice(field, usage string) *DynamicSliceFlag[string] {
	item := NewDynamicSliceItemImpl(
		field,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
		g.fs.defaultDelimiter,
	)

	g.items[field] = item

	return addDynamicSlice(g.fs, g.prefix, field, usage, item)
}
