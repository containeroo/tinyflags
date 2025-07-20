package dynamic

import (
	"slices"

	"github.com/containeroo/tinyflags/internal/core"
)

// Group represents a dynamic flag group such as `--http.alpha.port`.
// Each instance (e.g., "alpha") may set fields defined under this group.
type Group struct {
	fs     FlagSetRef                   // reference to parent FlagSet
	prefix string                       // shared prefix, e.g., "http"
	items  map[string]core.DynamicValue // field name → dynamic value storage
}

// NewGroup creates a new dynamic group with the given prefix.
func NewGroup(fs FlagSetRef, prefix string) *Group {
	return &Group{
		fs:     fs,
		prefix: prefix,
		items:  make(map[string]core.DynamicValue),
	}
}

// Instances returns a sorted list of all instance IDs seen across fields.
func (g *Group) Instances() []string {
	seen := map[string]struct{}{}
	for _, v := range g.items {
		if d, ok := v.(core.DynamicItemValues); ok {
			for id := range d.ValuesAny() {
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
