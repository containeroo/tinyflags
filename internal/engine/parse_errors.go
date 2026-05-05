package engine

import (
	"fmt"
	"os"
)

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
