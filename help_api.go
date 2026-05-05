package tinyflags

import "io"

// HideEnvs disables all env-var annotations in help output.
func (f *FlagSet) HideEnvs() { f.impl.HideEnvs() }

// Title sets the main title shown in usage output.
func (f *FlagSet) Title(s string) { f.Help().Title(s) }

// Authors sets the list of authors printed in help output.
func (f *FlagSet) Authors(s string) { f.Help().Authors(s) }

// Description sets the top description section of the help output.
func (f *FlagSet) Description(s string) { f.Help().Description(s) }

// Note sets the bottom note section of the help output.
func (f *FlagSet) Note(s string) { f.Help().Note(s) }

// HelpText sets the --help text.
func (f *FlagSet) HelpText(s string) { f.Help().HelpText(s) }

// DisableHelp disables the automatic --help flag.
func (f *FlagSet) DisableHelp() { f.Help().DisableHelp() }

// DisableVersion disables the automatic --version flag.
func (f *FlagSet) DisableVersion() { f.Help().DisableVersion() }

// SortedFlags enables sorted help output for static flags.
func (f *FlagSet) SortedFlags() { f.impl.SortedFlags(true) }

// SortedGroups enables sorted help output for dynamic groups.
func (f *FlagSet) SortedGroups() { f.impl.SortedGroups(true) }

// SetOneOfGroupVerbose toggles verbose OneOfGroup error messages.
func (f *FlagSet) SetOneOfGroupVerbose(enable bool) { f.impl.SetOneOfGroupVerbose(enable) }

// OneOfGroupVerbose reports whether OneOfGroup errors include conflicting flags.
func (f *FlagSet) OneOfGroupVerbose() bool { return f.impl.OneOfGroupVerbose() }

// SetOutput changes the destination writer for usage and error messages.
func (f *FlagSet) SetOutput(w io.Writer) { f.impl.SetOutput(w) }

// Output returns the current output writer.
func (f *FlagSet) Output() io.Writer { return f.impl.Output() }

// PrintUsage renders the top usage line.
func (f *FlagSet) PrintUsage(w io.Writer, mode FlagPrintMode) {
	f.impl.PrintUsage(w, mode)
}

// PrintTitle renders the title above all help content.
func (f *FlagSet) PrintTitle(w io.Writer) { f.impl.PrintTitle(w) }

// PrintAuthors renders the author line if set.
func (f *FlagSet) PrintAuthors(w io.Writer) { f.impl.PrintAuthors(w) }

// PrintDescription renders the full description block.
func (f *FlagSet) PrintDescription(w io.Writer, indent, width int) {
	f.impl.PrintDescription(w, indent, width)
}

// PrintStaticDefaults renders all static flag usage lines.
func (f *FlagSet) PrintStaticDefaults(w io.Writer, indent, col, width int) {
	f.impl.PrintStaticDefaults(w, indent, col, width)
}

// PrintDynamicDefaults renders all dynamic flag usage lines.
func (f *FlagSet) PrintDynamicDefaults(w io.Writer, indent, col, width int) {
	f.impl.PrintDynamicDefaults(w, indent, col, width)
}

// PrintNotes renders the notes section, if configured.
func (f *FlagSet) PrintNotes(w io.Writer, indent, width int) {
	f.impl.PrintNotes(w, indent, width)
}
