package dynamic

// defineDynamicScalar creates and registers a dynamic scalar flag of type T.
func defineDynamicScalar[T any](
	g *Group,
	field string,
	parse func(string) (T, error),
	format func(T) string,
) *ScalarFlag[T] {
	item := NewDynamicScalarValue(field, parse, format)
	if err := g.fs.RegisterDynamic(g.prefix, field, item); err != nil {
		panic(err)
	}
	return &ScalarFlag[T]{item: item}
}

// defineDynamicSlice creates and registers a dynamic slice flag of type T.
func defineDynamicSlice[T any](
	g *Group,
	field string,
	parse func(string) (T, error),
	format func(T) string,
) *SliceFlag[T] {
	item := NewDynamicSliceValue(field, parse, format, g.fs.DefaultDelimiter())
	if err := g.fs.RegisterDynamic(g.prefix, field, item); err != nil {
		panic(err)
	}
	return &SliceFlag[T]{item: item}
}
