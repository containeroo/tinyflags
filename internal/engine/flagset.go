package engine

import (
	"fmt"
	"io"
	"os"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
)

// FlagSet manages the definition, parsing, and usage output of command-line flags.
type FlagSet struct {
	name               string                                  // name of the application or command (used in usage output).
	errorHandling      ErrorHandling                           // errorHandling determines what happens when parsing fails.
	flags              map[string]*core.BaseFlag               // flags holds all registered named flags by their long name.
	registered         []*core.BaseFlag                        // registered keeps the order in which flags were added (for ordered output).
	groups             []*core.MutualGroup                     // groups holds mutual exclusion groups (e.g. only one of a set of flags is allowed).
	dynamic            map[string]map[string]core.DynamicValue // dynamic holds all dynamically defined flags grouped by group name and field name.
	positional         []string                                // positional stores remaining positional arguments after flag parsing.
	requiredPositional int                                     // requiredPositional defines how many positional arguments must be provided.
	envPrefix          string                                  // envPrefix is the optional prefix applied to environment variable lookups (e.g. "APP_").
	getEnv             func(string) string                     // getEnv is the function used to look up environment variables (default: os.Getenv).
	ignoreInvalidEnv   bool                                    // ignoreInvalidEnv skips unknown or invalid environment overrides.
	defaultDelimiter   string                                  // defaultDelimiter is the global delimiter for slice flags (default: ",").
	title              string                                  // title is printed before the list of flags in usage output.
	desc               string                                  // desc is printed as a prolog above the flags.
	notes              string                                  // notes is printed as an epilog below the flags.
	versionString      string                                  // versionString is shown when --version is triggered.
	usagePrintMode     FlagPrintMode                           // usagePrintMode controls what is printed in PrintUsage.
	descMaxLen         int                                     // descMaxLen controls the max line length before wrapping.
	descIndent         int                                     // descIndent controls the left indent of flag descriptions.
	output             io.Writer                               // output is where usage output is written (default: os.Stdout).
	enableHelp         bool                                    // enableHelp toggles whether the built-in --help flag is added automatically.
	enableVer          bool                                    // enableVer toggles whether the built-in --version flag is added automatically.
	showHelp           *bool                                   // showHelp is a pointer to the parsed --help flag value, if enabled.
	showVersion        *bool                                   // showVersion is a pointer to the parsed --version flag value, if enabled.
	Usage              func()                                  // Usage is the customizable function for printing usage. Defaults to printing title, description, flags, and notes.
	sortFlags          bool                                    // sortFlags determines whether flags are printed in sorted order.
}

// NewFlagSet creates a new FlagSet with the given name and error handling policy.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	fs := &FlagSet{
		name:             name,
		flags:            make(map[string]*core.BaseFlag),
		getEnv:           os.Getenv,
		errorHandling:    errorHandling,
		ignoreInvalidEnv: false,
		enableHelp:       true,
		enableVer:        true,
		descIndent:       40,
		descMaxLen:       100,
		defaultDelimiter: ",",
		usagePrintMode:   PrintFlags,
		title:            "Flags:",
		output:           os.Stdout,
	}

	// Define a default usage function
	fs.Usage = func() {
		out := fs.Output()
		fs.PrintUsage(out, fs.usagePrintMode)
		fs.PrintTitle(out)
		fs.PrintDescription(out, fs.descMaxLen)
		fs.PrintDefaults()
		fs.PrintNotes(out, fs.descMaxLen)
	}

	return fs
}

// Version sets the version string to enable the --version flag.
func (f *FlagSet) Version(s string) {
	f.versionString = s
	f.enableVer = true
}

// EnvPrefix sets a prefix to be prepended to all environment variables.
func (f *FlagSet) EnvPrefix(prefix string) { f.envPrefix = prefix }

// Title sets the usage section title.
func (f *FlagSet) Title(s string) { f.title = s }

// Description sets the prolog text shown above the flags.
func (f *FlagSet) Description(s string) { f.desc = s }

// Note sets the epilog text shown below the flags.
func (f *FlagSet) Note(s string) { f.notes = s }

// DisableHelp disables the automatic --help flag.
func (f *FlagSet) DisableHelp() { f.enableHelp = false }

// DisableVersion disables the automatic --version flag.
func (f *FlagSet) DisableVersion() {
	f.enableVer = false
	f.versionString = ""
}

// Sorted enables or disables sorted help output.
func (f *FlagSet) Sorted(enable bool) { f.sortFlags = enable }

// SetOutput sets the writer for help output.
func (f *FlagSet) SetOutput(w io.Writer) { f.output = w }

// Output returns the configured help writer.
func (f *FlagSet) Output() io.Writer { return f.output }

// IgnoreInvalidEnv controls whether unknown environment values cause errors.
func (f *FlagSet) IgnoreInvalidEnv(enable bool) { f.ignoreInvalidEnv = enable }

// SetGetEnvFn sets a custom environment lookup function.
func (f *FlagSet) SetGetEnvFn(fn func(string) string) { f.getEnv = fn }

// Globaldelimiter sets the delimiter used for slice flags.
func (f *FlagSet) Globaldelimiter(s string) { f.defaultDelimiter = s }

// RequirePositional sets how many positional arguments are required.
func (f *FlagSet) RequirePositional(n int) { f.requiredPositional = n }

// Args returns all remaining positional arguments.
func (f *FlagSet) Args() []string { return f.positional }

// Arg returns the positional argument at index i.
func (f *FlagSet) Arg(i int) (string, bool) {
	if i >= 0 && i < len(f.positional) {
		return f.positional[i], true
	}
	return "", false
}

// DescriptionMaxLen sets the maximum line width for wrapped descriptions.
func (f *FlagSet) DescriptionMaxLen(line int) { f.descMaxLen = line }

// DescriptionIndent sets the number of spaces before descriptions.
func (f *FlagSet) DescriptionIndent(indent int) { f.descIndent = indent }

// DefaultDelimiter returns the current default delimiter.
func (f *FlagSet) DefaultDelimiter() string { return f.defaultDelimiter }

// RegisterFlag registers a static flag.
func (f *FlagSet) RegisterFlag(name string, bf *core.BaseFlag) {
	f.flags[name] = bf
	f.registered = append(f.registered, bf)
}

// RegisterDynamic registers a dynamic flag for the given group and field.
func (f *FlagSet) RegisterDynamic(group, field string, val core.DynamicValue) error {
	if f.dynamic == nil {
		f.dynamic = make(map[string]map[string]core.DynamicValue)
	}
	if _, ok := f.dynamic[group]; !ok {
		f.dynamic[group] = make(map[string]core.DynamicValue)
	}
	if _, exists := f.dynamic[group][field]; exists {
		return fmt.Errorf("dynamic flag already registered: %s.%s", group, field)
	}
	f.dynamic[group][field] = val
	return nil
}

// DynamicGroup creates a new dynamic group with the given prefix.
func (f *FlagSet) DynamicGroup(name string) *dynamic.Group {
	return dynamic.NewGroup(f, name)
}

// GetGroup returns a mutual exclusion group by name (creating it if necessary).
func (f *FlagSet) GetGroup(name string) *core.MutualGroup {
	for _, g := range f.groups {
		if g.Name == name {
			return g
		}
	}
	group := &core.MutualGroup{Name: name}
	f.groups = append(f.groups, group)
	return group
}

// AddGroup manually adds a mutual group.
func (f *FlagSet) AddGroup(name string, g *core.MutualGroup) {
	f.groups = append(f.groups, g)
}

// Groups returns all mutual exclusion groups.
func (f *FlagSet) Groups() []*core.MutualGroup {
	return f.groups
}

// AttachToGroup connects a flag to a mutual exclusion group.
func (f *FlagSet) AttachToGroup(bf *core.BaseFlag, group string) {
	g := f.GetGroup(group)
	g.Flags = append(g.Flags, bf)
	bf.Group = g
}

// maybeAddBuiltinFlags adds --help and --version if enabled and not already defined.
func (f *FlagSet) maybeAddBuiltinFlags() {
	if f.enableHelp && f.showHelp == nil {
		if _, exists := f.flags["help"]; !exists {
			f.showHelp = f.BoolP("help", "h", false, "show help").DisableEnv().Value()
		}
	}
	if f.enableVer && f.showVersion == nil && f.versionString != "" {
		if _, exists := f.flags["version"]; !exists {
			f.showVersion = f.Bool("version", false, "show version").DisableEnv().Value()
		}
	}
}
