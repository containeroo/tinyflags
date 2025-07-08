package tinyflags

import "errors"

type ErrorHandling int

const (
	// ContinueOnError will return an err from Parse() if an error is found
	ContinueOnError ErrorHandling = iota
	// ExitOnError will call os.Exit(2) if an error is found when parsing
	ExitOnError
	// PanicOnError will panic() if an error is found when parsing flags
	PanicOnError
)

// HelpRequested is returned when the built-in help flag (-h or --help) is triggered.
type HelpRequested struct {
	Message string // The full help message to show the user
}

// Error returns the help message, satisfying the error interface.
func (e *HelpRequested) Error() string {
	return e.Message
}

// VersionRequested is returned when the built-in version flag (--version) is triggered.
type VersionRequested struct {
	Version string // The version string to show the user
}

// Error returns the version string, satisfying the error interface.
func (e *VersionRequested) Error() string {
	return e.Version
}

// IsHelpRequested checks if the error is a HelpRequested sentinel.
func IsHelpRequested(err error) bool {
	var helpErr *HelpRequested
	return errors.As(err, &helpErr)
}

// IsVersionRequested checks if the error is a VersionRequested sentinel
func IsVersionRequested(err error) bool {
	var versionErr *VersionRequested
	return errors.As(err, &versionErr)
}

// RequestHelp returns an error with type HelpRequested and the given message.
func RequestHelp(msg string) error {
	return &HelpRequested{Message: msg}
}

// RequestVersion returns an error with type VersionRequested and the given message.
func RequestVersion(msg string) error {
	return &VersionRequested{Version: msg}
}
