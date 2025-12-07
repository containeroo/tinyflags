package engine

import (
	"fmt"
	"io"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
)

// FlagPrintMode defines how usage should be rendered.
type FlagPrintMode int

const (
	PrintShort FlagPrintMode = iota // e.g. -v
	PrintLong                       // e.g. --verbose
	PrintBoth                       // e.g. -v|--verbose
	PrintFlags                      // prints only [flags]
	PrintNone                       // prints nothing
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
		fmt.Println(f.StaticUsageNote())
	}
}

// PrintDynamicDefaults renders all dynamic groups.
func (f *FlagSet) PrintDynamicDefaults(w io.Writer, indent, startCol, maxWidth int) {
	for _, group := range f.dynamicGroups() {
		if group.IsHidden() {
			continue
		}
		name := group.Name()

		// Title
		if title := group.TitleText(); title != "" {
			fmt.Fprintf(w, "\n%s\n", title) // nolint:errcheck
		}
		// Description
		if desc := group.DescriptionText(); desc != "" {
			writeIndented(w, wrapText(desc, maxWidth-indent), 0, maxWidth)
		}
		idPlaceholder := group.GetPlaceholder()
		if idPlaceholder == "" {
			idPlaceholder = "<ID>"
		}

		// Match static behavior: give the description column a generous minimum
		descWidth := max(maxWidth-indent-startCol-1, 100)

		for _, fl := range group.DynamicFlags() {
			flagLine := formatDynamicFlagLine(name, idPlaceholder, fl)
			desc := buildFlagDescription(fl, f.hideEnvs, f.envPrefix)

			// One-liner if it fits
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

			// Otherwise wrap only the description and align under the desc column
			wrapped := wrapText(desc, descWidth)
			lines := strings.Split(wrapped, "\n")

			// First line with flag label + first chunk of desc
			_, _ = fmt.Fprintf(w,
				"%s%-*s %s\n",
				strings.Repeat(" ", indent),
				startCol,
				flagLine,
				lines[0],
			)

			// Continuation lines under the description column
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
		fmt.Println(f.DynamicUsageNote())
	}
}

// splitFlags separates user-defined and built-in flags.
func (f *FlagSet) splitFlags() (userFlags []*core.BaseFlag, versionFlag, helpFlag *core.BaseFlag) {
	for _, fl := range f.staticFlags() {
		switch {
		case fl.Name == "version" && f.enableVer:
			versionFlag = fl
		case fl.Name == "help" && f.enableHelp:
			helpFlag = fl
		default:
			userFlags = append(userFlags, fl)
		}
	}
	return
}

// orderedStaticFlags returns all static flags in desired order.
func (f *FlagSet) dynamicGroups() []*dynamic.Group {
	if f.sortGroups {
		return f.OrderedDynamicGroups()
	}
	return f.dynamicGroupsOrder
}

// orderedStaticFlags returns all static flags in desired order.
func (f *FlagSet) staticFlags() []*core.BaseFlag {
	if f.sortFlags {
		return f.OrderedStaticFlags()
	}
	return f.staticFlagsOrder
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

	// Width for description (after flag column)
	descWidth := max(maxWidth-indent-startCol-1, 100)

	// Wrap only the description
	wrapped := wrapText(desc, descWidth)
	lines := strings.Split(wrapped, "\n")

	// First line: print flag line and first part of desc
	fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), startCol, flagLine, lines[0]) // nolint:errcheck

	// Remaining lines: align under desc column
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

// buildFlagDescription creates the full help text for a flag, including metadata
// such as allowed values, default, environment variable, deprecation, and group info.
func buildFlagDescription(flag *core.BaseFlag, globalHideEnvs bool, name string) string {
	desc := flag.Usage

	// Determine allowed values:
	// If explicitly set via flag.Allowed, use those.
	// Otherwise, if it's a strict bool, default to "true,false".
	var allowed []string
	if len(flag.Allowed) > 0 {
		allowed = append(allowed, flag.Allowed...)
	} else if bv, ok := flag.Value.(core.StrictBool); ok && bv.IsStrictBool() {
		allowed = append(allowed, "true", "false")
	}
	if !flag.HideAllowed && len(allowed) > 0 {
		desc += " (Allowed: " + strings.Join(allowed, ", ") + ")"
	}

	// Determine whether to print the default value:
	// For non-strict bools and counters, suppress the default entirely.
	showDefault := true
	if bv, ok := flag.Value.(core.StrictBool); ok && !bv.IsStrictBool() {
		showDefault = false
	}
	if _, ok := flag.Value.(core.Incrementable); ok {
		showDefault = false
	}

	// Append deprecation notice, if applicable.
	if flag.Deprecated != "" {
		desc += " [DEPRECATED: " + flag.Deprecated + "]"
	}

	// Append default value, if available and allowed.
	if flag.Value != nil && showDefault && !flag.HideDefault {
		if def := flag.Value.Default(); def != "" {
			desc += " (Default: " + def + ")"
		}
	}

	// Append required marker, if not hidden.
	if !flag.HideRequires && len(flag.Requires) > 0 {
		desc += " (Requires: " + strings.Join(flag.Requires, ", ") + ")"
	}

	// If no EnvKey is explicitly set and it's allowed, generate one from prefix and flag name.
	if shouldInjectEnvKey(flag, globalHideEnvs, name) {
		flag.EnvKey = strings.ToUpper(name + "_" + strings.ReplaceAll(flag.Name, "-", "_"))
	}

	// Append environment variable name, if applicable.
	if shouldShowEnv(flag, globalHideEnvs) {
		desc += " (Env: " + flag.EnvKey + ")"
	}

	// Append required marker, if not hidden.
	if !flag.HideRequired && flag.Required {
		desc += " (Required)"
	}

	// Append group info if part of a One Of group.
	if flag.OneOfGroup != nil && !flag.OneOfGroup.IsHidden() {
		desc += buildGroupInfo(flag.OneOfGroup)
	}

	// Append group info if part of a Require Together group.
	if flag.AllOrNone != nil && !flag.AllOrNone.IsHidden() {
		desc += buildRequireGroupInfo(flag.AllOrNone)
	}

	return desc
}

// shouldInjectEnvKey decides whether to compute EnvKey from prefix.
func shouldInjectEnvKey(flag *core.BaseFlag, globalHideEnvs bool, prefix string) bool {
	return !globalHideEnvs && !flag.DisableEnv && !flag.HideEnv && flag.EnvKey == "" && prefix != ""
}

// shouldShowEnv decides whether to include EnvKey in help.
func shouldShowEnv(flag *core.BaseFlag, globalHideEnvs bool) bool {
	return !globalHideEnvs && !flag.DisableEnv && !flag.HideEnv && flag.EnvKey != ""
}

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

// buildGroupInfo returns group info suffix if flag belongs to a one of group.
func buildGroupInfo(group *core.OneOfGroupGroup) string {
	var b strings.Builder
	b.WriteString(" [Group: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	b.WriteString(" (One Of)")
	if group.IsRequired() {
		b.WriteString(" - required")
	}
	b.WriteString("]")
	return b.String()
}

// buildRequireGroupInfo returns group info suffix if flag belongs to a require-together group.
func buildRequireGroupInfo(group *core.AllOrNoneGroup) string {
	var b strings.Builder
	b.WriteString(" [Group: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	b.WriteString(" (All Or None)")
	if group.IsRequired() {
		b.WriteString(" - required")
	}
	b.WriteString("]")
	return b.String()
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
