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

func Get[T any](g *Group, field, id string) (T, error) {
	var zero T

	item, ok := g.Items()[field]
	if !ok {
		return zero, fmt.Errorf("field not registered: %q", field)
	}
	v, ok := item.(*DynamicScalarValue[T])
	if !ok {
		return zero, fmt.Errorf("field %q has wrong type", field)
	}
	val, ok := v.values[id]
	if !ok {
		return zero, fmt.Errorf("value for field %q not found for id %q", field, id)
	}
	return val, nil
}

func MustGet[T any](g *Group, field, id string) T {
	val, err := Get[T](g, field, id)
	if err != nil {
		panic(err)
	}
	return val
}
