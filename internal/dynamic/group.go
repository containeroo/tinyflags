package dynamic

import (
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
