package engine

import (
	"fmt"
	"io"

	"github.com/containeroo/tinyflags/internal/help"
)

// PrintUsage prints a usage line depending on mode.
func (f *FlagSet) PrintUsage(w io.Writer, mode FlagPrintMode) {
	fmt.Fprint(w, "Usage: "+f.name) // nolint:errcheck

	switch mode {
	case PrintNone:
		fmt.Fprintln(w) // nolint:errcheck
		return
	case PrintFlags:
		fmt.Fprintln(w, " [flags]") // nolint:errcheck
		return
	}

	userFlags, versionFlag, helpFlag := f.splitFlags()

	for _, fl := range userFlags {
		if fl.Hidden {
			continue
		}
		printUsageToken(w, fl, mode)
	}
	if versionFlag != nil {
		printUsageToken(w, versionFlag, mode)
	}
	if helpFlag != nil {
		printUsageToken(w, helpFlag, mode)
	}
	fmt.Fprintln(w) // nolint:errcheck
}

// PrintTitle writes usage title heading.
func (f *FlagSet) PrintTitle(w io.Writer) {
	if f.title != "" {
		fmt.Fprintln(w, f.title) // nolint:errcheck
	}
}

// PrintAuthors writes usage author heading.
func (f *FlagSet) PrintAuthors(w io.Writer) {
	if f.authors != "" {
		fmt.Fprintln(w, "Authors: "+f.authors) // nolint:errcheck
	}
}

// PrintDescription renders description block above flags.
func (f *FlagSet) PrintDescription(w io.Writer, indent, maxWidth int) {
	if f.desc != "" {
		help.WriteIndented(w, f.desc, indent, maxWidth)
	}
}

// PrintNotes renders notes block below flags.
func (f *FlagSet) PrintNotes(w io.Writer, indent, maxWidth int) {
	if f.notes != "" {
		help.WriteIndented(w, f.notes, indent, maxWidth)
	}
}

// PrintStaticDefaults renders all statically registered flags.
func (f *FlagSet) PrintStaticDefaults(w io.Writer, indent, startCol, maxWidth int) {
	help.PrintStaticDefaults(w, f.staticFlags(), indent, startCol, maxWidth, f.hideEnvs, f.envPrefix, f.StaticUsageNote())
}

// PrintDynamicDefaults renders all dynamic groups.
func (f *FlagSet) PrintDynamicDefaults(w io.Writer, indent, startCol, maxWidth int) {
	help.PrintDynamicDefaults(w, f.dynamicGroups(), indent, startCol, maxWidth, f.hideEnvs, f.envPrefix, f.DynamicUsageNote())
}
