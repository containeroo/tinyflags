package core

// Registry manages static flags, dynamic groups, and defaults.
type Registry interface {
	RegisterFlag(name string, bf *BaseFlag)
	GetOneOfGroup(name string) *OneOfGroupGroup
	OneOfGroups() []*OneOfGroupGroup
	DefaultDelimiter() string
	GetAllOrNoneGroup(name string) *AllOrNoneGroup
}

// Value parses and holds a single CLI value.
type Value interface {
	Set(string) error // parse from string
	Get() any         // retrieve stored value
	Changed() bool    // was flag explicitly set?
	Default() string  // default value as string
}

// DefaultFinalizer applies default-only finalization for unset values.
type DefaultFinalizer interface {
	ApplyDefaultFinalize()
}

// GroupItem holds a single flag and its value for a dynamic group.
type GroupItem struct {
	Value DynamicValue
	Flag  *BaseFlag
}

// DynamicValue accepts keyed values (e.g. --http.alpha.port).
type DynamicValue interface {
	Set(id, val string) error
	FieldName() string
	GetAny(id string) (any, bool)
}

// DynamicItemValues exposes all parsed dynamic entries.
type DynamicItemValues interface {
	ValuesAny() map[string]any
}

// StrictBool supports --flag / --no-flag syntax.
type StrictBool interface {
	IsStrictBool() bool
}

// SliceMarker tags slice-type flags (no methods).
type SliceMarker interface {
	IsSlice()
}

// Incrementable allows repeated use to increment a counter flag.
type Incrementable interface {
	Value
	Increment() error
}
