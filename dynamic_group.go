package tinyflags

import "slices"

// DynamicGroup defines a prefix for dynamic flags (e.g., "http").
// Flags under this group are declared using `group.String(...)`, etc.
type DynamicGroup struct {
	fs     *FlagSet                // parent flag set
	prefix string                  // e.g. "http"
	items  map[string]DynamicValue // field name → dynamic value item
}

// Instances returns all seen instance IDs (e.g., alpha, beta).
func (g *DynamicGroup) Instances() []string {
	seen := make(map[string]struct{})
	for _, v := range g.items {
		if impl, ok := v.(dynamicItemValues); ok {
			for id := range impl.ValuesAny() {
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
