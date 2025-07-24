package dynamic

import (
	"fmt"
	"slices"

	"github.com/containeroo/tinyflags/internal/core"
)

// Group manages a set of dynamic flags under one prefix.
type Group struct {
	fs     FlagSetRef                   // parent flagset
	prefix string                       // e.g. "http"
	items  map[string]core.DynamicValue // field → parser
}

// NewGroup starts a new dynamic group.
func NewGroup(fs FlagSetRef, prefix string) *Group {
	return &Group{fs: fs, prefix: prefix, items: map[string]core.DynamicValue{}}
}

func (g *Group) Name() string {
	return g.prefix
}

func (g *Group) Items() map[string]core.DynamicValue {
	return g.items
}

// Instances returns all seen IDs, sorted.
func (g *Group) Instances() []string {
	seen := map[string]struct{}{}
	for _, v := range g.items {
		if di, ok := v.(core.DynamicItemValues); ok {
			for id := range di.ValuesAny() {
				seen[id] = struct{}{}
			}
		}
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	return ids
}

func (g *Group) Lookup(field string) (core.DynamicValue, bool) {
	val, ok := g.items[field]
	return val, ok
}

// Get returns all fields for a given ID as a map of field → value.
func (g *Group) Get(id string) map[string]any {
	out := make(map[string]any)
	for field, val := range g.items {
		if v, ok := val.GetAny(id); ok {
			out[field] = v
		}
	}
	return out
}

func Get[T any](g *Group, id, flag string) (T, error) {
	var zero T

	item, ok := g.Items()[flag]
	if !ok {
		return zero, fmt.Errorf("field not registered: %q", flag)
	}
	v, ok := item.(*DynamicScalarValue[T])
	if !ok {
		return zero, fmt.Errorf("field %q has wrong type", flag)
	}
	val, ok := v.values[id]
	if !ok {
		return zero, fmt.Errorf("value for field %q not found for id %q", flag, id)
	}
	return val, nil
}

func MustGet[T any](g *Group, id, flag string) T {
	val, err := Get[T](g, flag, id)
	if err != nil {
		panic(err)
	}
	return val
}

func GetOrDefault[T any](g *Group, id, flag string) T {
	item, ok := g.Items()[flag]
	if !ok {
		panic(fmt.Sprintf("field not registered: %q", flag))
	}
	v, ok := item.(*DynamicScalarValue[T])
	if !ok {
		panic(fmt.Sprintf("field %q has wrong type", flag))
	}
	val, ok := v.values[id]
	if !ok {
		return v.def
	}
	return val
}
