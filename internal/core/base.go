package core

// BaseFlag contains all metadata and behavior for a single flag.
// It is the unified representation of scalar, slice, and dynamic flags.
type BaseFlag struct {
	Name       string       // Name is the long name of the flag (e.g. "verbose").
	Short      string       // Short is the short alias (single letter) for the flag (e.g. "v").
	Usage      string       // Usage is the description shown in help output.
	Value      Value        // Value holds the actual value implementation (scalar, slice, or dynamic).
	Hidden     bool         // Hidden marks the flag as hidden from help output.
	Group      *MutualGroup // Group points to the mutual exclusion group this flag belongs to (optional).
	DisableEnv bool         // DisableEnv disables environment variable resolution for this flag.
	EnvKey     string       // EnvKey is the custom environment variable name to use (if not empty).
	Deprecated string       // Deprecated is an optional message shown if the flag is deprecated.
	Required   bool         // Required marks the flag as required (must be explicitly set).
	Metavar    string       // Metavar is a placeholder name used in help output for the value (e.g. "FILE").
	Allowed    []string     // Allowed is an optional list of allowed string values (used for help only).
}
