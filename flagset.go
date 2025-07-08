package tinyflags

import (
	"fmt"
	"io"
	"os"
)

// FlagSet stores configuration and registered flags for parsing.
type FlagSet struct {
	envPrefix          string               // optional ENV key prefix (e.g., APP_)
	errorHandling      ErrorHandling        // behavior on parse errors
	flags              map[string]*baseFlag // all registered flags by name
	enableHelp         bool                 // whether to add built-in --help
	enableVer          bool                 // whether to add built-in --version
	showHelp           *bool                // internal help flag pointer
	showVersion        *bool                // internal version flag pointer
	groups             []*mutualGroup       // registered mutual exclusion groups
	versionString      string               // version string for --version
	getEnv             func(string) string  // env lookup (defaults to os.Getenv)
	defaultDelimiter   string               // default delimiter for slice values
	positional         []string             // captured positional args
	requiredPositional int                  // number of required positional args
	ignoreInvalidEnv   bool                 // ignore bad ENV overrides
	name               string               // command name for usage
	title              string               // optional usage title
	desc               string               // optional text before usage
	notes              string               // optional text after usage
	sortFlags          bool                 // whether to sort flags in usage
	Usage              func()               // custom usage printer
	usagePrintMode     FlagPrintMode        // how to print usage
	descMaxLen         int                  // max line length for descriptions
	descIndent         int                  // indentation for descriptions
	output             io.Writer            // output writer for usage/help
	registered         []*baseFlag          // ordered list of flags
}

// NewFlagSet creates a new FlagSet with a name and error behavior.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	fs := &FlagSet{
		name:             name,
		flags:            make(map[string]*baseFlag),
		getEnv:           os.Getenv,
		errorHandling:    errorHandling,
		ignoreInvalidEnv: false,
		enableHelp:       true,
		enableVer:        true,
		descIndent:       40,
		descMaxLen:       100,
		usagePrintMode:   PrintFlags,
		title:            "\nFlags:",
	}
	fs.Usage = func() {
		out := fs.Output()
		fs.PrintUsage(out, fs.usagePrintMode)
		fs.PrintTitle(out)
		fs.PrintDescription(out, 100)
		fs.PrintDefaults()
		fs.PrintNotes(out, 100)
	}
	fs.SetOutput(os.Stdout)
	return fs
}

// maybeAddBuiltinFlags adds --help and --version if enabled and not user-defined.
func (f *FlagSet) maybeAddBuiltinFlags() {
	// Only add default help if not already defined by the user
	if f.enableHelp && f.showHelp == nil {
		if _, exists := f.flags["help"]; !exists {
			var help bool
			f.BoolVarP(&help, "help", "h", false, "show help").DisableEnv()
			f.showHelp = &help
		}
	}
	// Only add version if a version string was actually set and flag is not already present
	if f.enableVer && f.showVersion == nil && f.versionString != "" {
		if _, exists := f.flags["version"]; !exists {
			var ver bool
			f.BoolVar(&ver, "version", false, "show version").DisableEnv()
			f.showVersion = &ver
		}
	}
}

// EnvPrefix sets the environment prefix for all flags.
func (f *FlagSet) EnvPrefix(prefix string) { f.envPrefix = prefix }

// Version enables the version flag with the given string.
func (f *FlagSet) Version(s string) {
	f.versionString = s
	f.enableVer = true
}

// Name returns the program name.
func (f *FlagSet) Name() string { return f.name }

// Title sets the optional title for usage output.
func (f *FlagSet) Title(s string) { f.title = s }

// Description sets optional text before usage help.
func (f *FlagSet) Description(s string) { f.desc = s }

// Note sets optional text after usage help.
func (f *FlagSet) Note(s string) { f.notes = s }

// DisableHelp disables automatic help flag registration.
func (f *FlagSet) DisableHelp() { f.enableHelp = false }

// DisableVersion disables automatic version flag registration.
func (f *FlagSet) DisableVersion() {
	f.enableVer = false
	f.versionString = "" // ensures it won't show up
}

// Sorted enables or disables sorting of flag help output.
func (f *FlagSet) Sorted(enable bool) { f.sortFlags = enable }

// SetOutput sets the writer for help and usage messages.
func (f *FlagSet) SetOutput(w io.Writer) { f.output = w }

// Output returns the writer for usage and error output.
func (f *FlagSet) Output() io.Writer { return f.output }

// IgnoreInvalidEnv skips env vars that cannot be parsed.
func (f *FlagSet) IgnoreInvalidEnv(enable bool) { f.ignoreInvalidEnv = enable }

// SetGetEnvFn overrides the env lookup function.
func (f *FlagSet) SetGetEnvFn(fn func(string) string) { f.getEnv = fn }

// Globaldelimiter sets the default delimiter for slice flags.
func (f *FlagSet) Globaldelimiter(s string) { f.defaultDelimiter = s }

// RequirePositional sets the number of required positional args.
func (f *FlagSet) RequirePositional(n int) {
	f.requiredPositional = n
}

// Args returns all captured positional arguments.
func (f *FlagSet) Args() []string {
	return f.positional
}

// Arg returns the i-th positional argument.
func (f *FlagSet) Arg(i int) (string, bool) {
	if i >= 0 && i < len(f.positional) {
		return f.positional[i], true
	}
	return "", false
}

// Group adds a flag to a named mutual exclusion group.
func (f *FlagSet) Group(name string, flag *baseFlag) {
	// Don't re-register the flag here.
	if name != "" {
		for _, g := range f.groups {
			if g.name == name {
				g.flags = append(g.flags, flag)
				flag.group = g
				return
			}
		}
		g := &mutualGroup{name: name, flags: []*baseFlag{flag}}
		f.groups = append(f.groups, g)
		flag.group = g
	}
}

// UsagePrintMode sets the mode for printing usage.
func (f *FlagSet) UsagePrintMode(mode FlagPrintMode) {
	f.usagePrintMode = mode
}

// DescriptionMaxLen sets the max line length for descriptions.
func (f *FlagSet) DescriptionMaxLen(line int) {
	f.descMaxLen = line
}

// DescriptionIndent sets the indentation for descriptions.
func (f *FlagSet) DescriptionIndent(indent int) {
	f.descIndent = indent
}

// Get returns the parsed value for a named flag.
func (f *FlagSet) Get(name string) (any, error) {
	flag, ok := f.flags[name]
	if !ok {
		return nil, fmt.Errorf("flag %q not found", name)
	}
	return flag.value.Get(), nil
}

// MustGet is like Get but panics if the flag is missing.
func (f *FlagSet) MustGet(name string) any {
	v, err := f.Get(name)
	if err != nil {
		panic(err)
	}
	return v
}

// GetAs returns the typed value for a flag or an error.
func GetAs[T any](fs *FlagSet, name string) (T, error) {
	v, err := fs.Get(name)
	if err != nil {
		var zero T
		return zero, err
	}
	return v.(T), nil
}
