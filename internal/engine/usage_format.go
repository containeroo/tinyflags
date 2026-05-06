package engine

import (
	"fmt"
	"io"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
)

// calcStaticUsageColumn calculates the maximum length of a flag line.
func (f *FlagSet) calcStaticUsageColumn(padding int) int {
	maxFlagLen := 0
	for _, fl := range f.staticFlags() {
		var b strings.Builder
		formatStaticFlagNames(&b, fl)
		if meta := fl.UsagePlaceholder(); meta != "" {
			b.WriteString(" ")
			b.WriteString(meta)
		}
		line := b.String()
		if l := len(line); l > maxFlagLen {
			maxFlagLen = l
		}
	}
	return maxFlagLen + padding
}

// calcDynamicUsageColumn calculates the maximum dynamic flag line length.
func (f *FlagSet) calcDynamicUsageColumn(padding int) int {
	maxLen := 0
	for _, group := range f.dynamicGroups() {
		idPlaceholder := group.GetPlaceholder()
		if idPlaceholder == "" {
			idPlaceholder = "<ID>"
		}
		for _, fl := range group.DynamicFlags() {
			line := formatDynamicFlagLine(group.Name(), idPlaceholder, fl)
			if len(line) > maxLen {
				maxLen = len(line)
			}
		}
	}
	return maxLen + padding
}

// printFlagUsage renders a single usage line with wrapping and alignment.
func printFlagUsage(w io.Writer, layout usageLayout, globalHideEnvs bool, flag *core.BaseFlag, prefix string) {
	var b strings.Builder
	formatStaticFlagNames(&b, flag)

	if meta := flag.UsagePlaceholder(); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}
	flagLine := b.String()
	desc := buildFlagDescription(flag, globalHideEnvs, prefix)
	layout.writeWrappedRow(w, flagLine, desc)
}

// formatStaticFlagNames builds the flag name string for help output.
func formatStaticFlagNames(b *strings.Builder, flag *core.BaseFlag) {
	if flag.Short != "" {
		b.WriteString("-")
		b.WriteString(flag.Short)
		b.WriteString(", ")
	} else {
		b.WriteString("    ")
	}
	b.WriteString("--")
	b.WriteString(flag.Name)
}

// formatDynamicFlagLine builds the full flag line string for a dynamic flag.
func formatDynamicFlagLine(groupName, idPlaceholder string, fl *core.BaseFlag) string {
	var b strings.Builder

	b.WriteString("--")
	b.WriteString(groupName)
	b.WriteString(".")
	b.WriteString(idPlaceholder)
	b.WriteString(".")
	b.WriteString(fl.Name)

	if meta := fl.UsagePlaceholder(); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}

	return b.String()
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
