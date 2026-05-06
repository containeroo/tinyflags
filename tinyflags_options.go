package tinyflags

import "github.com/containeroo/tinyflags/internal/engine"

// HelpOptions groups high-level help/metadata settings.
type HelpOptions struct{ impl *engine.FlagSet }

// Title sets the help title.
func (h *HelpOptions) Title(s string) { h.impl.Title(s) }

// Authors sets the help authors text.
func (h *HelpOptions) Authors(s string) { h.impl.Authors(s) }

// Description sets the help description text.
func (h *HelpOptions) Description(s string) { h.impl.Description(s) }

// Note sets the help footer note.
func (h *HelpOptions) Note(s string) { h.impl.Note(s) }

// HelpText sets the built-in help flag text.
func (h *HelpOptions) HelpText(s string) { h.impl.HelpText(s) }

// DisableHelp disables the built-in help flag.
func (h *HelpOptions) DisableHelp() { h.impl.DisableHelp() }

// DisableVersion disables the built-in version flag.
func (h *HelpOptions) DisableVersion() { h.impl.DisableVersion() }

// VersionText sets the built-in version flag text.
func (h *HelpOptions) VersionText(s string) { h.impl.VersionText(s) }

// LayoutOptions groups usage/indent/width configuration.
type LayoutOptions struct{ impl *engine.FlagSet }

// SetDescIndent sets the description indentation.
func (l *LayoutOptions) SetDescIndent(n int) { l.impl.SetDescIndent(n) }

// SetDescWidth sets the description wrap width.
func (l *LayoutOptions) SetDescWidth(max int) { l.impl.SetDescWidth(max) }

// SetStaticUsageIndent sets the static usage indentation.
func (l *LayoutOptions) SetStaticUsageIndent(n int) { l.impl.SetStaticUsageIndent(n) }

// SetStaticUsageColumn sets the static usage description column.
func (l *LayoutOptions) SetStaticUsageColumn(col int) {
	l.impl.SetStaticUsageColumn(col)
}

// SetStaticUsageWidth sets the static usage wrap width.
func (l *LayoutOptions) SetStaticUsageWidth(maxWidth int) {
	l.impl.SetStaticUsageWidth(maxWidth)
}

// SetStaticUsageNote sets the static usage note.
func (l *LayoutOptions) SetStaticUsageNote(s string) { l.impl.SetStaticUsageNote(s) }

// SetDynamicUsageIndent sets the dynamic usage indentation.
func (l *LayoutOptions) SetDynamicUsageIndent(n int) { l.impl.SetDynamicUsageIndent(n) }

// SetDynamicUsageColumn sets the dynamic usage description column.
func (l *LayoutOptions) SetDynamicUsageColumn(col int) {
	l.impl.SetDynamicUsageColumn(col)
}

// SetDynamicUsageWidth sets the dynamic usage wrap width.
func (l *LayoutOptions) SetDynamicUsageWidth(max int) { l.impl.SetDynamicUsageWidth(max) }

// SetDynamicUsageNote sets the dynamic usage note.
func (l *LayoutOptions) SetDynamicUsageNote(s string) { l.impl.SetDynamicUsageNote(s) }

// SetNoteIndent sets the note indentation.
func (l *LayoutOptions) SetNoteIndent(n int) { l.impl.SetNoteIndent(n) }

// SetNoteWidth sets the note wrap width.
func (l *LayoutOptions) SetNoteWidth(max int) { l.impl.SetNoteWidth(max) }

// Help returns grouped helpers for configuring help/usage output.
func (f *FlagSet) Help() *HelpOptions { return &HelpOptions{impl: f.impl} }

// Layout returns grouped helpers for configuring usage layout.
func (f *FlagSet) Layout() *LayoutOptions { return &LayoutOptions{impl: f.impl} }
