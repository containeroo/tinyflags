package dynamic

import (
	"strconv"

	"github.com/containeroo/tinyflags/internal/builder"
	"github.com/containeroo/tinyflags/internal/core"
)

// BoolValue wraps a DynamicScalarValue[bool] and lets the FSM see IsStrictBool.
type BoolValue struct {
	item       *DynamicScalarValue[bool]
	strictMode *bool
}

func (s *BoolValue) Set(id, raw string) error {
	return s.item.Set(id, raw)
}

func (s *BoolValue) IsStrictBool() bool {
	return *s.strictMode
}

// BoolFlag is the builder for per-ID boolean flags.
type BoolFlag struct {
	*builder.DynamicFlag[bool]
	item       *DynamicScalarValue[bool]
	strictMode bool
}

// Strict requires explicit `=true|false`.
func (b *BoolFlag) Strict() *BoolFlag {
	b.strictMode = true
	return b
}

// Bool registers a dynamic boolean field under this group.
func (g *Group) Bool(field, usage string) *BoolFlag {
	// create the raw parser/storage
	item := NewDynamicScalarValue(field, strconv.ParseBool, strconv.FormatBool)

	// wrap it so it also implements StrictBool
	flagVal := &BoolValue{item: item, strictMode: new(bool)}

	// register for dynamic‐value parsing
	if err := g.fs.RegisterDynamic(g.prefix, field, flagVal); err != nil {
		panic(err)
	}

	// register a BaseFlag so it shows up in help
	bf := &core.BaseFlag{Name: field, Usage: usage}
	g.fs.RegisterFlag(field, bf)

	// build the fluent API
	df := builder.NewDynamicFlag[bool](g.fs, bf)

	// return the builder, wiring strictMode pointer
	return &BoolFlag{
		DynamicFlag: df,
		item:        item,
		strictMode:  *flagVal.strictMode,
	}
}

func (f *BoolFlag) Get(id string) (bool, bool) {
	val, ok := f.item.values[id]
	return val, ok
}

func (f *BoolFlag) MustGet(id string) bool {
	val, ok := f.Get(id)
	if !ok {
		panic("value not set")
	}
	return val
}

// Values returns all stored values.
func (f *BoolFlag) Values() map[string]bool {
	return f.item.values
}

// ValuesAny returns values as a generic map.
func (f *BoolFlag) ValuesAny() map[string]any {
	m := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		m[k] = v
	}
	return m
}
