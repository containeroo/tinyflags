package dynamic

import "github.com/containeroo/tinyflags/internal/core"

// FlagSetRef is the subset of FlagSet needed by dynamic flags.
type FlagSetRef interface {
	RegisterFlag(name string, bf *core.BaseFlag)
	AttachToOneOfGroup(*core.BaseFlag, string)
	GetOneOfGroup(name string) *core.OneOfGroupGroup
	OneOfGroups() []*core.OneOfGroupGroup
	DefaultDelimiter() string
	LookupFlag(name string) *core.BaseFlag
	GetAllOrNoneGroup(name string) *core.AllOrNoneGroup
}
