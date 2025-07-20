package engine

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/containeroo/tinyflags/internal/core"
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
func (f *FlagSet) PrintDefaults() {
	w := tabwriter.NewWriter(f.Output(), 2, 4, 2, ' ', 0)
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
		printFlagUsage(w, f.descIndent, f.descMaxLen, fl, f.envPrefix)
	}
}

// printUnsortedFlags prints flags in registration order.
func (f *FlagSet) printUnsortedFlags(w io.Writer, userFlags []*core.BaseFlag, versionFlag, helpFlag *core.BaseFlag) {
	for _, fl := range userFlags {
		printFlagUsage(w, f.descIndent, f.descMaxLen, fl, f.envPrefix)
	}
	if versionFlag != nil {
		printFlagUsage(w, f.descIndent, f.descMaxLen, versionFlag, f.envPrefix)
	}
	if helpFlag != nil {
		printFlagUsage(w, f.descIndent, f.descMaxLen, helpFlag, f.envPrefix)
	}
}

// printFlagUsage renders a single flag's full usage line.
func printFlagUsage(w io.Writer, descIndent, maxDesc int, flag *core.BaseFlag, prefix string) {
	var b strings.Builder

	if flag.Short != "" {
		b.WriteString("  -")
		b.WriteString(flag.Short)
		b.WriteString(", ")
	} else {
		b.WriteString("      ")
	}
	b.WriteString("--")
	b.WriteString(flag.Name)

	if meta := getMetavar(flag); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}

	flagLine := b.String()
	desc := flag.Usage

	if len(flag.Allowed) > 0 {
		desc += " (Allowed: " + strings.Join(flag.Allowed, ", ") + ")"
	}
	if flag.Deprecated != "" {
		desc += " [DEPRECATED: " + flag.Deprecated + "]"
	}
	if def := flag.Value.Default(); def != "" {
		desc += " (Default: " + def + ")"
		if bv, ok := flag.Value.(core.StrictBool); ok && bv.IsStrictBool() {
			desc += " (Strict)"
		}
	}
	if !flag.DisableEnv && flag.EnvKey == "" && prefix != "" {
		flag.EnvKey = strings.ToUpper(prefix + "_" + strings.ReplaceAll(flag.Name, "-", "_"))
	}
	if !flag.DisableEnv && flag.EnvKey != "" {
		desc += " (Env: " + flag.EnvKey + ")"
	}
	if flag.Required {
		desc += " (Required)"
	}
	if flag.Group != nil && !flag.Group.IsHidden() {
		desc += " (Group: "
		if flag.Group.TitleText() != "" {
			desc += flag.Group.TitleText()
		} else {
			desc += flag.Group.Name
		}
		if flag.Group.IsRequired() {
			desc += ", required"
		}
		desc += ")"
	}

	// Decide whether to wrap or keep on same line
	if len(flagLine)+len(desc) <= maxDesc {
		fmt.Fprintf(w, "%-*s %s\n", descIndent, flagLine, desc) // nolint:errcheck
		return
	}

	// Wrap description if too long
	wrapped := wrapText(desc, maxDesc-descIndent-1)
	lines := strings.Split(wrapped, "\n")

	fmt.Fprintf(w, "%-*s %s\n", descIndent, flagLine, lines[0]) // nolint:errcheck
	for _, l := range lines[1:] {
		fmt.Fprintf(w, "%-*s %s\n", descIndent, "", l) // nolint:errcheck
	}
}

// getMetavar returns the placeholder for a flag value shown in help.
func getMetavar(flag *core.BaseFlag) string {
	isBool := false
	isStrict := false

	if bv, ok := flag.Value.(core.StrictBool); ok {
		isBool = true
		isStrict = bv.IsStrictBool()
	}

	if isBool && !isStrict {
		return ""
	}

	var meta string
	if flag.Metavar != "" {
		meta = flag.Metavar
	} else if len(flag.Allowed) > 0 {
		meta = "<" + strings.Join(flag.Allowed, "|") + ">"
	} else if isStrict {
		meta = "<true|false>"
	} else {
		meta = strings.ToUpper(flag.Name)
	}

	if _, ok := flag.Value.(core.SliceMarker); ok {
		meta += "..."
	}
	return meta
}

// printUsageToken writes a short usage token (e.g. -v, --verbose) to the usage line.
func printUsageToken(w io.Writer, fl *core.BaseFlag, mode FlagPrintMode) {
	meta := getMetavar(fl)

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
