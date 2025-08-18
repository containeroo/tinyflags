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
	if err := f.checkRequired(); err != nil {
		return f.handleError(err)
	}
	if err := f.checkOneOfGroups(); err != nil {
		return f.handleError(err)
	}
	if err := f.checkAllOrNone(); err != nil {
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
		if envKey == "" && f.envPrefix != "" {
			envKey = strings.ToUpper(f.envPrefix + "_" + strings.ReplaceAll(fl.Name, "-", "_"))
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

// checkOneOfGroups ensures only one flag per group is set.
func (f *FlagSet) checkOneOfGroups() error {
	for _, g := range f.oneOfGroup {
		var conflicting []string

		// Check individual flags
		for _, fl := range g.Flags {
			if fl.Value.Changed() {
				conflicting = append(conflicting, "--"+fl.Name)
			}
		}

		// Check require-together groups
		for _, grp := range g.RequiredGroups {
			set := 0
			for _, fl := range grp.Flags {
				if fl.Value.Changed() {
					set++
				}
			}
			if set > 0 && set != len(grp.Flags) {
				// Skip partially set group (handled by checkAllOrNone)
				continue
			}
			if set == len(grp.Flags) {
				// Show group as comma-separated flag list
				var names []string
				for _, fl := range grp.Flags {
					names = append(names, "--"+fl.Name)
				}
				conflicting = append(conflicting, fmt.Sprintf("[%s]", strings.Join(names, ", ")))
			}
		}

		if len(conflicting) > 1 {
			return fmt.Errorf(
				"only one of the flags in group %q may be used: %s",
				g.Name, strings.Join(conflicting, " vs "),
			)
		}
	}
	return nil
}

// checkAllOrNone ensures all required flags were set.
func (f *FlagSet) checkAllOrNone() error {
	for _, g := range f.allOrNoneGroup {
		set := 0
		for _, fl := range g.Flags {
			if fl.Value.Changed() {
				set++
			}
		}
		if set > 0 && set != len(g.Flags) {
			names := make([]string, 0, len(g.Flags))
			for _, fl := range g.Flags {
				names = append(names, "--"+fl.Name)
			}
			return fmt.Errorf("flags %s must be set together", strings.Join(names, ", "))
		}
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
