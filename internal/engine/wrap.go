package engine

import "strings"

// wrapText wraps the input string s into lines no longer than width while
// preserving explicit newlines provided by the caller.
func wrapText(s string, width int) string {
	if width <= 0 {
		return s
	}

	lines := strings.Split(s, "\n")
	var wrapped []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			wrapped = append(wrapped, "")
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			wrapped = append(wrapped, "")
			continue
		}

		cur := words[0]
		for _, word := range words[1:] {
			if len(cur)+len(word)+1 > width {
				wrapped = append(wrapped, cur)
				cur = word
				continue
			}
			cur += " " + word
		}
		wrapped = append(wrapped, cur)
	}

	return strings.Join(wrapped, "\n")
}
