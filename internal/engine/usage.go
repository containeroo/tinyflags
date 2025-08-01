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
		writeIndented(w, wrapText(f.desc, maxWidth-indent), indent)
	}
}

// PrintNotes renders notes block below flags.
func (f *FlagSet) PrintNotes(w io.Writer, indent, maxWidth int) {
	if f.notes != "" {
		writeIndented(w, wrapText(f.notes, maxWidth-indent), indent)
	}
}

// PrintStaticDefaults renders all statically registered flags.
func (f *FlagSet) PrintStaticDefaults(w io.Writer, indent, startCol, maxWidth int) {
	for _, fl := range f.staticFlags() {
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
			writeIndented(w, wrapText(desc, maxWidth-indent), 0)
		}

		idPlaceholder := group.GetPlaceholder()
		if idPlaceholder == "" {
			idPlaceholder = "<ID>"
		}

		for _, fl := range group.DynamicFlags() {
			flagLine := formatDynamicFlagLine(name, idPlaceholder, fl)
			desc := buildFlagDescription(fl, f.hideEnvs, f.envPrefix)
			if len(flagLine)+len(desc) <= maxWidth-startCol {
				fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), startCol, flagLine, desc) // nolint:errcheck
				continue
			}

			wrapped := wrapText(desc, maxWidth-indent-startCol-1)
			lines := strings.Split(wrapped, "\n")
			fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), startCol, flagLine, lines[0]) // nolint:errcheck
			for _, l := range lines[1:] {
				fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), startCol, "", l) // nolint:errcheck
			}
		}

		if note := group.NoteText(); note != "" {
			writeIndented(w, wrapText(note, maxWidth-indent), 0)
		}
	}

	if f.DynamicUsageNote() != "" {
		fmt.Println(f.DynamicUsageNote())
	}
}

// PrintGroups renders usage help for mutual and require-together groups.
func (f *FlagSet) PrintGroups(w io.Writer, indent, maxWidth int) {
	// Required-together groups
	for _, g := range f.requiredTogether {
		if g.IsHidden() {
			continue
		}
		if g.TitleText() != "" {
			fmt.Fprintf(w, "\n%s\n", g.TitleText()) // nolint:errcheck
		}
		if g.IsRequired() {
			fmt.Fprintf(w, "%s(Required)\n", strings.Repeat(" ", indent)) // nolint:errcheck
		}
		if len(g.Flags) > 0 {
			for _, fl := range g.Flags {
				desc := buildFlagDescription(fl, f.hideEnvs, f.envPrefix)
				printGroupFlagUsage(w, indent, fl, desc, maxWidth)
			}
		}
	}

	// Mutual exclusive groups
	for _, g := range f.mutualGroups {
		if g.IsHidden() {
			continue
		}
		if g.TitleText() != "" {
			fmt.Fprintf(w, "\n%s\n", g.TitleText()) // nolint:errcheck
		} else {
			fmt.Fprintf(w, "\nMutually Exclusive Group: %s\n", g.Name) // nolint:errcheck
		}
		if g.IsRequired() {
			fmt.Fprintf(w, "%s(Exactly one required)\n", strings.Repeat(" ", indent)) // nolint:errcheck
		}
		for _, fl := range g.Flags {
			desc := buildFlagDescription(fl, f.hideEnvs, f.envPrefix)
			printGroupFlagUsage(w, indent, fl, desc, maxWidth)
		}
		for _, grp := range g.RequiredGroups {
			names := make([]string, 0, len(grp.Flags))
			for _, fl := range grp.Flags {
				names = append(names, "--"+fl.Name)
			}
			groupLabel := "[" + strings.Join(names, ", ") + "]"
			fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), 20, groupLabel, "(Required together)") // nolint:errcheck
		}
	}
}

// printGroupFlagUsage prints one flag in a group with indent and wrapping.
func printGroupFlagUsage(w io.Writer, indent int, flag *core.BaseFlag, desc string, maxWidth int) {
	var b strings.Builder
	formatStaticFlagNames(&b, flag)
	flagLine := b.String()

	if meta := getPlaceholder(flag); meta != "" {
		flagLine += " " + meta
	}

	if len(flagLine)+len(desc) <= maxWidth-indent {
		fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), 20, flagLine, desc) // nolint:errcheck
		return
	}

	wrapped := wrapText(desc, maxWidth-indent-21)
	lines := strings.Split(wrapped, "\n")
	fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), 20, flagLine, lines[0]) // nolint:errcheck
	for _, l := range lines[1:] {
		fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), 20, "", l) // nolint:errcheck
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

	if len(flagLine)+len(desc) <= maxWidth-startCol {
		fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), startCol, flagLine, desc) // nolint:errcheck
		return
	}

	wrapped := wrapText(desc, maxWidth-indent-startCol-1)
	lines := strings.Split(wrapped, "\n")
	fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), startCol, flagLine, lines[0]) // nolint:errcheck
	for _, l := range lines[1:] {
		fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), startCol, "", l) // nolint:errcheck
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
	placeholder := strings.ToUpper(flag.Name)
	if _, ok := flag.Value.(core.SliceMarker); ok {
		placeholder += "..."
	}
	return placeholder
}

// buildFlagDescription creates the descriptive string for a flag.
func buildFlagDescription(flag *core.BaseFlag, globalHideEnvs bool, name string) string {
	desc := flag.Usage
	if len(flag.Allowed) > 0 {
		desc += " (Allowed: " + strings.Join(flag.Allowed, ", ") + ")"
	}
	if bv, ok := flag.Value.(core.StrictBool); ok && bv.IsStrictBool() {
		desc += " (Allowed: true, false)"
	}
	if flag.Deprecated != "" {
		desc += " [DEPRECATED: " + flag.Deprecated + "]"
	}
	if flag.Value != nil {
		if def := flag.Value.Default(); def != "" {
			desc += " (Default: " + def + ")"
		}
	}
	if shouldInjectEnvKey(flag, globalHideEnvs, name) {
		flag.EnvKey = strings.ToUpper(name + "_" + strings.ReplaceAll(flag.Name, "-", "_"))
	}
	if shouldShowEnv(flag, globalHideEnvs) {
		desc += " (Env: " + flag.EnvKey + ")"
	}
	if !flag.HideRequired && flag.Required {
		desc += " (Required)"
	}
	if flag.MutualGroup != nil && !flag.MutualGroup.IsHidden() {
		desc += buildGroupInfo(flag.MutualGroup)
	}
	if flag.RequiredTogether != nil && !flag.RequiredTogether.IsHidden() {
		desc += buildRequireGroupInfo(flag.RequiredTogether)
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
		b := new(strings.Builder)
		formatStaticFlagNames(b, fl)
		if meta := getPlaceholder(fl); meta != "" {
			b.WriteString(" ")
			b.WriteString(meta)
		}
		line := b.String()
		if len(line) > maxFlagLen {
			maxFlagLen = len(line)
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

// buildGroupInfo returns group info suffix if flag belongs to a mutual group.
func buildGroupInfo(group *core.MutualExlusiveGroup) string {
	var b strings.Builder
	b.WriteString(" (Group: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	if group.IsRequired() {
		b.WriteString(", required")
	}
	b.WriteString(")")
	return b.String()
}

// buildRequireGroupInfo returns group info suffix if flag belongs to a require-together group.
func buildRequireGroupInfo(group *core.RequiredTogetherGroup) string {
	var b strings.Builder
	b.WriteString(" (Require Together: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	if group.IsRequired() {
		b.WriteString(", required")
	}
	b.WriteString(")")
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

// writeIndented prints each line with the given indentation.
func writeIndented(w io.Writer, text string, indent int) {
	lines := strings.Split(text, "\n")
	prefix := strings.Repeat(" ", indent)
	for _, line := range lines {
		fmt.Fprintf(w, "%s%s\n", prefix, line) // nolint:errcheck
	}
}
