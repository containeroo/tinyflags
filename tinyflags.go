// Package tinyflags provides a high-level API for defining and parsing
// CLI flags with support for dynamic groups, custom types, and rich usage output.
package tinyflags

import (
	"io"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
	"github.com/containeroo/tinyflags/internal/engine"
)

// ErrorHandling is re-exported from internal engine.
type ErrorHandling = engine.ErrorHandling

const (
	ContinueOnError = engine.ContinueOnError
	ExitOnError     = engine.ExitOnError
	PanicOnError    = engine.PanicOnError
)

// Common user errors.
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

// FlagSet is the public wrapper around the internal flag engine.
// It provides the main API for defining flags and parsing CLI input.
type FlagSet struct {
	impl *engine.FlagSet

	// Usage is an optional function that overrides default help output.
	// It is called automatically when --help is requested.
	Usage func()
}

// NewFlagSet creates a new FlagSet with the given name and error handling mode.
func NewFlagSet(name string, handling ErrorHandling) *FlagSet {
	return &FlagSet{impl: engine.NewFlagSet(name, handling)}
}

// Parse parses the provided CLI args and environment variables.
func (f *FlagSet) Parse(args []string) error {
	if f.Usage != nil {
		f.impl.Usage = f.Usage
	}
	return f.impl.Parse(args)
}

// Name returns the program name.
func (f *FlagSet) Name() string { return f.impl.Name() }

// Version sets the --version output string.
func (f *FlagSet) Version(s string) { f.impl.Version(s) }

// EnvPrefix sets a prefix for resolving environment variables.
func (f *FlagSet) EnvPrefix(s string) { f.impl.EnvPrefix(s) }

// Authors adds author information to the help header.
func (f *FlagSet) Authors(s string) { f.impl.Authors(s) }

// Title sets the program title for help output.
func (f *FlagSet) Title(s string) { f.impl.Title(s) }

// Description adds a paragraph to the help description.
func (f *FlagSet) Description(s string) { f.impl.Description(s) }

// Note appends a note to the bottom of the help output.
func (f *FlagSet) Note(s string) { f.impl.Note(s) }

// DisableHelp disables automatic registration of the --help flag.
func (f *FlagSet) DisableHelp() { f.impl.DisableHelp() }

// DisableVersion disables automatic registration of the --version flag.
func (f *FlagSet) DisableVersion() { f.impl.DisableVersion() }

// SortedFlags enables sorting of static flag help output.
func (f *FlagSet) SortedFlags() { f.impl.SortedFlags(true) }

// SortedGroups enables sorting of dynamic groups in help output.
func (f *FlagSet) SortedGroups() { f.impl.SortedGroups(true) }

// SetOutput changes the writer used for help and error messages.
func (f *FlagSet) SetOutput(w io.Writer) { f.impl.SetOutput(w) }

// Output returns the current output writer.
func (f *FlagSet) Output() io.Writer { return f.impl.Output() }

// IgnoreInvalidEnv skips invalid values from environment variables.
func (f *FlagSet) IgnoreInvalidEnv(b bool) { f.impl.IgnoreInvalidEnv(b) }

// SetGetEnvFn overrides how environment variables are looked up.
func (f *FlagSet) SetGetEnvFn(fn func(string) string) { f.impl.SetGetEnvFn(fn) }

// Globaldelimiter sets the default slice delimiter for all slice flags.
func (f *FlagSet) Globaldelimiter(s string) { f.impl.Globaldelimiter(s) }

// DefaultDelimiter returns the current slice delimiter.
func (f *FlagSet) DefaultDelimiter() string {
	return f.impl.DefaultDelimiter()
}

// RequirePositional enforces a minimum number of positional arguments.
func (f *FlagSet) RequirePositional(n int) { f.impl.RequirePositional(n) }

// Args returns all leftover positional arguments.
func (f *FlagSet) Args() []string { return f.impl.Args() }

// Arg returns the nth positional argument, if present.
func (f *FlagSet) Arg(i int) (string, bool) { return f.impl.Arg(i) }

// GetGroup retrieves a named mutual-exclusion group.
func (f *FlagSet) GetGroup(name string) *core.MutualGroup { return f.impl.GetGroup(name) }

// Groups returns all defined mutual-exclusion groups.
func (f *FlagSet) Groups() []*core.MutualGroup {
	return f.impl.Groups()
}

// AttachToGroup assigns a static flag to a mutual-exclusion group.
func (f *FlagSet) AttachToGroup(bf *core.BaseFlag, group string) {
	f.impl.AttachToGroup(bf, group)
}

// LookupFlag returns the static flag registered under the given name.
func (f *FlagSet) LookupFlag(name string) *core.BaseFlag {
	return f.impl.LookupFlag(name)
}

// DescriptionMaxLen sets the wrapping width for help descriptions.
func (f *FlagSet) DescriptionMaxLen(n int) { f.impl.DescriptionMaxLen(n) }

// DescriptionIndent sets the indentation width for help descriptions.
func (f *FlagSet) DescriptionIndent(n int) { f.impl.DescriptionIndent(n) }

// PrintDefaults prints static and dynamic flags with their default values.
func (f *FlagSet) PrintDefaults(w io.Writer, width int) { f.impl.PrintDefaults(w, width) }

// PrintUsage writes help output in the given format mode.
func (f *FlagSet) PrintUsage(w io.Writer, mode FlagPrintMode) { f.impl.PrintUsage(w, mode) }

// PrintTitle writes the help title section.
func (f *FlagSet) PrintTitle(w io.Writer) { f.impl.PrintTitle(w) }

// PrintNotes writes additional help notes, wrapped at the given width.
func (f *FlagSet) PrintNotes(w io.Writer, width int) { f.impl.PrintNotes(w, width) }

// PrintDescription writes the description section, wrapped at the given width.
func (f *FlagSet) PrintDescription(w io.Writer, width int) { f.impl.PrintDescription(w, width) }

// PrintAuthors writes the authors section.
func (f *FlagSet) PrintAuthors(w io.Writer) { f.impl.PrintAuthors(w) }

// DynamicGroup creates or returns a dynamic flag group by name.
func (f *FlagSet) DynamicGroup(name string) *dynamic.Group {
	return f.impl.DynamicGroup(name)
}

// DynamicGroups returns all registered dynamic flag groups.
func (f *FlagSet) DynamicGroups() []*dynamic.Group {
	return f.impl.DynamicGroups()
}

// GetDynamic returns the typed value for a dynamic flag field by ID.
func GetDynamic[T any](group *dynamic.Group, id string, flag string) (T, error) {
	return dynamic.Get[T](group, id, flag)
}

// MustGetDynamic returns the typed value for a dynamic field or panics if unset.
func MustGetDynamic[T any](group *dynamic.Group, id string, flag string) T {
	return dynamic.MustGet[T](group, id, flag)
}

// GetOrDefaultDynamic returns the value for a dynamic field or its default.
func GetOrDefaultDynamic[T any](group *dynamic.Group, id, flag string) T {
	return dynamic.GetOrDefault[T](group, id, flag)
}
