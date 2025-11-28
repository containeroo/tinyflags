package tinyflags

import "github.com/containeroo/tinyflags/internal/engine"

// HelpOptions groups high-level help/metadata settings.
type HelpOptions struct{ impl *engine.FlagSet }

func (h *HelpOptions) Title(s string)       { h.impl.Title(s) }
func (h *HelpOptions) Authors(s string)     { h.impl.Authors(s) }
func (h *HelpOptions) Description(s string) { h.impl.Description(s) }
func (h *HelpOptions) Note(s string)        { h.impl.Note(s) }
func (h *HelpOptions) HelpText(s string)    { h.impl.HelpText(s) }
func (h *HelpOptions) DisableHelp()         { h.impl.DisableHelp() }
func (h *HelpOptions) DisableVersion()      { h.impl.DisableVersion() }
func (h *HelpOptions) VersionText(s string) { h.impl.VersionText(s) }

// LayoutOptions groups usage/indent/width configuration.
type LayoutOptions struct{ impl *engine.FlagSet }

func (l *LayoutOptions) SetDescIndent(n int)        { l.impl.SetDescIndent(n) }
func (l *LayoutOptions) SetDescWidth(max int)       { l.impl.SetDescWidth(max) }
func (l *LayoutOptions) SetStaticUsageIndent(n int) { l.impl.SetStaticUsageIndent(n) }
func (l *LayoutOptions) SetStaticUsageColumn(col int) {
	l.impl.SetStaticUsageColumn(col)
}
func (l *LayoutOptions) SetStaticUsageWidth(maxWidth int) {
	l.impl.SetStaticUsageWidth(maxWidth)
}
func (l *LayoutOptions) SetStaticUsageNote(s string) { l.impl.SetStaticUsageNote(s) }

func (l *LayoutOptions) SetDynamicUsageIndent(n int) { l.impl.SetDynamicUsageIndent(n) }
func (l *LayoutOptions) SetDynamicUsageColumn(col int) {
	l.impl.SetDynamicUsageColumn(col)
}
func (l *LayoutOptions) SetDynamicUsageWidth(max int) { l.impl.SetDynamicUsageWidth(max) }
func (l *LayoutOptions) SetDynamicUsageNote(s string) { l.impl.SetDynamicUsageNote(s) }

func (l *LayoutOptions) SetNoteIndent(n int)  { l.impl.SetNoteIndent(n) }
func (l *LayoutOptions) SetNoteWidth(max int) { l.impl.SetNoteWidth(max) }

// Help returns grouped helpers for configuring help/usage output.
func (f *FlagSet) Help() *HelpOptions { return &HelpOptions{impl: f.impl} }

// Layout returns grouped helpers for configuring usage layout.
func (f *FlagSet) Layout() *LayoutOptions { return &LayoutOptions{impl: f.impl} }
