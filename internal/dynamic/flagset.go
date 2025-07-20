package dynamic

import "github.com/containeroo/tinyflags/internal/core"

// FlagSetRef defines the interface needed by dynamic flag types
// to register themselves and attach to mutual groups.
type FlagSetRef interface {
	RegisterDynamic(group, field string, v core.DynamicValue) error
	AttachToGroup(*core.BaseFlag, string)
	DefaultDelimiter() string
}
