package engine

import (
	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
)

// dynamicGroups returns all dynamic groups in desired order for internal consumers.
func (f *FlagSet) dynamicGroups() []*dynamic.Group {
	if f.sortGroups {
		return f.OrderedDynamicGroups()
	}
	return f.dynamicGroupsOrder
}

// staticFlags returns all static flags in desired order for internal consumers.
func (f *FlagSet) staticFlags() []*core.BaseFlag {
	if f.sortFlags {
		return f.OrderedStaticFlags()
	}
	return f.staticFlagsOrder
}

// staticSectionOrder returns named help sections in display order.
func (f *FlagSet) staticSectionOrder() []string {
	order := make([]string, 0, len(f.sectionOrder))
	seen := make(map[string]struct{}, len(f.sectionOrder))

	appendSection := func(name string) {
		if name == "" {
			return
		}
		if _, exists := seen[name]; exists {
			return
		}
		seen[name] = struct{}{}
		order = append(order, name)
	}

	for _, name := range f.sectionOrder {
		appendSection(name)
	}
	for _, flag := range f.staticFlagsOrder {
		appendSection(flag.Section)
	}

	return order
}
