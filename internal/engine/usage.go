package engine

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
)

// FlagPrintMode defines how usage should be rendered.
type FlagPrintMode int

const (
	PrintShort FlagPrintMode = iota
	PrintLong
	PrintBoth
	PrintFlags
	PrintNone
)

// PrintUsage prints a usage line depending on mode.
func (f *FlagSet) PrintUsage(w io.Writer, mode FlagPrintMode) {
	fmt.Fprint(w, "Usage: "+f.name)

	switch mode {
	case PrintNone:
		fmt.Fprintln(w)
		return
	case PrintFlags:
		fmt.Fprintln(w, " [flags]")
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
	fmt.Fprintln(w)
}

// PrintTitle writes usage title heading.
func (f *FlagSet) PrintTitle(w io.Writer) {
	if f.title != "" {
		fmt.Fprintln(w, f.title)
	}
}

// PrintAuthors writes usage author heading.
func (f *FlagSet) PrintAuthors(w io.Writer) {
	if f.authors != "" {
		fmt.Fprintln(w, "Authors: "+f.authors)
	}
}

// PrintDescription renders description block above flags.
func (f *FlagSet) PrintDescription(w io.Writer, width int) {
	if f.desc != "" {
		fmt.Fprintln(w, wrapText(f.desc, width))
	}
}

// PrintNotes renders notes block below flags.
func (f *FlagSet) PrintNotes(w io.Writer, width int) {
	if f.notes != "" {
		fmt.Fprintln(w, wrapText(f.notes, width))
	}
}

// PrintDefaults prints both static and dynamic flags.
func (f *FlagSet) PrintDefaults(w io.Writer, width int) {
	f.printStaticDefaults(w, width)
	f.printDynamicDefaults(w, width)
}

// printStaticDefaults renders all statically registered flags.
func (f *FlagSet) printStaticDefaults(w io.Writer, width int) {
	tw := tabwriter.NewWriter(w, 2, 4, 2, ' ', 0)

	staticFlags := f.orderedFlags()

	if f.sortFlags {
		sort.Slice(staticFlags, func(i, j int) bool {
			return staticFlags[i].Name < staticFlags[j].Name
		})
	}

	for _, fl := range staticFlags {
		printFlagUsage(tw, f.descIndent, f.descMaxLen, f.hideEnvs, fl, f.envPrefix)
	}
	tw.Flush() // nolint:errcheck
}

// printDynamicDefaults renders all dynamic groups.
func (f *FlagSet) printDynamicDefaults(w io.Writer, width int) {
	type dynGroup struct {
		name string
		g    *dynamic.Group
	}
	var groups []dynGroup
	for name, g := range f.dynamicGroups {
		if g.IsHidden() {
			continue
		}
		groups = append(groups, dynGroup{name: name, g: g})
	}

	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].g.IsGroupSorted() && groups[j].g.IsGroupSorted() {
			return groups[i].name < groups[j].name
		}
		return false
	})

	for _, group := range groups {
		name := group.name
		g := group.g

		// Title
		if title := g.TitleText(); title != "" {
			fmt.Fprintf(w, "\n%s\n", title)
		}

		// Description
		if desc := g.DescriptionText(); desc != "" {
			fmt.Fprintln(w, wrapText(desc, width))
		}

		// Placeholder for ID
		idPlaceholder := g.GetPlaceholder()
		if idPlaceholder == "" {
			idPlaceholder = "<ID>"
		}

		// Flag definitions (not instances)
		flags := g.Flags()
		if g.IsFlagSorted() {
			sort.Slice(flags, func(i, j int) bool {
				return flags[i].Name < flags[j].Name
			})
		}

		tw := tabwriter.NewWriter(w, 2, 4, 2, ' ', 0)
		for _, fl := range flags {
			// Build flag path: --group.<ID>.flagname
			var b strings.Builder
			b.WriteString("      --")
			b.WriteString(name)
			b.WriteString(".")
			b.WriteString(idPlaceholder)
			b.WriteString(".")
			b.WriteString(fl.Name)
			if meta := getPlaceholder(fl); meta != "" {
				b.WriteString(" ")
				b.WriteString(meta)
			}
			flagLine := b.String()
			desc := buildFlagDescription(fl, f.hideEnvs, f.envPrefix)
			fmt.Fprintf(tw, "%-*s %s\n", f.descIndent, flagLine, desc)
		}
		tw.Flush() // nolint:errcheck

		// Group note
		if note := g.NoteText(); note != "" {
			fmt.Fprintln(w, wrapText(note, width))
		}
	}
}

// splitFlags separates user-defined and built-in flags.
func (f *FlagSet) splitFlags() (userFlags []*core.BaseFlag, versionFlag, helpFlag *core.BaseFlag) {
	for _, fl := range f.orderedFlags() {
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

// orderedFlags returns all static flags in desired order.
func (f *FlagSet) orderedFlags() []*core.BaseFlag {
	if f.sortFlags {
		var all []*core.BaseFlag
		for _, fl := range f.staticFlags {
			all = append(all, fl)
		}
		sort.Slice(all, func(i, j int) bool {
			return all[i].Name < all[j].Name
		})
		return all
	}
	return f.registered
}

// printFlagUsage renders a single usage line.
func printFlagUsage(w io.Writer, descIndent, maxDesc int, globalHideEnvs bool, flag *core.BaseFlag, prefix string) {
	var b strings.Builder

	formatFlagNames(&b, flag)

	if meta := getPlaceholder(flag); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}
	flagLine := b.String()
	desc := buildFlagDescription(flag, globalHideEnvs, prefix)

	if len(flagLine)+len(desc) <= maxDesc {
		fmt.Fprintf(w, "%-*s %s\n", descIndent, flagLine, desc)
		return
	}

	wrapped := wrapText(desc, maxDesc-descIndent-1)
	lines := strings.Split(wrapped, "\n")
	fmt.Fprintf(w, "%-*s %s\n", descIndent, flagLine, lines[0])
	for _, l := range lines[1:] {
		fmt.Fprintf(w, "%-*s %s\n", descIndent, "", l)
	}
}

// formatFlagNames builds flag names (e.g. -v, --verbose).
func formatFlagNames(b *strings.Builder, flag *core.BaseFlag) {
	if flag.Short != "" {
		b.WriteString("  -")
		b.WriteString(flag.Short)
		b.WriteString(", ")
	} else {
		b.WriteString("      ")
	}
	b.WriteString("--")
	b.WriteString(flag.Name)
}

// getPlaceholder returns help placeholder for a flag value.
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

// buildFlagDescription renders the description of a flag.
func buildFlagDescription(flag *core.BaseFlag, globalHideEnvs bool, prefix string) string {
	desc := flag.Usage
	if len(flag.Allowed) > 0 {
		desc += " (Allowed: " + strings.Join(flag.Allowed, ", ") + ")"
	}
	if flag.Deprecated != "" {
		desc += " [DEPRECATED: " + flag.Deprecated + "]"
	}
	if flag.Value != nil {
		if def := flag.Value.Default(); def != "" {
			desc += " (Default: " + def + ")"
			if bv, ok := flag.Value.(core.StrictBool); ok && bv.IsStrictBool() {
				desc += " (Strict)"
			}
		}
	}
	if shouldInjectEnvKey(flag, globalHideEnvs, prefix) {
		flag.EnvKey = strings.ToUpper(prefix + "_" + strings.ReplaceAll(flag.Name, "-", "_"))
	}
	if shouldShowEnv(flag, globalHideEnvs) {
		desc += " (Env: " + flag.EnvKey + ")"
	}
	if !flag.HideRequired && flag.Required {
		desc += " (Required)"
	}
	if flag.Group != nil && !flag.Group.IsHidden() {
		desc += buildGroupInfo(flag.Group)
	}
	return desc
}

func shouldInjectEnvKey(flag *core.BaseFlag, globalHideEnvs bool, prefix string) bool {
	return !globalHideEnvs && !flag.DisableEnv && !flag.HideEnv && flag.EnvKey == "" && prefix != ""
}

func shouldShowEnv(flag *core.BaseFlag, globalHideEnvs bool) bool {
	return !globalHideEnvs && !flag.DisableEnv && !flag.HideEnv && flag.EnvKey != ""
}

func buildGroupInfo(group *core.MutualGroup) string {
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

// printUsageToken prints usage token like `-v`, `--verbose`, or both.
func printUsageToken(w io.Writer, fl *core.BaseFlag, mode FlagPrintMode) {
	meta := getPlaceholder(fl)
	switch mode {
	case PrintShort:
		if fl.Short != "" {
			fmt.Fprintf(w, " -%s", fl.Short)
			if meta != "" {
				fmt.Fprintf(w, " %s", meta)
			}
		}
	case PrintLong:
		fmt.Fprintf(w, " --%s", fl.Name)
		if meta != "" {
			fmt.Fprintf(w, " %s", meta)
		}
	case PrintBoth:
		if fl.Short != "" {
			fmt.Fprintf(w, " -%s|--%s", fl.Short, fl.Name)
		} else {
			fmt.Fprintf(w, " --%s", fl.Name)
		}
		if meta != "" {
			fmt.Fprintf(w, " %s", meta)
		}
	}
}
