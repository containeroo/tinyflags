package tinyflags

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

// FlagPrintMode defines how flags are displayed in usage.
type FlagPrintMode int

const (
	PrintShort FlagPrintMode = iota // Show only short flags: -x
	PrintLong                       // Show only long flags: --xyz
	PrintBoth                       // Show both: -x|--xyz
	PrintFlags                      // Prints "Flags"
	PrintNone                       // Hide all flags
)

// PrintDefaults prints the detailed help for all flags.
func (f *FlagSet) PrintDefaults() {
	out := f.Output()
	w := tabwriter.NewWriter(out, 2, 4, 2, ' ', 0)

	userFlags, versionFlag, helpFlag := f.splitFlags()

	if f.sortFlags {
		f.printSortedFlags(w, userFlags, versionFlag, helpFlag)
	} else {
		f.printUnsortedFlags(w, userFlags, versionFlag, helpFlag)
	}

	w.Flush() // nolint:errcheck
}

// PrintUsage prints the one-line usage header based on mode.
func (f *FlagSet) PrintUsage(w io.Writer, mode FlagPrintMode) {
	fmt.Fprint(w, "Usage: "+f.Name()) // nolint:errcheck

	switch mode {
	case PrintNone:
		fmt.Fprintln(w) // nolint:errcheck
		return
	case PrintFlags:
		fmt.Fprintln(w, " [flags]") // nolint:errcheck
		return
	}

	var userFlags []*baseFlag
	var versionFlag, helpFlag *baseFlag

	for _, fl := range f.orderedFlags() {
		switch {
		case fl.name == "version" && f.enableVer:
			versionFlag = fl
		case fl.name == "help" && f.enableHelp:
			helpFlag = fl
		default:
			userFlags = append(userFlags, fl)
		}
	}

	for _, fl := range userFlags {
		if fl.hidden {
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

// PrintTitle prints the configured title above the usage.
func (f *FlagSet) PrintTitle(w io.Writer) {
	if f.title != "" {
		fmt.Fprintln(w, f.title) // nolint:errcheck
	}
}

// PrintNotes prints epilog text below the usage block.
func (f *FlagSet) PrintNotes(w io.Writer, width int) {
	if f.notes != "" {
		fmt.Fprintln(w, wrapText(f.notes, width)) // nolint:errcheck
	}
}

// PrintDescription prints prolog text above the flags.
func (f *FlagSet) PrintDescription(w io.Writer, width int) {
	if f.desc != "" {
		fmt.Fprintln(w, wrapText(f.desc, width)) // nolint:errcheck
	}
}

func (f *FlagSet) splitFlags() (userFlags []*baseFlag, versionFlag, helpFlag *baseFlag) {
	for _, fl := range f.orderedFlags() {
		switch {
		case fl.name == "version" && f.enableVer:
			versionFlag = fl
		case fl.name == "help" && f.enableHelp:
			helpFlag = fl
		default:
			userFlags = append(userFlags, fl)
		}
	}
	return
}

func (f *FlagSet) printSortedFlags(w io.Writer, userFlags []*baseFlag, versionFlag, helpFlag *baseFlag) {
	all := append([]*baseFlag{}, userFlags...)
	if versionFlag != nil {
		all = append(all, versionFlag)
	}
	if helpFlag != nil {
		all = append(all, helpFlag)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].name < all[j].name })
	for _, fl := range all {
		printFlagUsage(w, f.descIndent, f.descMaxLen, fl, f.envPrefix)
	}
}

func (f *FlagSet) printUnsortedFlags(w io.Writer, userFlags []*baseFlag, versionFlag, helpFlag *baseFlag) {
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

// printFlagUsage writes a single flag's help entry.
func printFlagUsage(w io.Writer, descIndent, maxDesc int, flag *baseFlag, prefix string) {
	var b strings.Builder

	// Format flag name
	if flag.short != "" {
		b.WriteString("  -")
		b.WriteString(flag.short)
		b.WriteString(", ")
	} else {
		b.WriteString("      ")
	}
	b.WriteString("--")
	b.WriteString(flag.name)

	if meta := getMetavar(flag); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}

	flagLine := b.String()

	// Build description metadata
	desc := flag.usage
	if len(flag.allowed) > 0 {
		desc += " (Allowed: " + strings.Join(flag.allowed, ", ") + ")"
	}
	if flag.deprecated != "" {
		desc += " [DEPRECATED: " + flag.deprecated + "]"
	}
	if def := flag.value.Default(); def != "" {
		if bv, ok := flag.value.(*BoolValue); ok && bv.IsStrictBool() {
			desc += " (Default: " + def + ")"
		} else if !ok {
			desc += " (Default: " + def + ")"
		}
	}
	if !flag.disableEnv && flag.envKey == "" && prefix != "" {
		flag.envKey = strings.ToUpper(prefix + "_" + strings.ReplaceAll(flag.name, "-", "_"))
	}
	if !flag.disableEnv && flag.envKey != "" {
		desc += " (Env: " + flag.envKey + ")"
	}
	if flag.required {
		desc += " (Required)"
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

func getMetavar(flag *baseFlag) string {
	isBool := false
	isStrict := false

	if bv, ok := flag.value.(*BoolValue); ok {
		isBool = true
		isStrict = bv.IsStrictBool()
	}

	if isBool && !isStrict {
		return ""
	}

	var meta string
	if flag.metavar != "" {
		meta = flag.metavar
	} else if len(flag.allowed) > 0 {
		meta = "<" + strings.Join(flag.allowed, "|") + ">"
	} else if isStrict {
		meta = "<true|false>"
	} else {
		meta = strings.ToUpper(flag.name)
	}

	if _, ok := flag.value.(SliceMarker); ok {
		meta += "..."
	}

	return meta
}

func printUsageToken(w io.Writer, fl *baseFlag, mode FlagPrintMode) {
	meta := getMetavar(fl)

	switch mode {
	case PrintShort:
		if fl.short != "" {
			fmt.Fprintf(w, " -%s", fl.short) // nolint:errcheck
			if meta != "" {
				fmt.Fprintf(w, " %s", meta) // nolint:errcheck
			}
		}
	case PrintLong:
		fmt.Fprintf(w, " --%s", fl.name) // nolint:errcheck
		if meta != "" {
			fmt.Fprintf(w, " %s", meta) // nolint:errcheck
		}
	case PrintBoth:
		if fl.short != "" {
			fmt.Fprintf(w, " -%s|--%s", fl.short, fl.name) // nolint:errcheck
		} else {
			fmt.Fprintf(w, " --%s", fl.name) // nolint:errcheck
		}
		if meta != "" {
			fmt.Fprintf(w, " %s", meta) // nolint:errcheck
		}
	}
}

// orderedFlags returns all flags in sorted or registration order.
func (f *FlagSet) orderedFlags() []*baseFlag {
	var userFlags []*baseFlag
	var versionFlag, helpFlag *baseFlag
	var all []*baseFlag

	if f.sortFlags {
		// collect and sort
		for _, fl := range f.flags {
			all = append(all, fl)
		}
		sort.Slice(all, func(i, j int) bool { return all[i].name < all[j].name })
	} else {
		// use registered order
		all = f.registered
	}

	for _, fl := range all {
		switch fl.name {
		case "version":
			if f.enableVer {
				versionFlag = fl
			}
		case "help":
			if f.enableHelp {
				helpFlag = fl
			}
		default:
			userFlags = append(userFlags, fl)
		}
	}

	if versionFlag != nil {
		userFlags = append(userFlags, versionFlag)
	}
	if helpFlag != nil {
		userFlags = append(userFlags, helpFlag)
	}
	return userFlags
}
