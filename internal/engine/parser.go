package engine

import (
	"fmt"
	"os"
	"strings"

	"github.com/containeroo/tinyflags/internal/utils"
)

// Parse parses CLI arguments, env vars, built-in help/version, and validations.
func (f *FlagSet) Parse(args []string) error {
	f.maybeAddBuiltinFlags()

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

// parseArgs processes CLI arguments and sets flags or positional args.
func (f *FlagSet) parseArgs(args []string) error {
	positional, err := parseArgsWithFSM(f, args)
	if err != nil {
		return err
	}
	f.positional = append(f.positional, positional...)

	if f.requiredPositional > 0 && len(f.positional) < f.requiredPositional {
		if f.showHelp != nil && *f.showHelp {
			return nil
		}
		if f.showVersion != nil && *f.showVersion {
			return nil
		}

		return fmt.Errorf("expected at least %d positional argument%s, got %d",
			f.requiredPositional,
			utils.PluralSuffix(f.requiredPositional),
			len(f.positional))
	}
	return nil
}

// parseEnv loads unset flags from environment variables.
func (f *FlagSet) parseEnv() error {
	for _, fl := range f.staticFlagsMap {
		if fl.Value == nil {
			// dynamically‐registered flags aren’t loaded from ENV
			continue
		}
		if fl.Value.Changed() {
			continue
		}
		envKey := fl.EnvKey
		if envKey == "" {
			envKey = f.envKeyFunc(f.envPrefix, fl.Name)
		}
		val := f.getEnv(envKey)
		if val == "" {
			continue
		}
		if err := fl.Value.Set(val); err != nil {
			if f.ignoreInvalidEnv {
				continue
			}
			return fmt.Errorf("invalid environment value for %s: %w", fl.Name, err)
		}
	}
	return nil
}

// checkRequirements ensures all flags which requrie others are set.
func (f *FlagSet) checkRequirements() error {
	for _, fl := range f.staticFlagsMap {
		// only enforce for flags that are actually set/changed
		if fl.Value == nil || !fl.Value.Changed() {
			continue
		}
		for _, req := range fl.Requires {
			rfl, ok := f.staticFlagsMap[req]
			if !ok || rfl.Value == nil || !rfl.Value.Changed() {
				return fmt.Errorf("--%s requires --%s", fl.Name, req)
			}
		}
	}
	return nil
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
	for _, g := range f.oneOfGroup {
		selections := 0
		for _, fl := range g.Flags {
			if fl.Value.Changed() {
				selections++
			}
		}
		for _, grp := range g.RequiredGroups {
			changed := 0
			for _, fl := range grp.Flags {
				if fl.Value.Changed() {
					changed++
				}
			}
			if changed == len(grp.Flags) && len(grp.Flags) > 0 {
				selections++
			}
		}

		if selections > 1 {
			return fmt.Errorf("only one of the flags in group %q may be used", g.Name)
		}
		if g.IsRequired() && selections == 0 {
			return fmt.Errorf("one of the flags in group %q must be set", g.Name)
		}
	}
	return nil
}

// checkAllOrNone ensures all required flags were set.
func (f *FlagSet) checkAllOrNone() error {
	for _, g := range f.allOrNoneGroup {
		changed := 0
		for _, fl := range g.Flags {
			if fl.Value.Changed() {
				changed++
			}
		}
		if changed > 0 && changed != len(g.Flags) {
			names := make([]string, 0, len(g.Flags))
			for _, fl := range g.Flags {
				names = append(names, "--"+fl.Name)
			}
			return fmt.Errorf("flags %s must be set together", strings.Join(names, ", "))
		}
		if g.IsRequired() && changed == 0 {
			names := make([]string, 0, len(g.Flags))
			for _, fl := range g.Flags {
				names = append(names, "--"+fl.Name)
			}
			return fmt.Errorf("flags %s must be set together", strings.Join(names, ", "))
		}
	}
	return nil
}

// validatePositionals ensures all positional arguments are valid.
func (f *FlagSet) validatePositionals() error {
	if len(f.positional) == 0 {
		return nil
	}
	if f.validatePositional == nil {
		return nil
	}
	for _, arg := range f.positional {
		if err := f.validatePositional(arg); err != nil {
			return err
		}
	}
	return nil
}

// finalizePositionals mutates positional arguments.
func (f *FlagSet) finalizePositionals() error {
	if len(f.positional) == 0 {
		return nil
	}
	if f.finalizePositional == nil {
		return nil
	}
	for i, arg := range f.positional {
		f.positional[i] = f.finalizePositional(arg)
	}
	return nil
}

// checkRequired ensures all required flags were set.
func (f *FlagSet) checkRequired() error {
	for _, fl := range f.staticFlagsMap {
		if fl.Required && !fl.Value.Changed() {
			return fmt.Errorf("flag --%s is required", fl.Name)
		}
	}
	return nil
}

// checkRequiredDynamic ensures all required dynamic flags are set for each
// existing instance of every dynamic group. Errors use --group.id.flag format.
func (f *FlagSet) checkRequiredDynamic() error {
	if len(f.dynamicGroups()) == 0 {
		return nil
	}

	for _, g := range f.dynamicGroups() {
		ids := g.Instances() // []string of all IDs seen (from parsed values)
		if len(ids) == 0 {
			// No instances present → nothing to enforce for this group.
			continue
		}

		items := g.Items() // map[field]core.GroupItem
		if len(items) == 0 {
			continue
		}

		for _, id := range ids {
			for field, item := range items {
				bf := item.Flag
				if bf == nil || !bf.Required {
					continue
				}
				// Was a value set for this id?
				if _, ok := item.Value.GetAny(id); !ok {
					return fmt.Errorf("flag --%s.%s.%s is required", g.Name(), id, field)
				}
			}
		}
	}
	return nil
}

// handleError responds to errors based on the configured mode.
func (f *FlagSet) handleError(err error) error {
	switch f.errorHandling {
	case ContinueOnError:
		return err
	case ExitOnError:
		fmt.Fprintf(f.Output(), "Error: %v\n", err) // nolint:errcheck
		os.Exit(2)
	case PanicOnError:
		panic(err)
	}
	return err
}
