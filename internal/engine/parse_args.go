package engine

import (
	"errors"
	"fmt"

	"github.com/containeroo/tinyflags/internal/utils"
)

// parseArgs processes CLI arguments and sets flags or positional args.
func (f *FlagSet) parseArgs(args []string) error {
	positional, err := runArgParserFSM(f, args)
	if err != nil {
		return err
	}
	f.positional = append(f.positional, positional...)

	if f.showHelp != nil && *f.showHelp {
		return nil
	}
	if f.showVersion != nil && *f.showVersion {
		return nil
	}

	return f.validatePositionalCount()
}

// validatePositionalCount enforces configured positional argument count constraints.
func (f *FlagSet) validatePositionalCount() error {
	count := len(f.positional)

	if f.exactPositional && f.minPositional == 0 && count > 0 {
		return errors.New("accepts no positional arguments")
	}

	if f.exactPositional && count != f.minPositional {
		return fmt.Errorf(
			"expected exactly %d positional argument%s, got %d",
			f.minPositional,
			utils.PluralSuffix(f.minPositional),
			count,
		)
	}

	if f.minPositional > 0 && count < f.minPositional {
		return fmt.Errorf(
			"expected at least %d positional argument%s, got %d",
			f.minPositional,
			utils.PluralSuffix(f.minPositional),
			count,
		)
	}

	if f.maxPositionalSet && count > f.maxPositional {
		return fmt.Errorf(
			"expected at most %d positional argument%s, got %d",
			f.maxPositional,
			utils.PluralSuffix(f.maxPositional),
			count,
		)
	}

	return nil
}
