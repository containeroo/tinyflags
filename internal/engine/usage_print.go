package engine

import (
	"fmt"
	"io"
	"strings"
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
		writeIndented(w, f.desc, indent, maxWidth)
	}
}

// PrintNotes renders notes block below flags.
func (f *FlagSet) PrintNotes(w io.Writer, indent, maxWidth int) {
	if f.notes != "" {
		writeIndented(w, f.notes, indent, maxWidth)
	}
}

// PrintStaticDefaults renders all statically registered flags.
func (f *FlagSet) PrintStaticDefaults(w io.Writer, indent, startCol, maxWidth int) {
	var lastSection string
	for _, fl := range f.staticFlags() {
		if fl.Hidden {
			continue
		}
		if fl.Section != "" && fl.Section != lastSection {
			fmt.Fprintf(w, "\n%s:\n", fl.Section) // nolint:errcheck
			lastSection = fl.Section
		}
		printFlagUsage(w, indent, startCol, maxWidth, f.hideEnvs, fl, f.envPrefix)
	}

	if f.StaticUsageNote() != "" {
		fmt.Fprintln(w, f.StaticUsageNote()) // nolint:errcheck
	}
}

// PrintDynamicDefaults renders all dynamic groups.
func (f *FlagSet) PrintDynamicDefaults(w io.Writer, indent, startCol, maxWidth int) {
	for _, group := range f.dynamicGroups() {
		if group.IsHidden() {
			continue
		}
		name := group.Name()

		if title := group.TitleText(); title != "" {
			fmt.Fprintf(w, "\n%s\n", title) // nolint:errcheck
		}
		if desc := group.DescriptionText(); desc != "" {
			writeIndented(w, wrapText(desc, maxWidth-indent), 0, maxWidth)
		}
		idPlaceholder := group.GetPlaceholder()
		if idPlaceholder == "" {
			idPlaceholder = "<ID>"
		}

		descWidth := max(maxWidth-indent-startCol-1, 100)

		for _, fl := range group.DynamicFlags() {
			flagLine := formatDynamicFlagLine(name, idPlaceholder, fl)
			desc := buildFlagDescription(fl, f.hideEnvs, f.envPrefix)

			if len(desc) <= descWidth {
				_, _ = fmt.Fprintf(w,
					"%s%-*s %s\n",
					strings.Repeat(" ", indent),
					startCol,
					flagLine,
					desc,
				)
				continue
			}

			wrapped := wrapText(desc, descWidth)
			lines := strings.Split(wrapped, "\n")

			_, _ = fmt.Fprintf(w,
				"%s%-*s %s\n",
				strings.Repeat(" ", indent),
				startCol,
				flagLine,
				lines[0],
			)

			padding := strings.Repeat(" ", indent+startCol+1)
			for _, l := range lines[1:] {
				fmt.Fprintf(w, "%s%s\n", padding, l) // nolint:errcheck
			}
		}

		if note := group.NoteText(); note != "" {
			writeIndented(w, wrapText(note, maxWidth-indent), indent, maxWidth)
		}
	}

	if f.DynamicUsageNote() != "" {
		fmt.Fprintln(w, f.DynamicUsageNote()) // nolint:errcheck
	}
}
