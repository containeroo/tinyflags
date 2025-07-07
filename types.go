package tinyflags

// baseFlag stores metadata and state for a single registered flag.
type baseFlag struct {
	name       string       // Long name of the flag (e.g. "port")
	short      string       // Optional shorthand name (e.g. "p")
	usage      string       // Help text shown in usage
	value      Value        // Concrete value implementation for parsing and storage
	hidden     bool         // Whether to hide this flag from help output
	group      *mutualGroup // Optional mutual exclusion group this flag belongs to
	disableEnv bool         // Whether to disable env key lookup
	envKey     string       // Optional environment variable name to override this flag
	deprecated string       // Deprecation reason, if any
	required   bool         // Whether this flag is required
	metavar    string       // Placeholder name shown in help output (e.g. "FILE")
	allowed    []string     // allowed values for this flag
}

// mutualGroup represents a group of flags that are mutually exclusive.
// Only one flag in the group can be set at a time.
type mutualGroup struct {
	name  string      // Name of the group (used for error messages and grouping)
	flags []*baseFlag // All flags that are members of the group
}

// Value is the interface implemented by all flag value types.
// It provides a way to parse and retrieve values from strings.
type Value interface {
	Set(string) error // Parses and sets the value from a string input
	Get() any         // Returns the current value as interface{}
	Default() string  // Returns the default value as a string
	IsChanged() bool  // Reports whether the flag value was explicitly set
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
