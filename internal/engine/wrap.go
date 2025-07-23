package engine

import "strings"

// wrapText wraps the input string s into lines no longer than width.
func wrapText(s string, width int) string {
	words := strings.Fields(s) // split into words by whitespace
	if len(words) == 0 {
		return ""
	}

	var lines []string
	line := words[0] // start with the first word

	// add words to the current line until width is exceeded
	for _, word := range words[1:] {
		if len(line)+len(word)+1 > width {
			lines = append(lines, line) // push current line
			line = word                 // start new line
			continue
		}
		line += " " + word // append word to current line
	}

	lines = append(lines, line)      // push the final line
	return strings.Join(lines, "\n") // join lines with newline
}
