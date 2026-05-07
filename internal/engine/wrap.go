package engine

import "github.com/containeroo/tinyflags/internal/help"

// wrapText wraps the input string s into lines no longer than width while
// preserving explicit newlines provided by the caller.
func wrapText(s string, width int) string {
	return help.WrapText(s, width)
}
