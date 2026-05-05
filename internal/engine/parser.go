package engine

import (
	"strings"
)

// Parse parses CLI arguments, env vars, built-in help/version, and validations.
func (f *FlagSet) Parse(args []string) error {
	f.maybeAddBuiltinFlags()
	f.resetParseState()

	if f.beforeParse != nil {
		var err error
		args, err = f.beforeParse(args)
		if err != nil {
			return f.handleError(err)
		}
	}

	if err := f.parseArgs(args); err != nil {
		return f.handleError(err)
	}

	// Check if help was requested
	if f.enableHelp && f.showHelp != nil && *f.showHelp {
		var buf strings.Builder
		f.SetOutput(&buf)
		f.Usage()
		return &HelpRequested{Message: buf.String()}
	}

	// Check if version was requested
	if f.enableVer && f.showVersion != nil && *f.showVersion {
		return &VersionRequested{Version: f.versionString}
	}

	// Load values from env and validate
	if err := f.parseEnv(); err != nil {
		return f.handleError(err)
	}
	f.applyDefaultFinalizers()
	if err := f.checkRequired(); err != nil { // static
		return f.handleError(err)
	}
	if err := f.checkRequiredDynamic(); err != nil { // NEW
		return f.handleError(err)
	}
	if err := f.checkOneOfGroups(); err != nil {
		return f.handleError(err)
	}
	if err := f.checkAllOrNone(); err != nil {
		return f.handleError(err)
	}
	if err := f.checkRequirements(); err != nil {
		return f.handleError(err)
	}
	if err := f.checkPositionals(); err != nil {
		return f.handleError(err)
	}
	return nil
}
