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

// BuildFlagDescription builds the help text for a static flag, including metadata suffixes.
func BuildFlagDescription(flag *core.BaseFlag, globalHideEnvs bool, prefix string) string {
	desc := buildFlagDescriptionPrefix(flag)

	flag.ResolveUsageEnvKey(prefix, globalHideEnvs)
	if flag.ShouldShowUsageEnv(globalHideEnvs) {
		desc += " (env: " + flag.EnvKey + ")"
	}

	return finishFlagDescription(desc, flag)
}

// BuildDynamicFlagDescription builds help text for a dynamic flag.
func BuildDynamicFlagDescription(flag *core.BaseFlag, globalHideEnvs bool, prefix, groupName, idPlaceholder string) string {
	desc := buildFlagDescriptionPrefix(flag)

	if envKey := dynamicUsageEnvKey(flag, globalHideEnvs, prefix, groupName, idPlaceholder); envKey != "" {
		desc += " (env: " + envKey + ")"
	}

	return finishFlagDescription(desc, flag)
}

// buildFlagDescriptionPrefix assembles metadata shared by static and dynamic help text.
func buildFlagDescriptionPrefix(flag *core.BaseFlag) string {
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

	return desc
}

// finishFlagDescription appends required and group annotations to a description.
func finishFlagDescription(desc string, flag *core.BaseFlag) string {
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

// dynamicUsageEnvKey returns the displayable environment key for a dynamic flag.
func dynamicUsageEnvKey(flag *core.BaseFlag, globalHideEnvs bool, prefix, groupName, idPlaceholder string) string {
	if flag == nil || globalHideEnvs || flag.DisableEnv || flag.HideEnv || prefix == "" {
		return ""
	}
	return core.DynamicEnvKey(prefix, groupName, idPlaceholder, flag.Name)
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
func PrintStaticDefaults(
	w io.Writer,
	flags []*core.BaseFlag,
	sectionOrder []string,
	indent, startCol, maxWidth int,
	hideEnvs bool,
	envPrefix, note string,
) {
	layout := newLayout(indent, startCol, maxWidth)
	sections := groupStaticFlags(flags)
	renderedSection := false

	for _, name := range sectionOrder {
		sectionFlags := sections[name]
		if len(sectionFlags) == 0 {
			continue
		}

		if renderedSection || name != "" {
			fmt.Fprintln(w) // nolint:errcheck
		}
		if name != "" {
			fmt.Fprintf(w, "%s:\n", name) // nolint:errcheck
		}
		for _, flag := range sectionFlags {
			printFlagUsage(w, layout, hideEnvs, flag, envPrefix)
		}
		renderedSection = true
	}

	if note != "" {
		fmt.Fprintln(w, note) // nolint:errcheck
	}
}

// groupStaticFlags groups visible static flags by section while preserving input order.
// The empty section name represents the unnamed section and renders without a heading.
func groupStaticFlags(flags []*core.BaseFlag) map[string][]*core.BaseFlag {
	sections := make(map[string][]*core.BaseFlag)

	for _, flag := range flags {
		if flag.Hidden {
			continue
		}
		sections[flag.Section] = append(sections[flag.Section], flag)
	}

	return sections
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
			desc := BuildDynamicFlagDescription(fl, hideEnvs, envPrefix, name, idPlaceholder)

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

// layout controls indentation and column widths for rendered help rows.
type layout struct {
	indent   int // Number of leading spaces for each row.
	startCol int // Width reserved for the flag label column.
	maxWidth int // Maximum total row width before wrapping.
}

// newLayout creates a layout with the supplied rendering dimensions.
func newLayout(indent, startCol, maxWidth int) layout {
	return layout{indent: indent, startCol: startCol, maxWidth: maxWidth}
}

// descriptionWidth returns the available width for a help description.
func (l layout) descriptionWidth() int {
	if l.maxWidth <= 0 {
		return 100
	}
	return max(l.maxWidth-l.indent-l.startCol-1, 1)
}

// writeWrappedRow renders a label and wraps its description beneath it.
func (l layout) writeWrappedRow(w io.Writer, label, desc string) {
	wrapped := WrapText(desc, l.descriptionWidth())
	lines := strings.Split(wrapped, "\n")

	fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", l.indent), l.startCol, label, lines[0]) // nolint:errcheck

	padding := strings.Repeat(" ", l.indent+l.startCol+1)
	for _, line := range lines[1:] {
		fmt.Fprintf(w, "%s%s\n", padding, line) // nolint:errcheck
	}
}

// writeIndented renders wrapped text with this layout's indentation.
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

// printFlagUsage renders one static flag's usage row.
func printFlagUsage(w io.Writer, layout layout, globalHideEnvs bool, flag *core.BaseFlag, prefix string) {
	var b strings.Builder
	formatStaticFlagNames(&b, flag)
	if meta := flag.UsagePlaceholder(); meta != "" {
		b.WriteString(" ")
		b.WriteString(meta)
	}
	layout.writeWrappedRow(w, b.String(), BuildFlagDescription(flag, globalHideEnvs, prefix))
}

// formatStaticFlagNames appends the short and long names of a static flag.
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

// formatDynamicFlagLine formats the command-line spelling of a dynamic flag.
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

// buildGroupInfo formats one one-of group annotation for help text.
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

// buildRequireGroupInfo formats one all-or-none group annotation for help text.
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
