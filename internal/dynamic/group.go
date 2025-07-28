package dynamic

import (
	"fmt"
	"slices"
	"sort"

	"github.com/containeroo/tinyflags/internal/core"
)

// Group manages a set of dynamic flags under one prefix.
type Group struct {
	fs          FlagSetRef                // Parent flagset reference
	name        string                    // Group name prefix (e.g. "http")
	items       map[string]core.GroupItem // Registered flags and their dynamic values
	itemOrder   []*core.BaseFlag          // Preserved registration order
	sortFlags   bool                      // Whether internal flag sorting is enabled
	hidden      bool                      // Whether to hide the group from help
	title       string                    // Group section title
	description string                    // Group description
	placeholder string                    // Optional usage placeholder for ID
	notes       string                    // Optional help note
}

// NewGroup starts a new dynamic group.
func NewGroup(fs FlagSetRef, prefix string) *Group {
	return &Group{
		fs:    fs,
		name:  prefix,
		items: map[string]core.GroupItem{},
	}
}

// SortFlags enables sorting of flags within the group.
func (g *Group) SortFlags() *Group {
	g.sortFlags = true
	return g
}

// Hidden marks the group as hidden from usage output.
func (g *Group) Hidden() *Group {
	g.hidden = true
	return g
}

// Title sets the group section title.
func (g *Group) Title(s string) *Group {
	g.title = s
	return g
}

// Description sets the group description.
func (g *Group) Description(s string) *Group {
	g.description = s
	return g
}

// Note sets an optional help note for the group.
func (g *Group) Note(s string) *Group {
	g.notes = s
	return g
}

// Placeholder sets the placeholder string for ID in usage.
func (g *Group) Placeholder(s string) *Group {
	g.placeholder = s
	return g
}

// GetPlaceholder returns the current placeholder string.
func (g *Group) GetPlaceholder() string {
	return g.placeholder
}

// Name returns the group prefix.
func (g *Group) Name() string {
	return g.name
}

// Items returns all registered group items.
func (g *Group) Items() map[string]core.GroupItem {
	return g.items
}

// Flags returns flags in registration order.
func (g *Group) Flags() []*core.BaseFlag {
	return g.itemOrder
}

// DynamicFlags returns flags
func (g *Group) DynamicFlags() []*core.BaseFlag {
	if g.sortFlags {
		out := make([]*core.BaseFlag, len(g.itemOrder))
		copy(out, g.itemOrder)
		sort.Slice(out, func(i, j int) bool {
			return out[i].Name < out[j].Name
		})
		return out
	}
	return g.itemOrder
}

// IsFlagSorted reports whether internal flags should be sorted.
func (g *Group) IsFlagSorted() bool {
	return g.sortFlags
}

// IsHidden reports whether the group is hidden.
func (g *Group) IsHidden() bool {
	return g.hidden
}

// TitleText returns the group title.
func (g *Group) TitleText() string {
	return g.title
}

// DescriptionText returns the group description.
func (g *Group) DescriptionText() string {
	return g.description
}

// NoteText returns the group notes.
func (g *Group) NoteText() string {
	return g.notes
}

// Instances returns a sorted list of all seen instance IDs.
func (g *Group) Instances() []string {
	seen := map[string]struct{}{}
	for _, item := range g.items {
		if di, ok := item.Value.(core.DynamicItemValues); ok {
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

// Get returns all flag values for a given instance ID.
func (g *Group) Get(id string) map[string]any {
	out := make(map[string]any)
	for field, item := range g.items {
		if v, ok := item.Value.GetAny(id); ok {
			out[field] = v
		}
	}
	return out
}

// Lookup retrieves the dynamic value interface for a given field.
func (g *Group) Lookup(field string) (core.DynamicValue, bool) {
	item, ok := g.items[field]
	if !ok {
		return nil, false
	}
	return item.Value, true
}

// LookupFlag retrieves the base flag metadata for a given field.
func (g *Group) LookupFlag(field string) *core.BaseFlag {
	item, ok := g.items[field]
	if !ok {
		return nil
	}
	return item.Flag
}

// Get retrieves a typed value for the given ID and flag field.
func Get[T any](g *Group, id, flag string) (T, error) {
	var zero T
	item, ok := g.Items()[flag]
	if !ok {
		return zero, fmt.Errorf("dynamic flag %q is not registered", flag)
	}
	v, ok := item.Value.GetAny(id)
	if !ok {
		return zero, fmt.Errorf("no value set for --%s.%s.%s", g.Name(), id, flag)
	}
	typed, ok := v.(T)
	if !ok {
		return zero, fmt.Errorf("dynamic flag %q has unexpected type", flag)
	}
	return typed, nil
}

// MustGet is like Get but panics on failure.
func MustGet[T any](g *Group, id, flag string) T {
	val, err := Get[T](g, id, flag)
	if err != nil {
		panic(err)
	}
	return val
}

// GetOrDefault returns the typed value or default if not set.
func GetOrDefault[T any](g *Group, id, flag string) T {
	item, ok := g.Items()[flag]
	if !ok {
		panic(fmt.Sprintf("dynamic flag %q is not registered", flag))
	}
	v, _ := item.Value.GetAny(id)
	typed, ok := v.(T)
	if !ok {
		panic(fmt.Sprintf("dynamic flag %q has unexpected type", flag))
	}
	return typed
}
