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
	PrintShort FlagPrintMode = iota // Only short flags (e.g., -v)
	PrintLong                       // Only long flags (e.g., --verbose)
	PrintBoth                       // Combined short|long (e.g., -v|--verbose)
	PrintFlags                      // Just "Usage: [flags]"
	PrintNone                       // No usage tokens shown
)

// PrintUsage prints a usage line including flags depending on mode.
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

// PrintTitle writes the usage title heading.
func (f *FlagSet) PrintTitle(w io.Writer) {
	if f.title != "" {
		fmt.Fprintln(w, f.title) // nolint:errcheck
	}
}

// PrintAuthors writes the usage authors heading.
func (f *FlagSet) PrintAuthors(w io.Writer) {
	if f.authors != "" {
		fmt.Fprintln(w, "Authors: "+f.authors) // nolint:errcheck
	}
}

// PrintDescription renders the prolog text above flags.
func (f *FlagSet) PrintDescription(w io.Writer, width int) {
	if f.desc != "" {
		fmt.Fprintln(w, wrapText(f.desc, width)) // nolint:errcheck
	}
}

// PrintNotes renders the epilog text below flags.
func (f *FlagSet) PrintNotes(w io.Writer, width int) {
	if f.notes != "" {
		fmt.Fprintln(w, wrapText(f.notes, width)) // nolint:errcheck
	}
}

// PrintDefaults prints all visible flags with their descriptions.
func (f *FlagSet) PrintDefaults(out io.Writer, width int) {
	w := tabwriter.NewWriter(out, 2, 4, 2, ' ', 0)
	userFlags, versionFlag, helpFlag := f.splitFlags()

	if f.sortFlags {
		f.printSortedFlags(w, userFlags, versionFlag, helpFlag)
	} else {
		f.printUnsortedFlags(w, userFlags, versionFlag, helpFlag)
	}
	w.Flush() // nolint:errcheck
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

// printSortedFlags prints all flags alphabetically.
func (f *FlagSet) printSortedFlags(w io.Writer, userFlags []*core.BaseFlag, versionFlag, helpFlag *core.BaseFlag) {
	all := append([]*core.BaseFlag{}, userFlags...)
	if versionFlag != nil {
		all = append(all, versionFlag)
	}
	if helpFlag != nil {
		all = append(all, helpFlag)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })
	for _, fl := range all {
		printFlagUsage(w, f.descIndent, f.descMaxLen, f.hideEnvs, fl, f.envPrefix)
	}
}

// printUnsortedFlags prints flags in registration order.
func (f *FlagSet) printUnsortedFlags(w io.Writer, userFlags []*core.BaseFlag, versionFlag, helpFlag *core.BaseFlag) {
	for _, fl := range userFlags {
		printFlagUsage(w, f.descIndent, f.descMaxLen, f.hideEnvs, fl, f.envPrefix)
	}
	if versionFlag != nil {
		printFlagUsage(w, f.descIndent, f.descMaxLen, f.hideEnvs, versionFlag, f.envPrefix)
	}
	if helpFlag != nil {
		printFlagUsage(w, f.descIndent, f.descMaxLen, f.hideEnvs, helpFlag, f.envPrefix)
	}
}

// printFlagUsage renders a single flag's full usage line.
func printFlagUsage(w io.Writer, descIndent, maxDesc int, globalHideEnvs bool, flag *core.BaseFlag, prefix string) {
	var b strings.Builder

	//	formatFlagNames(&b, flag)
	if meta := getPlaceholder(flag); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}
	flagLine := b.String()

	desc := buildFlagDescription(flag, globalHideEnvs, prefix)

	if len(flagLine)+len(desc) <= maxDesc {
		fmt.Fprintf(w, "%-*s %s\n", descIndent, flagLine, desc) // nolint:errcheck
		return
	}

	wrapped := wrapText(desc, maxDesc-descIndent-1)
	lines := strings.Split(wrapped, "\n")

	fmt.Fprintf(w, "%-*s %s\n", descIndent, flagLine, lines[0]) // nolint:errcheck
	for _, l := range lines[1:] {
		fmt.Fprintf(w, "%-*s %s\n", descIndent, "", l) // nolint:errcheck
	}
}

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

// getPlaceholder returns the placeholder for a flag value shown in help.
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

	var placeholder string
	if flag.Placeholder != "" {
		placeholder = flag.Placeholder
	} else if len(flag.Allowed) > 0 {
		placeholder = "<" + strings.Join(flag.Allowed, "|") + ">"
	} else if isStrict {
		placeholder = "<true|false>"
	} else {
		placeholder = strings.ToUpper(flag.Name)
	}

	if _, ok := flag.Value.(core.SliceMarker); ok {
		placeholder += "..."
	}
	return placeholder
}

// printUsageToken writes a short usage token (e.g. -v, --verbose) to the usage line.
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

// orderedFlags returns all registered flags in desired display order.
func (f *FlagSet) orderedFlags() []*core.BaseFlag {
	var all []*core.BaseFlag

	if f.sortFlags {
		for _, fl := range f.flags {
			all = append(all, fl)
		}
		sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })
	} else {
		all = f.registered
	}

	return all
}

// PrintDynamicDefaults renders all dynamic groups and their flags.
func (f *FlagSet) PrintDynamicDefaults(w io.Writer, width int) {
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

		if title := g.TitleText(); title != "" {
			fmt.Fprintf(w, "\n%s\n", title)
		}
		if desc := g.DescriptionText(); desc != "" {
			fmt.Fprintln(w, wrapText(desc, width))
		}

		ids := g.Instances()
		if len(ids) == 0 {
			// still print static flag layout if there’s no instance
			ids = []string{"<id>"}
		}

		for _, id := range ids {
			flags := g.Flags()
			if g.IsFlagSorted() {
				sort.Slice(flags, func(i, j int) bool {
					return flags[i].Name < flags[j].Name
				})
			}

			tw := tabwriter.NewWriter(w, 2, 4, 2, ' ', 0)
			for _, fl := range flags {
				fmt.Fprintf(tw, "      --%s.%s.%s", name, id, fl.Name)
				if meta := getPlaceholder(fl); meta != "" {
					fmt.Fprintf(tw, " %s", meta)
				}
				desc := buildFlagDescription(fl, f.hideEnvs, f.envPrefix)
				fmt.Fprintf(tw, "\t%s\n", desc)
			}
			tw.Flush()
		}

		if note := g.NoteText(); note != "" {
			fmt.Fprintln(w, wrapText(note, width))
		}
	}
}
