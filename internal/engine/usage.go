package engine

import "github.com/containeroo/tinyflags/internal/core"

// FlagPrintMode defines how usage should be rendered.
type FlagPrintMode int

const (
	PrintShort FlagPrintMode = iota // e.g. -v
	PrintLong                       // e.g. --verbose
	PrintBoth                       // e.g. -v|--verbose
	PrintFlags                      // prints only [flags]
	PrintNone                       // prints nothing
)

// splitFlags separates user-defined and built-in flags.
func (f *FlagSet) splitFlags() (userFlags []*core.BaseFlag, versionFlag, helpFlag *core.BaseFlag) {
	for _, fl := range f.staticFlags() {
		switch {
		case fl.Name == "version" && f.enableVer:
			versionFlag = fl
		case fl.Name == "help" && f.enableHelp:
			helpFlag = fl
		default:
			userFlags = append(userFlags, fl)
		}
	}
	return
}
