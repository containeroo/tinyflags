package core

// BaseFlag holds metadata for a single flag.
type BaseFlag struct {
	Name         string           // Long name (e.g. "verbose").
	Short        string           // Short alias (single letter, e.g. "v").
	Usage        string           // Brief description shown in help.
	Value        Value            // Underlying value handler (scalar, slice, or dynamic).
	Hidden       bool             // If true, omit from help.
	DisableEnv   bool             // If true, disallow ENV lookup.
	EnvKey       string           // Custom environment variable (overrides derived key).
	HideEnv      bool             // If true, hide ENV key from help.
	Deprecated   string           // If non‐empty, show deprecation notice.
	Required     bool             // Mark flag as required.
	HideRequired bool             // Hide “(Required)” in help.
	Placeholder  string           // Placeholder for the value (e.g. "FILE").
	Allowed      []string         // Allowed string values (help only).
	HideAllowed  bool             // Hide allowed values from help.
	OneOfGroup   *OneOfGroupGroup // OneOfGroup group membership.
	AllOrNone    *AllOrNoneGroup  // AllOrNone group membership.
	Requires     []string         // Names of flags this flag requires
	HideRequires bool             // Hide “(Requires)” in help.
}
