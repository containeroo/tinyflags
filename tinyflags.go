// Package tinyflags provides a high-level API for defining and parsing
// CLI flags with support for dynamic groups, custom types, and rich usage output.
package tinyflags

import (
	"errors"

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

// UsageError wraps a semantic parse error with rendered help text.
type UsageError struct {
	Err  error
	Help string
}

// Error returns the wrapped error message.
func (e *UsageError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// Unwrap exposes the semantic parse error.
func (e *UsageError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// CommandRequired is returned when one command requires a subcommand selection.
type CommandRequired struct {
	Command string
}

// Error returns the human-readable message for a missing required subcommand.
func (e *CommandRequired) Error() string {
	return `command "` + e.Command + `" requires a subcommand`
}

var (
	IsHelpRequested    = engine.IsHelpRequested
	IsVersionRequested = engine.IsVersionRequested
	RequestHelp        = engine.RequestHelp
	RequestVersion     = engine.RequestVersion
)

// IsCommandRequired checks whether err indicates a missing required subcommand.
func IsCommandRequired(err error) bool {
	var target *CommandRequired
	return errors.As(err, &target)
}

// HelpText extracts rendered help text from a usage-bearing error.
func HelpText(err error) (string, bool) {
	var usageErr *UsageError
	if !errors.As(err, &usageErr) || usageErr == nil || usageErr.Help == "" {
		return "", false
	}
	return usageErr.Help, true
}

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
	Flag[T any]  = core.Flag[T]  // Minimal flag handle interface
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
