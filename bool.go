package tinyflags

import "strconv"

// BoolValue holds the internal state of a boolean flag and whether it is strict.
type BoolValue struct {
	*ValueImpl[bool]
	strict bool
}

// IsStrictBool reports whether the flag requires an explicit value (--flag=true/false).
func (b *BoolValue) IsStrictBool() bool {
	return b.strict
}

// BoolFlag provides fluent builder methods for boolean flags,
// including support for .Strict() to require explicit values.
type BoolFlag struct {
	*Flag[bool] // embeds core builder methods like Env(), Required(), etc.
	val         *BoolValue
}

// Strict marks this boolean flag as requiring an explicit value.
func (b *BoolFlag) Strict() *BoolFlag {
	b.val.strict = true
	return b
}

// Bool defines a boolean flag with a default value.
func (f *FlagSet) Bool(name string, def bool, usage string) *BoolFlag {
	return f.BoolVarP(new(bool), name, "", def, usage)
}

// BoolP defines a boolean flag with a short name and a default value.
func (f *FlagSet) BoolP(name, short string, def bool, usage string) *BoolFlag {
	return f.BoolVarP(new(bool), name, short, def, usage)
}

// BoolVar defines a boolean flag and binds it to the given pointer.
func (f *FlagSet) BoolVar(ptr *bool, name string, def bool, usage string) *BoolFlag {
	return f.BoolVarP(ptr, name, "", def, usage)
}

// BoolVarP defines a boolean flag with a short name and binds it to the given pointer.
func (f *FlagSet) BoolVarP(ptr *bool, name, short string, def bool, usage string) *BoolFlag {
	val := &BoolValue{
		ValueImpl: NewValueImpl(
			ptr,
			def,
			strconv.ParseBool,
			strconv.FormatBool,
		),
	}
	flag := addScalar(f, name, short, usage, val, ptr)
	return &BoolFlag{Flag: flag, val: val}
}

// Bool defines a dynamic boolean flag under the group (e.g. --http.alpha.debug=true).
func (g *DynamicGroup) Bool(field, usage string) *DynamicBoolFlag {
	item := NewDynamicItemImpl(
		field,
		strconv.ParseBool,
		strconv.FormatBool,
	)

	g.items[field] = item

	// Register with baseFlag and dynamic map
	return addDynamicBool(g.fs, g.prefix, field, usage, item)
}
