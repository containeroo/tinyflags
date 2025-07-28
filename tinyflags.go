package tinyflags

import (
	"io"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
	"github.com/containeroo/tinyflags/internal/engine"
)

// ErrorHandling is re-exported from internal engine
type ErrorHandling = engine.ErrorHandling

const (
	ContinueOnError = engine.ContinueOnError
	ExitOnError     = engine.ExitOnError
	PanicOnError    = engine.PanicOnError
)

// HelpRequested, VersionRequested, and helpers
type (
	HelpRequested    = engine.HelpRequested
	VersionRequested = engine.VersionRequested
)

type FlagPrintMode = engine.FlagPrintMode

const (
	PrintNone  = engine.PrintNone
	PrintFlags = engine.PrintFlags
	PrintShort = engine.PrintShort
	PrintLong  = engine.PrintLong
	PrintBoth  = engine.PrintBoth
)

var (
	IsHelpRequested    = engine.IsHelpRequested
	IsVersionRequested = engine.IsVersionRequested
	RequestHelp        = engine.RequestHelp
	RequestVersion     = engine.RequestVersion
)

// FlagSet is the public wrapper around internal engine.FlagSet
// It provides the full user-facing API for defining and parsing CLI flags.
type FlagSet struct {
	impl *engine.FlagSet

	// Usage can be overridden by the user to customize help output.
	// It is executed when --help is triggered or from user code.
	Usage func()
}

// NewFlagSet creates a new flag set.
func NewFlagSet(name string, handling ErrorHandling) *FlagSet {
	return &FlagSet{impl: engine.NewFlagSet(name, handling)}
}

// Parse processes args and environment variables.
func (f *FlagSet) Parse(args []string) error {
	if f.Usage != nil {
		f.impl.Usage = f.Usage
	}
	return f.impl.Parse(args)
}

// Name returns the name of the application.
func (f *FlagSet) Name() string { return f.impl.Name() }

// Version sets the version string for --version output.
func (f *FlagSet) Version(s string) { f.impl.Version(s) }

// EnvPrefix sets a prefix for deriving environment-variable names.
func (f *FlagSet) EnvPrefix(s string) { f.impl.EnvPrefix(s) }

// Authors adds author information to the help text.
func (f *FlagSet) Authors(s string) { f.impl.Authors(s) }

// Title sets the program title in the help header.
func (f *FlagSet) Title(s string) { f.impl.Title(s) }

// Description adds a description paragraph to the help.
func (f *FlagSet) Description(s string) { f.impl.Description(s) }

// Note appends a note paragraph to the help footer.
func (f *FlagSet) Note(s string) { f.impl.Note(s) }

// DisableHelp turns off automatic help flag registration.
func (f *FlagSet) DisableHelp() { f.impl.DisableHelp() }

// DisableVersion turns off automatic version flag registration.
func (f *FlagSet) DisableVersion() { f.impl.DisableVersion() }

// Sorted enables or disables sorted flag output.
func (f *FlagSet) Sorted(b bool) { f.impl.Sorted(b) }

// SetOutput redirects all help and error output.
func (f *FlagSet) SetOutput(w io.Writer) { f.impl.SetOutput(w) }

// Output returns the current output writer.
func (f *FlagSet) Output() io.Writer { return f.impl.Output() }

// IgnoreInvalidEnv skips invalid environment-variable values.
func (f *FlagSet) IgnoreInvalidEnv(b bool) { f.impl.IgnoreInvalidEnv(b) }

// SetGetEnvFn replaces how environment variables are looked up.
func (f *FlagSet) SetGetEnvFn(fn func(string) string) { f.impl.SetGetEnvFn(fn) }

// Globaldelimiter sets the default delimiter for slice flags.
func (f *FlagSet) Globaldelimiter(s string) { f.impl.Globaldelimiter(s) }

// GetGroup returns the named mutual-exclusion group.
func (f *FlagSet) GetGroup(name string) *core.MutualGroup { return f.impl.GetGroup(name) }

// RequirePositional enforces a minimum number of positional args.
func (f *FlagSet) RequirePositional(n int) { f.impl.RequirePositional(n) }

// Args returns leftover positional arguments.
func (f *FlagSet) Args() []string { return f.impl.Args() }

// Arg returns the nth positional argument, if present.
func (f *FlagSet) Arg(i int) (string, bool) { return f.impl.Arg(i) }

// DescriptionMaxLen sets the max width for description text.
func (f *FlagSet) DescriptionMaxLen(n int) { f.impl.DescriptionMaxLen(n) }

// DescriptionIndent sets the indent width for descriptions.
func (f *FlagSet) DescriptionIndent(n int) { f.impl.DescriptionIndent(n) }

// PrintDefaults prints all defined flags and their defaults.
func (f *FlagSet) PrintDefaults(w io.Writer, width int) { f.impl.PrintDefaults(w, width) }

// func (f *FlagSet) PrintDynamicDefaults(w io.Writer, width int) { f.impl.PrintDynamicDefaults(w, width) }
// PrintUsage writes usage text in the specified mode.
func (f *FlagSet) PrintUsage(w io.Writer, mode FlagPrintMode) { f.impl.PrintUsage(w, mode) }

// PrintTitle writes the help title.
func (f *FlagSet) PrintTitle(w io.Writer) { f.impl.PrintTitle(w) }

// PrintNotes writes help notes wrapped at width.
func (f *FlagSet) PrintNotes(w io.Writer, width int) { f.impl.PrintNotes(w, width) }

// PrintDescription writes the description wrapped at width.
func (f *FlagSet) PrintDescription(w io.Writer, width int) { f.impl.PrintDescription(w, width) }

// PrintAuthors writes the authors heading.
func (f *FlagSet) PrintAuthors(w io.Writer) { f.impl.PrintAuthors(w) }

// DynamicGroup returns or creates a dynamic flag group.
func (f *FlagSet) DynamicGroup(name string) *dynamic.Group {
	return f.impl.DynamicGroup(name)
}

// // GetDynamic retrieves a dynamic field value for a given ID with the correct type.
func GetDynamic[T any](group *dynamic.Group, id string, flag string) (T, error) {
	return dynamic.Get[T](group, id, flag)
}

// // MustGetDynamic panics if the field or id is missing.
func MustGetDynamic[T any](group *dynamic.Group, id string, flag string) T {
	return dynamic.MustGet[T](group, id, flag)
}

// // GetOrDefaultDynamic returns the default value if the field or id is missing.
func GetOrDefaultDynamic[T any](group *dynamic.Group, id, flag string) T {
	return dynamic.GetOrDefault[T](group, id, flag)
}

// DynamicGroups returns all dynamic flag groups.
func (f *FlagSet) DynamicGroups() []*dynamic.Group {
	return f.impl.DynamicGroups()
}

// DefaultDelimiter returns the slice-value separator.
func (f *FlagSet) DefaultDelimiter() string {
	return f.impl.DefaultDelimiter()
}

// Groups returns all mutual-exclusion groups.
func (f *FlagSet) Groups() []*core.MutualGroup {
	return f.impl.Groups()
}

// AttachToGroup connects a flag to a mutual-exclusion group.
func (f *FlagSet) AttachToGroup(bf *core.BaseFlag, group string) {
	f.impl.AttachToGroup(bf, group)
}

// LookupFlag retrieves a registered static flag.
func (f *FlagSet) LookupFlag(name string) *core.BaseFlag {
	return f.impl.LookupFlag(name)
}
