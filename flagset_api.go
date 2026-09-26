package tinyflags

import (
	"strings"

	"github.com/containeroo/tinyflags/internal/engine"
)

// Parse processes the given CLI args and populates all registered flags.
func (f *FlagSet) Parse(args []string) error {
	if f.Usage != nil {
		f.impl.Usage = f.Usage
	}
	return f.impl.Parse(args)
}

// BeforeParse installs a hook to mutate arguments before parsing.
func (f *FlagSet) BeforeParse(fn func([]string) ([]string, error)) {
	f.impl.BeforeParse(fn)
}

// Validate registers a cross-flag validator. Validators run in registration
// order after parsing, environment loading, defaults, and built-in checks.
func (f *FlagSet) Validate(fn func() error) {
	f.impl.Validate(fn)
}

// OnUnknownFlag installs a handler for unknown flags. Return nil to ignore.
func (f *FlagSet) OnUnknownFlag(fn func(string) error) {
	f.impl.OnUnknownFlag(fn)
}

// Name returns the flag set's name.
func (f *FlagSet) Name() string { return f.impl.Name() }

// Version sets the --version string.
func (f *FlagSet) Version(s string) { f.impl.Version(s) }

// VersionText sets the --version text.
func (f *FlagSet) VersionText(s string) { f.impl.VersionText(s) }

// EnvPrefix sets the prefix for automatic environment variable lookups.
// Commands inherit the nearest ancestor prefix unless explicitly configured.
// An empty prefix stops inheritance and disables automatic lookups; explicit
// per-flag Env keys still work. Descendants inherit this setting.
func (f *FlagSet) EnvPrefix(s string) { f.impl.EnvPrefix(s) }

// SetEnvKeyFunc sets a function to derive env keys from prefix+flag name.
func (f *FlagSet) SetEnvKeyFunc(fn engine.EnvKeyFunc) { f.impl.SetEnvKeyFunc(fn) }

// EnvKeyForFlag derives the env key for a flag.
func (f *FlagSet) EnvKeyForFlag(name string) string { return f.impl.EnvKeyForFlag(name) }

// NewReplacerEnvKeyFunc builds an EnvKeyFunc that:
// - returns "" when prefix is empty
// - applies the given replacer to the flag name
// - joins prefix + "_" + transformed name
// - upper-cases the result (if upper is true)
func (f *FlagSet) NewReplacerEnvKeyFunc(replacer *strings.Replacer, upper bool) engine.EnvKeyFunc {
	return engine.NewReplacerEnvKeyFunc(replacer, upper)
}

// FirstChanged returns the value of the first changed flag in the given order.
// If no flag was changed, it returns defaultValue and false.
func FirstChanged[T any](defaultValue T, flags ...Flag[T]) (T, bool) {
	return engine.FirstChanged(defaultValue, flags...)
}

// IgnoreInvalidEnv disables errors for unrecognized environment values.
func (f *FlagSet) IgnoreInvalidEnv(b bool) { f.impl.IgnoreInvalidEnv(b) }

// SetGetEnvFn overrides the function used to look up environment variables.
func (f *FlagSet) SetGetEnvFn(fn func(string) string) { f.impl.SetGetEnvFn(fn) }

// GlobalDelimiter sets the delimiter used for all slice flags.
func (f *FlagSet) GlobalDelimiter(s string) { f.impl.GlobalDelimiter(s) }

// DefaultDelimiter returns the delimiter used for slice flags.
func (f *FlagSet) DefaultDelimiter() string { return f.impl.DefaultDelimiter() }

// RequirePositional sets how many positional arguments must be present.
func (f *FlagSet) RequirePositional(n int) { f.impl.RequirePositional(n) }

// Args returns all leftover positional arguments.
func (f *FlagSet) Args() []string { return f.impl.Args() }

// Arg returns the i-th positional argument and whether it exists.
func (f *FlagSet) Arg(i int) (string, bool) { return f.impl.Arg(i) }

// SetPositionalValidate sets a function to validate positional arguments.
func (f *FlagSet) SetPositionalValidate(fn func(string) error) { f.impl.SetPositionalValidate(fn) }

// SetPositionalFinalize sets a function to finalize positional arguments.
func (f *FlagSet) SetPositionalFinalize(fn func(string) string) { f.impl.SetPositionalFinalize(fn) }

// OverriddenValues returns all flags that were explicitly set (args or env).
// Dynamic flags use the key format "group.id.flag".
func (f *FlagSet) OverriddenValues() map[string]any { return f.impl.OverriddenValues() }

// Origin returns the exact input that supplied name in the most recent parse.
// Unoverridden values report the default origin.
func (f *FlagSet) Origin(name string) ValueOrigin { return f.impl.Origin(name) }

// OverriddenOrigins returns the exact origin of each value explicitly set by a flag or environment variable.
// Dynamic flags use the key format "group.id.flag".
func (f *FlagSet) OverriddenOrigins() map[string]ValueOrigin { return f.impl.OverriddenOrigins() }
