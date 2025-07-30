package dynamic

import "github.com/containeroo/tinyflags/internal/core"

// FlagSetRef is the subset of FlagSet needed by dynamic flags.
type FlagSetRef interface {
	RegisterFlag(name string, bf *core.BaseFlag)
	AttachToMutualGroup(*core.BaseFlag, string)
	GetMutualGroup(name string) *core.MutualExlusiveGroup
	MutualGroups() []*core.MutualExlusiveGroup
	DefaultDelimiter() string
	LookupFlag(name string) *core.BaseFlag
	GetRequireTogetherGroup(name string) *core.RequiredTogetherGroup
}
