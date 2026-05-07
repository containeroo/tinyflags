package engine

import (
	"fmt"
	"io"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/help"
)

// calcStaticUsageColumn calculates the maximum length of a flag line.
func (f *FlagSet) calcStaticUsageColumn(padding int) int {
	return help.CalcStaticUsageColumn(f.staticFlags(), padding)
}

// calcDynamicUsageColumn calculates the maximum dynamic flag line length.
func (f *FlagSet) calcDynamicUsageColumn(padding int) int {
	return help.CalcDynamicUsageColumn(f.dynamicGroups(), padding)
}

// printUsageToken prints short, long, or combined flag usage.
func printUsageToken(w io.Writer, fl *core.BaseFlag, mode FlagPrintMode) {
	meta := fl.UsagePlaceholder()
	switch mode {
	case PrintShort:
		if fl.Short != "" {
			fmt.Fprintf(w, " -%s", fl.Short) // nolint:errcheck
			if meta != "" {
				fmt.Fprintf(w, " %s", meta) // nolint:errcheck
			}
		}
	case PrintLong:
		fmt.Fprintf(w, " --%s", fl.Name) // nolint:errcheck
		if meta != "" {
			fmt.Fprintf(w, " %s", meta) // nolint:errcheck
		}
	case PrintBoth:
		if fl.Short != "" {
			fmt.Fprintf(w, " -%s|--%s", fl.Short, fl.Name) // nolint:errcheck
		} else {
			fmt.Fprintf(w, " --%s", fl.Name) // nolint:errcheck
		}
		if meta != "" {
			fmt.Fprintf(w, " %s", meta) // nolint:errcheck
		}
	}
}
