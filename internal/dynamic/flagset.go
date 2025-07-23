package dynamic

import "github.com/containeroo/tinyflags/internal/core"

// FlagSetRef is the subset of FlagSet needed by dynamic flags.
type FlagSetRef interface {
	RegisterDynamic(group, field string, v core.DynamicValue) error
	RegisterFlag(name string, bf *core.BaseFlag)
	AttachToGroup(*core.BaseFlag, string)
	GetGroup(name string) *core.MutualGroup
	Groups() []*core.MutualGroup
	DefaultDelimiter() string
}
