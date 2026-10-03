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

	if request := f.builtinRequest(); request != nil {
		return request
	}

	if err := f.parseEnv(); err != nil {
		return f.handleError(err)
	}
	if err := f.validateParsedValues(); err != nil {
		return f.handleError(err)
	}
	return nil
}

// builtinRequest returns the requested built-in help or version response, if any.
func (f *FlagSet) builtinRequest() error {
	if f.enableHelp && f.showHelp != nil && *f.showHelp {
		var buf strings.Builder
		previousOutput := f.Output()
		f.SetOutput(&buf)
		defer f.SetOutput(previousOutput)
		f.Usage()
		return &HelpRequested{Message: buf.String()}
	}
	if f.enableVer && f.showVersion != nil && *f.showVersion {
		return &VersionRequested{Version: f.versionString}
	}
	return nil
}

// validateParsedValues applies finalization and runs every post-parse validation stage.
func (f *FlagSet) validateParsedValues() error {
	f.applyDefaultFinalizers()

	if err := f.checkRequired(); err != nil {
		return err
	}
	if err := f.checkRequiredDynamic(); err != nil {
		return err
	}
	if err := f.checkNotEmpty(); err != nil {
		return err
	}
	if err := f.checkNotEmptyDynamic(); err != nil {
		return err
	}
	if err := f.checkOneOfGroups(); err != nil {
		return err
	}
	if err := f.checkAllOrNone(); err != nil {
		return err
	}
	if err := f.checkRequirements(); err != nil {
		return err
	}
	if err := f.checkPositionals(); err != nil {
		return err
	}
	for _, validate := range f.validators {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}
