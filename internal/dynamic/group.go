package dynamic

import (
	"fmt"
	"slices"

	"github.com/containeroo/tinyflags/internal/core"
)

// Group manages a set of dynamic flags under one prefix.
type Group struct {
	fs          FlagSetRef                   // parent flagset
	prefix      string                       // e.g. "http"
	items       map[string]core.DynamicValue // field → parser
	sortGroup   bool                         // sort group items
	sortFlags   bool                         // sort flags
	hidden      bool                         // hide group from
	title       string                       // group title
	description string                       // group description
	notes       string                       // group notes
}

// NewGroup starts a new dynamic group.
func NewGroup(fs FlagSetRef, prefix string) *Group {
	return &Group{fs: fs, prefix: prefix, items: map[string]core.DynamicValue{}}
}

func (g *Group) SortGroup() *Group {
	if g.sortFlags {
		panic("cannot call SortGroup after SortFlags")
	}
	g.sortGroup = true
	return g
}

func (g *Group) SortFlags() *Group {
	if g.sortGroup {
		panic("cannot call SortFlags after SortGroup")
	}
	g.sortFlags = true
	return g
}

func (g *Group) Hidden() *Group {
	g.hidden = true
	return g
}

func (g *Group) Title(s string) *Group {
	g.title = s
	return g
}

func (g *Group) Description(s string) *Group {
	g.description = s
	return g
}

func (g *Group) Note(s string) *Group {
	g.notes = s
	return g
}

func (g *Group) Name() string                        { return g.prefix }
func (g *Group) Items() map[string]core.DynamicValue { return g.items }
func (g *Group) IsGroupSorted() bool                 { return g.sortGroup }
func (g *Group) IsFlagSorted() bool                  { return g.sortFlags }
func (g *Group) IsHidden() bool                      { return g.hidden }
func (g *Group) TitleText() string                   { return g.title }
func (g *Group) DescriptionText() string             { return g.description }
func (g *Group) NoteText() string                    { return g.notes }

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

func (g *Group) Get(id string) map[string]any {
	out := make(map[string]any)
	for field, val := range g.items {
		if v, ok := val.GetAny(id); ok {
			out[field] = v
		}
	}
	return out
}

func (g *Group) Lookup(field string) (core.DynamicValue, bool) {
	val, ok := g.items[field]
	return val, ok
}

func (g *Group) LookupFlag(field string) *core.BaseFlag {
	return g.fs.LookupFlag(field)
}

func (g *Group) Flags() []*core.BaseFlag {
	var out []*core.BaseFlag
	for name := range g.items {
		if fl := g.fs.LookupFlag(name); fl != nil {
			out = append(out, fl)
		}
	}
	return out
}

// Typed accessors
func Get[T any](g *Group, id, flag string) (T, error) {
	var zero T
	item, ok := g.Items()[flag]
	if !ok {
		return zero, fmt.Errorf("dynamic flag %q is not registered", flag)
	}
	v, ok := item.GetAny(id)
	if !ok {
		return zero, fmt.Errorf("no value set for --%s.%s.%s", g.Name(), id, flag)
	}
	typed, ok := v.(T)
	if !ok {
		return zero, fmt.Errorf("dynamic flag %q has unexpected type", flag)
	}
	return typed, nil
}

func MustGet[T any](g *Group, id, flag string) T {
	val, err := Get[T](g, id, flag)
	if err != nil {
		panic(err)
	}
	return val
}

func GetOrDefault[T any](g *Group, id, flag string) T {
	item, ok := g.Items()[flag]
	if !ok {
		panic(fmt.Sprintf("dynamic flag %q is not registered", flag))
	}
	v, _ := item.GetAny(id)
	typed, ok := v.(T)
	if !ok {
		panic(fmt.Sprintf("dynamic flag %q has unexpected type", flag))
	}
	return typed
}
