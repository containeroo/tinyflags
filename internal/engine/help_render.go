package engine

import "strings"

// RenderHelpText renders help output without mutating parse state.
func (f *FlagSet) RenderHelpText() string {
	if f == nil {
		return ""
	}

	f.maybeAddBuiltinFlags()

	var buf strings.Builder
	prevOutput := f.Output()
	f.SetOutput(&buf)
	defer f.SetOutput(prevOutput)

	f.Usage()
	return buf.String()
}
