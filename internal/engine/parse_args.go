package engine

import (
	"fmt"

	"github.com/containeroo/tinyflags/internal/utils"
)

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
