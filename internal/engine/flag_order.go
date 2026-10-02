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

// staticSectionOrder returns static help sections in display order.
// Named sections follow the explicit order first, then first-registration order.
// The unnamed section is appended last unless explicitly positioned.
func (f *FlagSet) staticSectionOrder() []string {
	order := make([]string, 0, len(f.sectionOrder)+1)
	seen := make(map[string]struct{}, len(f.sectionOrder)+1)
	unnamedExplicit := false
	hasUnnamed := false

	appendSection := func(name string) {
		if _, exists := seen[name]; exists {
			return
		}
		seen[name] = struct{}{}
		order = append(order, name)
	}

	for _, name := range f.sectionOrder {
		if name == "" {
			unnamedExplicit = true
		}
		appendSection(name)
	}

	for _, flag := range f.staticFlagsOrder {
		if flag.Section == "" {
			hasUnnamed = true
			continue
		}
		appendSection(flag.Section)
	}

	if hasUnnamed && !unnamedExplicit {
		appendSection("")
	}

	return order
}
