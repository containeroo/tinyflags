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
		if meta := getPlaceholder(fl); meta != "" {
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
func printFlagUsage(w io.Writer, indent, startCol, maxWidth int, globalHideEnvs bool, flag *core.BaseFlag, prefix string) {
	var b strings.Builder
	formatStaticFlagNames(&b, flag)

	if meta := getPlaceholder(flag); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}
	flagLine := b.String()
	desc := buildFlagDescription(flag, globalHideEnvs, prefix)

	descWidth := max(maxWidth-indent-startCol-1, 100)
	wrapped := wrapText(desc, descWidth)
	lines := strings.Split(wrapped, "\n")

	fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), startCol, flagLine, lines[0]) // nolint:errcheck

	padding := strings.Repeat(" ", indent+startCol+1)
	for _, line := range lines[1:] {
		fmt.Fprintf(w, "%s%s\n", padding, line) // nolint:errcheck
	}
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

	if meta := getPlaceholder(fl); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}

	return b.String()
}

// getPlaceholder returns the appropriate help placeholder string.
func getPlaceholder(flag *core.BaseFlag) string {
	isBool := false
	isStrict := false

	if bv, ok := flag.Value.(core.StrictBool); ok {
		isBool = true
		isStrict = bv.IsStrictBool()
	}

	if isBool && !isStrict {
		return ""
	}
	if flag.Placeholder != "" {
		return flag.Placeholder
	}
	if len(flag.Allowed) > 0 {
		return "<" + strings.Join(flag.Allowed, "|") + ">"
	}
	if isStrict {
		return "<true|false>"
	}
	if _, ok := flag.Value.(core.Incrementable); ok {
		return ""
	}
	placeholder := strings.ToUpper(flag.Name)
	if _, ok := flag.Value.(core.SliceMarker); ok {
		placeholder += "..."
	}
	return placeholder
}

// printUsageToken prints short, long, or combined flag usage.
func printUsageToken(w io.Writer, fl *core.BaseFlag, mode FlagPrintMode) {
	meta := getPlaceholder(fl)
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

// writeIndented writes text to w, indenting each line by indent spaces.
func writeIndented(w io.Writer, text string, indent, maxWidth int) {
	if text == "" {
		return
	}

	prefix := strings.Repeat(" ", indent)
	wrapped := wrapText(text, maxWidth-indent)

	for _, line := range strings.Split(wrapped, "\n") {
		if line == "" {
			fmt.Fprintln(w) // nolint:errcheck
			continue
		}
		fmt.Fprintf(w, "%s%s\n", prefix, line) // nolint:errcheck
	}
}
