package tinyflags

type baseFlag struct {
	name       string
	short      string
	usage      string
	value      Value
	hidden     bool
	group      *mutualGroup
	disableEnv bool
	envKey     string
	deprecated string
	required   bool
	metavar    string
	allowed    []string // allowed values for this flag
}

type mutualGroup struct {
	name  string
	flags []*baseFlag
}

type Value interface {
	Set(string) error
	Get() any
	Default() string
	IsChanged() bool
}

type DynamicValue interface {
	Set(id string, val string) error
}

type dynamicItemValues interface {
	ValuesAny() map[string]any
}

// StrictBool is an optional interface implemented by flags that can be used
// without an explicit value (e.g. --verbose sets to true).
type StrictBool interface {
	IsStrictBool() bool
}

// SliceMarker is a marker interface for slice-type flags.
// It is used internally to distinguish scalar vs. slice values.
type SliceMarker interface {
	isSlice()
}

// HasDelimiter is implemented by slice flag types that allow custom
// delimiters for splitting input strings into elements.
type HasDelimiter interface {
	SetDelimiter(string) // Updates the delimiter used during Set()
}
