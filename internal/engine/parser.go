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
	if err := f.checkMutualExclusion(); err != nil {
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
	for _, fl := range f.flags {
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

// checkMutualExclusion ensures only one flag per group is set.
func (f *FlagSet) checkMutualExclusion() error {
	for _, g := range f.groups {
		count := 0
		for _, fl := range g.Flags {
			if fl.Value.Changed() {
				count++
			}
		}
		if count > 1 {
			return fmt.Errorf("mutually exclusive flags used in group %q", g.Name)
		}
	}
	return nil
}

// checkRequired ensures all required flags were set.
func (f *FlagSet) checkRequired() error {
	for _, fl := range f.flags {
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
