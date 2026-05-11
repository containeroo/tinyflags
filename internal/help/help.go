package help

import (
	"fmt"
	"io"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
)

// WrapText wraps s to the given width while preserving explicit newlines.
func WrapText(s string, width int) string {
	if width <= 0 || s == "" {
		return s
	}

	var out []string
	for _, paragraph := range strings.Split(s, "\n") {
		if paragraph == "" {
			out = append(out, "")
			continue
		}
		var line strings.Builder
		for word := range strings.FieldsSeq(paragraph) {
			if line.Len() == 0 {
				line.WriteString(word)
				continue
			}
			if line.Len()+1+len(word) > width {
				out = append(out, line.String())
				line.Reset()
				line.WriteString(word)
				continue
			}
			line.WriteByte(' ')
			line.WriteString(word)
		}
		if line.Len() > 0 {
			out = append(out, line.String())
		}
	}
	return strings.Join(out, "\n")
}

// BuildFlagDescription builds the help text for a flag, including metadata suffixes.
func BuildFlagDescription(flag *core.BaseFlag, globalHideEnvs bool, prefix string) string {
	desc := flag.Usage

	allowed := flag.AllowedValues()
	if !flag.HideAllowed && len(allowed) > 0 {
		desc += " (allowed: " + strings.Join(allowed, ", ") + ")"
	}

	if flag.Deprecated != "" {
		desc += " (deprecated: " + flag.Deprecated + ")"
	}

	if flag.ShouldShowDefaultInHelp() {
		if def := flag.Value.Default(); def != "" {
			desc += " (default: " + def + ")"
		}
	}

	if !flag.HideRequires && len(flag.Requires) > 0 {
		desc += " (requires: " + strings.Join(flag.Requires, ", ") + ")"
	}

	flag.ResolveUsageEnvKey(prefix, globalHideEnvs)

	if flag.ShouldShowUsageEnv(globalHideEnvs) {
		desc += " (env: " + flag.EnvKey + ")"
	}

	if !flag.HideRequired && flag.Required {
		desc += " (required)"
	}

	for _, group := range flag.VisibleOneOfGroups() {
		desc += buildGroupInfo(group)
	}

	if flag.AllOrNone != nil && !flag.AllOrNone.IsHidden() {
		desc += buildRequireGroupInfo(flag.AllOrNone)
	}

	return desc
}

// CalcStaticUsageColumn calculates the maximum static flag label width.
func CalcStaticUsageColumn(flags []*core.BaseFlag, padding int) int {
	maxFlagLen := 0
	for _, fl := range flags {
		var b strings.Builder
		formatStaticFlagNames(&b, fl)
		if meta := fl.UsagePlaceholder(); meta != "" {
			b.WriteString(" ")
			b.WriteString(meta)
		}
		if l := len(b.String()); l > maxFlagLen {
			maxFlagLen = l
		}
	}
	return maxFlagLen + padding
}

// CalcDynamicUsageColumn calculates the maximum dynamic flag label width.
func CalcDynamicUsageColumn(groups []*dynamic.Group, padding int) int {
	maxLen := 0
	for _, group := range groups {
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

// WriteIndented writes wrapped text with a fixed indentation.
func WriteIndented(w io.Writer, text string, indent, maxWidth int) {
	newLayout(indent, 0, maxWidth).writeIndented(w, text)
}

// PrintStaticDefaults renders all static flags with help descriptions.
func PrintStaticDefaults(w io.Writer, flags []*core.BaseFlag, indent, startCol, maxWidth int, hideEnvs bool, envPrefix, note string) {
	layout := newLayout(indent, startCol, maxWidth)
	var lastSection string
	for _, fl := range flags {
		if fl.Hidden {
			continue
		}
		if fl.Section != "" && fl.Section != lastSection {
			fmt.Fprintf(w, "\n%s:\n", fl.Section) // nolint:errcheck
			lastSection = fl.Section
		}
		printFlagUsage(w, layout, hideEnvs, fl, envPrefix)
	}

	if note != "" {
		fmt.Fprintln(w, note) // nolint:errcheck
	}
}

// PrintDynamicDefaults renders all dynamic groups with help descriptions.
func PrintDynamicDefaults(w io.Writer, groups []*dynamic.Group, indent, startCol, maxWidth int, hideEnvs bool, envPrefix, note string) {
	layout := newLayout(indent, startCol, maxWidth)
	for _, group := range groups {
		if group.IsHidden() {
			continue
		}
		name := group.Name()

		if title := group.TitleText(); title != "" {
			fmt.Fprintf(w, "\n%s\n", title) // nolint:errcheck
		}
		if desc := group.DescriptionText(); desc != "" {
			newLayout(0, 0, maxWidth).writeIndented(w, WrapText(desc, maxWidth-indent))
		}
		idPlaceholder := group.GetPlaceholder()
		if idPlaceholder == "" {
			idPlaceholder = "<ID>"
		}

		for _, fl := range group.DynamicFlags() {
			flagLine := formatDynamicFlagLine(name, idPlaceholder, fl)
			desc := BuildFlagDescription(fl, hideEnvs, envPrefix)

			if len(desc) <= layout.descriptionWidth() {
				fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", indent), startCol, flagLine, desc) // nolint:errcheck
				continue
			}

			layout.writeWrappedRow(w, flagLine, desc)
		}

		if groupNote := group.NoteText(); groupNote != "" {
			newLayout(indent, 0, maxWidth).writeIndented(w, WrapText(groupNote, maxWidth-indent))
		}
	}

	if note != "" {
		fmt.Fprintln(w, note) // nolint:errcheck
	}
}

type layout struct {
	indent   int
	startCol int
	maxWidth int
}

func newLayout(indent, startCol, maxWidth int) layout {
	return layout{indent: indent, startCol: startCol, maxWidth: maxWidth}
}

func (l layout) descriptionWidth() int {
	return max(l.maxWidth-l.indent-l.startCol-1, 100)
}

func (l layout) writeWrappedRow(w io.Writer, label, desc string) {
	wrapped := WrapText(desc, l.descriptionWidth())
	lines := strings.Split(wrapped, "\n")

	fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", l.indent), l.startCol, label, lines[0]) // nolint:errcheck

	padding := strings.Repeat(" ", l.indent+l.startCol+1)
	for _, line := range lines[1:] {
		fmt.Fprintf(w, "%s%s\n", padding, line) // nolint:errcheck
	}
}

func (l layout) writeIndented(w io.Writer, text string) {
	if text == "" {
		return
	}

	prefix := strings.Repeat(" ", l.indent)
	wrapped := WrapText(text, l.maxWidth-l.indent)
	for _, line := range strings.Split(wrapped, "\n") {
		if line == "" {
			fmt.Fprintln(w) // nolint:errcheck
			continue
		}
		fmt.Fprintf(w, "%s%s\n", prefix, line) // nolint:errcheck
	}
}

func printFlagUsage(w io.Writer, layout layout, globalHideEnvs bool, flag *core.BaseFlag, prefix string) {
	var b strings.Builder
	formatStaticFlagNames(&b, flag)
	if meta := flag.UsagePlaceholder(); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}
	layout.writeWrappedRow(w, b.String(), BuildFlagDescription(flag, globalHideEnvs, prefix))
}

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

func formatDynamicFlagLine(groupName, idPlaceholder string, fl *core.BaseFlag) string {
	var b strings.Builder
	b.WriteString("--")
	b.WriteString(groupName)
	b.WriteString(".")
	b.WriteString(idPlaceholder)
	b.WriteString(".")
	b.WriteString(fl.Name)
	if meta := fl.UsagePlaceholder(); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}
	return b.String()
}

func buildGroupInfo(group *core.OneOfGroupGroup) string {
	var b strings.Builder
	b.WriteString(" [group: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	b.WriteString(" (one of")
	if group.IsRequired() {
		b.WriteString(", required")
	}
	b.WriteString(")]")
	return b.String()
}

func buildRequireGroupInfo(group *core.AllOrNoneGroup) string {
	var b strings.Builder
	b.WriteString(" [group: ")
	if group.TitleText() != "" {
		b.WriteString(group.TitleText())
	} else {
		b.WriteString(group.Name)
	}
	b.WriteString(" (all or none")
	if group.IsRequired() {
		b.WriteString(", required")
	}
	b.WriteString(")]")
	return b.String()
}
