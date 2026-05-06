package engine

import (
	"fmt"
	"io"
	"strings"
)

type usageLayout struct {
	indent   int
	startCol int
	maxWidth int
}

func newUsageLayout(indent, startCol, maxWidth int) usageLayout {
	return usageLayout{
		indent:   indent,
		startCol: startCol,
		maxWidth: maxWidth,
	}
}

func (l usageLayout) descriptionWidth() int {
	return max(l.maxWidth-l.indent-l.startCol-1, 100)
}

func (l usageLayout) writeWrappedRow(w io.Writer, label, desc string) {
	wrapped := wrapText(desc, l.descriptionWidth())
	lines := strings.Split(wrapped, "\n")

	fmt.Fprintf(w, "%s%-*s %s\n", strings.Repeat(" ", l.indent), l.startCol, label, lines[0]) // nolint:errcheck

	padding := strings.Repeat(" ", l.indent+l.startCol+1)
	for _, line := range lines[1:] {
		fmt.Fprintf(w, "%s%s\n", padding, line) // nolint:errcheck
	}
}

func (l usageLayout) writeIndented(w io.Writer, text string) {
	if text == "" {
		return
	}

	prefix := strings.Repeat(" ", l.indent)
	wrapped := wrapText(text, l.maxWidth-l.indent)

	for _, line := range strings.Split(wrapped, "\n") {
		if line == "" {
			fmt.Fprintln(w) // nolint:errcheck
			continue
		}
		fmt.Fprintf(w, "%s%s\n", prefix, line) // nolint:errcheck
	}
}
