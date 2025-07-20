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

// Parse triggers parsing of args and environment.
func (f *FlagSet) Parse(args []string) error {
	if f.Usage != nil {
		f.impl.Usage = f.Usage
	}
	return f.impl.Parse(args)
}

// Public API passthroughs
func (f *FlagSet) Version(s string)                       { f.impl.Version(s) }
func (f *FlagSet) EnvPrefix(s string)                     { f.impl.EnvPrefix(s) }
func (f *FlagSet) Title(s string)                         { f.impl.Title(s) }
func (f *FlagSet) Description(s string)                   { f.impl.Description(s) }
func (f *FlagSet) Note(s string)                          { f.impl.Note(s) }
func (f *FlagSet) DisableHelp()                           { f.impl.DisableHelp() }
func (f *FlagSet) DisableVersion()                        { f.impl.DisableVersion() }
func (f *FlagSet) Sorted(b bool)                          { f.impl.Sorted(b) }
func (f *FlagSet) SetOutput(w io.Writer)                  { f.impl.SetOutput(w) }
func (f *FlagSet) Output() io.Writer                      { return f.impl.Output() }
func (f *FlagSet) IgnoreInvalidEnv(b bool)                { f.impl.IgnoreInvalidEnv(b) }
func (f *FlagSet) SetGetEnvFn(fn func(string) string)     { f.impl.SetGetEnvFn(fn) }
func (f *FlagSet) Globaldelimiter(s string)               { f.impl.Globaldelimiter(s) }
func (f *FlagSet) GetGroup(name string) *core.MutualGroup { return f.impl.GetGroup(name) }
func (f *FlagSet) RequirePositional(n int)                { f.impl.RequirePositional(n) }
func (f *FlagSet) Args() []string                         { return f.impl.Args() }
func (f *FlagSet) Arg(i int) (string, bool)               { return f.impl.Arg(i) }
func (f *FlagSet) DescriptionMaxLen(n int)                { f.impl.DescriptionMaxLen(n) }
func (f *FlagSet) DescriptionIndent(n int)                { f.impl.DescriptionIndent(n) }
func (f *FlagSet) PrintDefaults()                         { f.impl.PrintDefaults() }
func (f *FlagSet) PrintUsage(w io.Writer, mode FlagPrintMode) {
	f.impl.PrintUsage(w, mode)
}
func (f *FlagSet) PrintTitle(w io.Writer)                  { f.impl.PrintTitle(w) }
func (f *FlagSet) PrintNotes(w io.Writer, width int)       { f.impl.PrintNotes(w, width) }
func (f *FlagSet) PrintDescription(w io.Writer, width int) { f.impl.PrintDescription(w, width) }

// Dynamic groups
func (f *FlagSet) DynamicGroup(name string) *dynamic.Group {
	return f.impl.DynamicGroup(name)
}

func (f *FlagSet) DefaultDelimiter() string {
	return f.impl.DefaultDelimiter()
}

func (f *FlagSet) RegisterDynamic(group, field string, val core.DynamicValue) error {
	return f.impl.RegisterDynamic(group, field, val)
}

func (f *FlagSet) RegisterFlag(name string, bf *core.BaseFlag) { f.impl.RegisterFlag(name, bf) }

// Mutual group passthrough (if needed)
func (f *FlagSet) Groups() []*core.MutualGroup {
	return f.impl.Groups()
}
