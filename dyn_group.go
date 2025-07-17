package tinyflags

import (
	"slices"
)

// DynamicGroup holds a set of dynamic fields under a shared prefix (e.g. "http").
// Fields are accessed using instance identifiers like --http.alpha.port.
type DynamicGroup struct {
	fs     *FlagSet                // parent flag set
	prefix string                  // used as prefix (e.g. "http")
	items  map[string]DynamicValue // field name → dynamic value (e.g. port, address)
}

// Instances returns a sorted list of unique instance IDs seen across all dynamic flags.
func (g *DynamicGroup) Instances() []string {
	seen := make(map[string]struct{})
	for _, v := range g.items {
		if dv, ok := v.(dynamicItemValues); ok {
			for id := range dv.ValuesAny() {
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
