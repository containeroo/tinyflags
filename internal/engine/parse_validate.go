package engine

import "github.com/containeroo/tinyflags/internal/validate"

// checkRequirements ensures all flag dependencies are satisfied.
func (f *FlagSet) checkRequirements() error {
	return validate.CheckRequirements(f.staticFlagsMap)
}

// checkPositionals ensures all positional arguments are valid.
func (f *FlagSet) checkPositionals() error {
	if err := f.validatePositionals(); err != nil {
		return f.handleError(err)
	}
	if err := f.finalizePositionals(); err != nil {
		return f.handleError(err)
	}
	return nil
}

// checkOneOfGroups ensures only one flag per group is set.
func (f *FlagSet) checkOneOfGroups() error {
	if f.showHelp != nil && *f.showHelp {
		return nil
	}
	if f.showVersion != nil && *f.showVersion {
		return nil
	}

	return validate.CheckOneOfGroups(f.oneOfGroup, f.oneOfVerbose)
}

// checkAllOrNone ensures all required flags were set.
func (f *FlagSet) checkAllOrNone() error {
	if f.showHelp != nil && *f.showHelp {
		return nil
	}
	if f.showVersion != nil && *f.showVersion {
		return nil
	}

	return validate.CheckAllOrNone(f.allOrNoneGroup)
}

// validatePositionals ensures all positional arguments are valid.
func (f *FlagSet) validatePositionals() error {
	return validate.ValidatePositionals(f.positional, f.validatePositional)
}

// finalizePositionals mutates positional arguments.
func (f *FlagSet) finalizePositionals() error {
	return validate.FinalizePositionals(f.positional, f.finalizePositional)
}

// checkRequired ensures all required flags were set.
func (f *FlagSet) checkRequired() error {
	return validate.CheckRequired(f.staticFlagsMap)
}

// checkRequiredDynamic ensures all required dynamic flags are set for each
// existing instance of every dynamic group. Errors use --group.id.flag format.
func (f *FlagSet) checkRequiredDynamic() error {
	return validate.CheckRequiredDynamic(f.dynamicGroups())
}
