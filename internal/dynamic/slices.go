package dynamic

import "strconv"

// StringSlice registers a dynamic string slice flag under this group.
func (g *Group) StringSlice(field, usage string) *SliceFlag[string] {
	return defineDynamicSlice(g,
		field,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)
}

// IntSlice registers a dynamic int slice flag under this group.
func (g *Group) IntSlice(field, usage string) *SliceFlag[int] {
	return defineDynamicSlice(g,
		field,
		strconv.Atoi,
		strconv.Itoa,
	)
}
