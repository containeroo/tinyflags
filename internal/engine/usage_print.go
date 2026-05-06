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
		newUsageLayout(indent, 0, maxWidth).writeIndented(w, f.desc)
	}
}

// PrintNotes renders notes block below flags.
func (f *FlagSet) PrintNotes(w io.Writer, indent, maxWidth int) {
	if f.notes != "" {
		newUsageLayout(indent, 0, maxWidth).writeIndented(w, f.notes)
	}
}

// PrintStaticDefaults renders all statically registered flags.
func (f *FlagSet) PrintStaticDefaults(w io.Writer, indent, startCol, maxWidth int) {
	layout := newUsageLayout(indent, startCol, maxWidth)
	var lastSection string
	for _, fl := range f.staticFlags() {
		if fl.Hidden {
			continue
		}
		if fl.Section != "" && fl.Section != lastSection {
			fmt.Fprintf(w, "\n%s:\n", fl.Section) // nolint:errcheck
			lastSection = fl.Section
		}
		printFlagUsage(w, layout, f.hideEnvs, fl, f.envPrefix)
	}

	if f.StaticUsageNote() != "" {
		fmt.Fprintln(w, f.StaticUsageNote()) // nolint:errcheck
	}
}

// PrintDynamicDefaults renders all dynamic groups.
func (f *FlagSet) PrintDynamicDefaults(w io.Writer, indent, startCol, maxWidth int) {
	layout := newUsageLayout(indent, startCol, maxWidth)
	for _, group := range f.dynamicGroups() {
		if group.IsHidden() {
			continue
		}
		name := group.Name()

		if title := group.TitleText(); title != "" {
			fmt.Fprintf(w, "\n%s\n", title) // nolint:errcheck
		}
		if desc := group.DescriptionText(); desc != "" {
			newUsageLayout(0, 0, maxWidth).writeIndented(w, wrapText(desc, maxWidth-indent))
		}
		idPlaceholder := group.GetPlaceholder()
		if idPlaceholder == "" {
			idPlaceholder = "<ID>"
		}

		for _, fl := range group.DynamicFlags() {
			flagLine := formatDynamicFlagLine(name, idPlaceholder, fl)
			desc := buildFlagDescription(fl, f.hideEnvs, f.envPrefix)

			if len(desc) <= layout.descriptionWidth() {
				_, _ = fmt.Fprintf(w,
					"%s%-*s %s\n",
					strings.Repeat(" ", indent),
					startCol,
					flagLine,
					desc,
				)
				continue
			}

			layout.writeWrappedRow(w, flagLine, desc)
		}

		if note := group.NoteText(); note != "" {
			newUsageLayout(indent, 0, maxWidth).writeIndented(w, wrapText(note, maxWidth-indent))
		}
	}

	if f.DynamicUsageNote() != "" {
		fmt.Fprintln(w, f.DynamicUsageNote()) // nolint:errcheck
	}
}
