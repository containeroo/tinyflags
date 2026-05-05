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
