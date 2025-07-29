// Package tinyflags provides a high-level API for defining and parsing
// CLI flags with support for dynamic groups, custom types, and rich usage output.
package tinyflags

import (
	"io"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
	"github.com/containeroo/tinyflags/internal/engine"
)

// ErrorHandling defines how parsing errors are handled.
type ErrorHandling = engine.ErrorHandling

const (
	ContinueOnError = engine.ContinueOnError // Continue and return error
	ExitOnError     = engine.ExitOnError     // Exit with error message
	PanicOnError    = engine.PanicOnError    // Panic on error
)

// Common user-triggered exit conditions.
type (
	HelpRequested    = engine.HelpRequested
	VersionRequested = engine.VersionRequested
)

var (
	IsHelpRequested    = engine.IsHelpRequested
	IsVersionRequested = engine.IsVersionRequested
	RequestHelp        = engine.RequestHelp
	RequestVersion     = engine.RequestVersion
)

// FlagPrintMode controls how the usage line is rendered.
type FlagPrintMode = engine.FlagPrintMode

const (
	PrintNone  = engine.PrintNone  // Omits usage line entirely
	PrintFlags = engine.PrintFlags // Prints: [flags]
	PrintShort = engine.PrintShort // Prints: -v
	PrintLong  = engine.PrintLong  // Prints: --verbose
	PrintBoth  = engine.PrintBoth  // Prints: -v|--verbose
)

// Exported types for advanced access.
type (
	DynamicGroup = dynamic.Group // Dynamic group of instance-scoped flags
	StaticFlag   = core.BaseFlag // Static flag definition metadata
)

// FlagSet is the user-facing flag parser and usage configurator.
type FlagSet struct {
	impl  *engine.FlagSet
	Usage func() // Optional custom usage function
}

// NewFlagSet creates a new flag set with the given name and error handling mode.
func NewFlagSet(name string, handling ErrorHandling) *FlagSet {
	return &FlagSet{impl: engine.NewFlagSet(name, handling)}
}

// Parse processes the given CLI args and populates all registered flags.
func (f *FlagSet) Parse(args []string) error {
	if f.Usage != nil {
		f.impl.Usage = f.Usage
	}
	return f.impl.Parse(args)
}

// Name returns the flag set's name.
func (f *FlagSet) Name() string { return f.impl.Name() }

// Version sets the --version string.
func (f *FlagSet) Version(s string) { f.impl.Version(s) }

// EnvPrefix sets a prefix for all environment variables.
func (f *FlagSet) EnvPrefix(s string) { f.impl.EnvPrefix(s) }

// Title sets the main title shown in usage output.
func (f *FlagSet) Title(s string) { f.impl.Title(s) }

// Authors sets the list of authors printed in help output.
func (f *FlagSet) Authors(s string) { f.impl.Authors(s) }

// Description sets the top description section of the help output.
func (f *FlagSet) Description(s string) { f.impl.Description(s) }

// Note sets the bottom note section of the help output.
func (f *FlagSet) Note(s string) { f.impl.Note(s) }

// DisableHelp disables the automatic --help flag.
func (f *FlagSet) DisableHelp() { f.impl.DisableHelp() }

// DisableVersion disables the automatic --version flag.
func (f *FlagSet) DisableVersion() { f.impl.DisableVersion() }

// HideEnvs disables all env-var annotations in help output.
func (f *FlagSet) HideEnvs() { f.impl.HideEnvs() }

// SortedFlags enables sorted help output for static flags.
func (f *FlagSet) SortedFlags() { f.impl.SortedFlags(true) }

// SortedGroups enables sorted help output for dynamic groups.
func (f *FlagSet) SortedGroups() { f.impl.SortedGroups(true) }

// SetOutput changes the destination writer for usage and error messages.
func (f *FlagSet) SetOutput(w io.Writer) { f.impl.SetOutput(w) }

// Output returns the current output writer.
func (f *FlagSet) Output() io.Writer { return f.impl.Output() }

// IgnoreInvalidEnv disables errors for unrecognized environment values.
func (f *FlagSet) IgnoreInvalidEnv(b bool) { f.impl.IgnoreInvalidEnv(b) }

// SetGetEnvFn overrides the function used to look up environment variables.
func (f *FlagSet) SetGetEnvFn(fn func(string) string) { f.impl.SetGetEnvFn(fn) }

// Globaldelimiter sets the delimiter used for all slice flags.
func (f *FlagSet) Globaldelimiter(s string) { f.impl.Globaldelimiter(s) }

// DefaultDelimiter returns the delimiter used for slice flags.
func (f *FlagSet) DefaultDelimiter() string { return f.impl.DefaultDelimiter() }

// RequirePositional sets how many positional arguments must be present.
func (f *FlagSet) RequirePositional(n int) { f.impl.RequirePositional(n) }

// Args returns all leftover positional arguments.
func (f *FlagSet) Args() []string { return f.impl.Args() }

// Arg returns the i-th positional argument and whether it exists.
func (f *FlagSet) Arg(i int) (string, bool) { return f.impl.Arg(i) }

// GetGroup retrieves or creates a named mutual exclusion group.
func (f *FlagSet) GetGroup(name string) *core.MutualGroup { return f.impl.GetGroup(name) }

// Groups returns all registered mutual exclusion groups.
func (f *FlagSet) Groups() []*core.MutualGroup { return f.impl.Groups() }

// AttachToGroup attaches a static flag to a mutual exclusion group.
func (f *FlagSet) AttachToGroup(flag *core.BaseFlag, group string) {
	f.impl.AttachToGroup(flag, group)
}

// LookupFlag retrieves a static flag by name.
func (f *FlagSet) LookupFlag(name string) *core.BaseFlag {
	return f.impl.LookupFlag(name)
}

// DynamicGroup registers or retrieves a dynamic group by name.
func (f *FlagSet) DynamicGroup(name string) *dynamic.Group {
	return f.impl.DynamicGroup(name)
}

// DynamicGroups returns all dynamic groups in registration order.
func (f *FlagSet) DynamicGroups() []*dynamic.Group {
	return f.impl.DynamicGroups()
}

// PrintUsage renders the top usage line.
func (f *FlagSet) PrintUsage(w io.Writer, mode FlagPrintMode) {
	f.impl.PrintUsage(w, mode)
}

// PrintTitle renders the title above all help content.
func (f *FlagSet) PrintTitle(w io.Writer) { f.impl.PrintTitle(w) }

// PrintAuthors renders the author line if set.
func (f *FlagSet) PrintAuthors(w io.Writer) { f.impl.PrintAuthors(w) }

// PrintDescription renders the full description block.
func (f *FlagSet) PrintDescription(w io.Writer, indent, width int) {
	f.impl.PrintDescription(w, indent, width)
}

// PrintStaticDefaults renders all static flag usage lines.
func (f *FlagSet) PrintStaticDefaults(w io.Writer, indent, col, width int) {
	f.impl.PrintStaticDefaults(w, indent, col, width)
}

// PrintDynamicDefaults renders all dynamic flag usage lines.
func (f *FlagSet) PrintDynamicDefaults(w io.Writer, indent, col, width int) {
	f.impl.PrintDynamicDefaults(w, indent, col, width)
}

// PrintNotes renders the notes section, if configured.
func (f *FlagSet) PrintNotes(w io.Writer, indent, width int) {
	f.impl.PrintNotes(w, indent, width)
}

// SetDescIndent sets the indentation for the description block.
func (f *FlagSet) SetDescIndent(n int) { f.impl.SetDescIndent(n) }

// DescIndent returns the current indent used for the description block.
func (f *FlagSet) DescIndent() int { return f.impl.DescIndent() }

// SetDescWidth sets the wrapping width for the description block.
func (f *FlagSet) SetDescWidth(max int) { f.impl.SetDescWidth(max) }

// DescWidth returns the wrapping width for the description block.
func (f *FlagSet) DescWidth() int { return f.impl.DescWidth() }

// SetStaticUsageIndent sets the indentation for static flag usage lines.
func (f *FlagSet) SetStaticUsageIndent(n int) { f.impl.SetStaticUsageIndent(n) }

// StaticUsageIndent returns the static usage indentation.
func (f *FlagSet) StaticUsageIndent() int { return f.impl.StaticUsageIndent() }

// SetStaticUsageColumn sets the column at which static flag descriptions begin.
func (f *FlagSet) SetStaticUsageColumn(col int) { f.impl.SetStaticUsageColumn(col) }

// StaticUsageColumn returns the description column for static flags.
func (f *FlagSet) StaticUsageColumn() int { return f.impl.StaticUsageColumn() }

// SetStaticUsageWidth sets the max wrapping width for static flag descriptions.
func (f *FlagSet) SetStaticUsageWidth(max int) { f.impl.SetStaticUsageWidth(max) }

// StaticUsageWidth returns the wrapping width for static flag descriptions.
func (f *FlagSet) StaticUsageWidth() int { return f.impl.StaticUsageWidth() }

// StaticAutoUsageColumn computes a good usage column for static flags.
func (f *FlagSet) StaticAutoUsageColumn(padding int) int {
	return f.impl.StaticAutoUsageColumn(padding)
}

// SetStaticUsageNote adds a note after the static flag block.
func (f *FlagSet) SetStaticUsageNote(s string) { f.impl.SetStaticUsageNote(s) }

// StaticUsageNote returns the static flag section note.
func (f *FlagSet) StaticUsageNote() string { return f.impl.StaticUsageNote() }

// SetDynamicUsageIndent sets the indentation for dynamic flag usage lines.
func (f *FlagSet) SetDynamicUsageIndent(n int) { f.impl.SetDynamicUsageIndent(n) }

// DynamicUsageIndent returns the dynamic flag usage indent.
func (f *FlagSet) DynamicUsageIndent() int { return f.impl.DynamicUsageIndent() }

// SetDynamicUsageColumn sets the column at which dynamic flag descriptions begin.
func (f *FlagSet) SetDynamicUsageColumn(col int) { f.impl.SetDynamicUsageColumn(col) }

// DynamicUsageColumn returns the description column for dynamic flags.
func (f *FlagSet) DynamicUsageColumn() int { return f.impl.DynamicUsageColumn() }

// SetDynamicUsageWidth sets the max wrapping width for dynamic flags.
func (f *FlagSet) SetDynamicUsageWidth(max int) { f.impl.SetDynamicUsageWidth(max) }

// DynamicUsageWidth returns the wrapping width for dynamic flag descriptions.
func (f *FlagSet) DynamicUsageWidth() int { return f.impl.DynamicUsageWidth() }

// DynamicAutoUsageColumn computes a good usage column for dynamic flags.
func (f *FlagSet) DynamicAutoUsageColumn(padding int) int {
	return f.impl.DynamicAutoUsageColumn(padding)
}

// SetDynamicUsageNote adds a note after the dynamic flag block.
func (f *FlagSet) SetDynamicUsageNote(s string) { f.impl.SetDynamicUsageNote(s) }

// DynamicUsageNote returns the dynamic flag section note.
func (f *FlagSet) DynamicUsageNote() string { return f.impl.DynamicUsageNote() }

// SetNoteIndent sets the indentation for help notes.
func (f *FlagSet) SetNoteIndent(n int) { f.impl.SetNoteIndent(n) }

// NoteIndent returns the note section indentation.
func (f *FlagSet) NoteIndent() int { return f.impl.NoteIndent() }

// SetNoteWidth sets the wrapping width for help notes.
func (f *FlagSet) SetNoteWidth(max int) { f.impl.SetNoteWidth(max) }

// NoteWidth returns the wrapping width for help notes.
func (f *FlagSet) NoteWidth() int { return f.impl.NoteWidth() }

// GetDynamic retrieves the value for a dynamic flag by group, id, and field name.
func GetDynamic[T any](group *dynamic.Group, id, flag string) (T, error) {
	return dynamic.Get[T](group, id, flag)
}

// MustGetDynamic retrieves the value for a dynamic flag or panics.
func MustGetDynamic[T any](group *dynamic.Group, id, flag string) T {
	return dynamic.MustGet[T](group, id, flag)
}

// GetOrDefaultDynamic returns the value for a dynamic flag or its default.
func GetOrDefaultDynamic[T any](group *dynamic.Group, id, flag string) T {
	return dynamic.GetOrDefault[T](group, id, flag)
}
