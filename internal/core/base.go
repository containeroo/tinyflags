package core

// BaseFlag holds metadata for a single flag.
type BaseFlag struct {
	Name         string       // Long name (e.g. "verbose").
	Short        string       // Short alias (single letter, e.g. "v").
	Usage        string       // Brief description shown in help.
	Value        Value        // Underlying value handler (scalar, slice, or dynamic).
	Hidden       bool         // If true, omit from help.
	DisableEnv   bool         // If true, disallow ENV lookup.
	HideEnv      bool         // If true, hide ENV key from help.
	EnvKey       string       // Custom environment variable (overrides derived key).
	Deprecated   string       // If non‐empty, show deprecation notice.
	Required     bool         // Mark flag as required.
	HideRequired bool         // Hide “(Required)” in help.
	Placeholder  string       // Placeholder for the value (e.g. "FILE").
	Allowed      []string     // Allowed string values (help only).
	Group        *MutualGroup // Mutual exclusion group membership.
}
