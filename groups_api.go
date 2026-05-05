package tinyflags

import (
	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
)

// AttachGroupToAllOrNone nests one AllOrNone group into another.
func (f *FlagSet) AttachGroupToAllOrNone(parent, child string) {
	f.impl.AttachGroupToAllOrNone(parent, child)
}

// AttachGroupToOneOf adds an AllOrNone group to a OneOf group.
func (f *FlagSet) AttachGroupToOneOf(group, name string) {
	f.impl.AttachGroupToOneOf(group, name)
}

// OneOfGroups returns all registered OneOfGroups groups.
func (f *FlagSet) OneOfGroups() []*core.OneOfGroupGroup { return f.impl.OneOfGroups() }

// AddOneOfGroup adds a AddOneOfGroup group.
func (f *FlagSet) AddOneOfGroup(name string, g *core.OneOfGroupGroup) {
	f.impl.AddOneOfGroup(name, g)
}

// GetOneOfGroup retrieves or creates a named OneOfGroup group.
func (f *FlagSet) GetOneOfGroup(name string) *core.OneOfGroupGroup {
	return f.impl.GetOneOfGroup(name)
}

// AttachToOneOfGroup attaches a static flag to a OneOfGroup group.
func (f *FlagSet) AttachToOneOfGroup(flag *core.BaseFlag, group string) {
	f.impl.AttachToOneOfGroup(flag, group)
}

// AllOrNoneGroup returns all registered AllOrNoneGroup group.
func (f *FlagSet) AllOrNoneGroup() []*core.AllOrNoneGroup {
	return f.impl.AllOrNoneGroups()
}

// AllOrNoneGroups returns all registered AllOrNone groups.
func (f *FlagSet) AllOrNoneGroups() []*core.AllOrNoneGroup {
	return f.impl.AllOrNoneGroups()
}

// AddAllOrNoneGroup adds a AllOrNoneGroup group.
func (f *FlagSet) AddAllOrNoneGroup(name string, g *core.AllOrNoneGroup) {
	f.impl.AddAllOrNoneGroup(name, g)
}

// GetAllOrNoneGroup retrieves or creates a named AllOrNoneGroup group.
func (f *FlagSet) GetAllOrNoneGroup(name string) *core.AllOrNoneGroup {
	return f.impl.GetAllOrNoneGroup(name)
}

// AttachToAllOrNoneGroup attaches a flag to a AllOrNoneGroup group.
func (f *FlagSet) AttachToAllOrNoneGroup(flag *core.BaseFlag, group string) {
	f.impl.AttachToAllOrNoneGroup(flag, group)
}

// LookupFlag retrieves a static flag by name.
func (f *FlagSet) LookupFlag(name string) *core.BaseFlag {
	return f.impl.LookupFlag(name)
}

// DynamicGroup registers or retrieves a dynamic group by name.
func (f *FlagSet) DynamicGroup(name string) *dynamic.Group {
	return f.impl.DynamicGroup(name)
}

// DynamicGroups returns all dynamic groups in registration order.
func (f *FlagSet) DynamicGroups() []*dynamic.Group {
	return f.impl.DynamicGroups()
}

// GetDynamic retrieves the value for a dynamic flag by group, id, and field name.
func GetDynamic[T any](group *dynamic.Group, id, flag string) (T, error) {
	return dynamic.Get[T](group, id, flag)
}

// MustGetDynamic retrieves the value for a dynamic flag or panics.
func MustGetDynamic[T any](group *dynamic.Group, id, flag string) T {
	return dynamic.MustGet[T](group, id, flag)
}

// GetOrDefaultDynamic returns the value for a dynamic flag or its default.
func GetOrDefaultDynamic[T any](group *dynamic.Group, id, flag string) T {
	return dynamic.GetOrDefault[T](group, id, flag)
}
