package dynamic

import "github.com/containeroo/tinyflags/internal/builder"

type BoolFlag struct {
	*builder.DynamicFlag[bool]
	item       *DynamicScalarValue[bool]
	strictMode *bool
}

func (b *BoolFlag) Strict() *BoolFlag {
	*b.strictMode = true
	return b
}

func (f *BoolFlag) Get(id string) (bool, bool) {
	val, ok := f.item.values[id]
	if !ok {
		return f.item.def, false
	}
	return val, true
}

func (f *BoolFlag) MustGet(id string) bool {
	val, ok := f.Get(id)
	if !ok {
		panic("missing value for bool flag: " + f.item.field + " (" + id + ")")
	}
	return val
}

func (f *BoolFlag) Has(id string) bool {
	_, ok := f.item.values[id]
	return ok
}

func (f *BoolFlag) Values() map[string]bool {
	return f.item.values
}

func (f *BoolFlag) ValuesAny() map[string]any {
	out := make(map[string]any, len(f.item.values))
	for k, v := range f.item.values {
		out[k] = v
	}
	return out
}
