package engine

import (
	"io"
	"os"
	"sort"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
)

// FlagSet manages the definition, parsing, and usage output of command-line flags.
type FlagSet struct {
	name               string                    // Application or command name (used in usage)
	errorHandling      ErrorHandling             // Behavior when parsing fails
	staticFlagsMap     map[string]*core.BaseFlag // All registered static flags by name
	staticFlagsOrder   []*core.BaseFlag          // Registration order of static flags
	dynamicGroupsOrder []*dynamic.Group          // Ordered list of dynamic groups
	dynamicGroupsMap   map[string]*dynamic.Group // All dynamic groups by name
	groups             []*core.MutualGroup       // Registered mutual exclusion groups
	positional         []string                  // Remaining non-flag arguments
	requiredPositional int                       // Required positional argument count
	envPrefix          string                    // Optional ENV prefix (e.g. "APP_")
	getEnv             func(string) string       // Function used to read ENV vars (default: os.Getenv)
	ignoreInvalidEnv   bool                      // Whether to ignore unknown ENV overrides
	defaultDelimiter   string                    // Global slice delimiter (default: ",")
	title              string                    // Printed above flags in usage
	desc               string                    // Prolog before flag list
	notes              string                    // Epilog after flag list
	versionString      string                    // --version output string
	usagePrintMode     FlagPrintMode             // Usage output mode
	descMaxLen         int                       // Max length before description wrapping
	descIndent         int                       // Left indent for wrapped lines
	output             io.Writer                 // Destination for help output
	enableHelp         bool                      // Whether --help is enabled
	enableVer          bool                      // Whether --version is enabled
	showHelp           *bool                     // Parsed value of --help
	showVersion        *bool                     // Parsed value of --version
	Usage              func()                    // Custom usage function
	sortFlags          bool                      // Sort static flags
	sortGroups         bool                      // Sort dynamic groups
	authors            string                    // Optional authors string
	hideEnvs           bool                      // Hide environment info in help
}

// NewFlagSet creates a new FlagSet with the given name and error handling policy.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	fs := &FlagSet{
		name:             name,
		staticFlagsMap:   make(map[string]*core.BaseFlag),
		getEnv:           os.Getenv,
		errorHandling:    errorHandling,
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
		fs.PrintAuthors(out)
		fs.PrintDescription(out, fs.descMaxLen)
		fs.PrintDefaults(out, fs.descMaxLen)
		fs.PrintNotes(out, fs.descMaxLen)
	}

	return fs
}

// Name returns the name of the application.
func (f *FlagSet) Name() string { return f.name }

// Version sets the version string to enable the --version flag.
func (f *FlagSet) Version(s string) {
	f.versionString = s
	f.enableVer = true
}

// EnvPrefix sets a prefix to be prepended to all environment variables.
func (f *FlagSet) EnvPrefix(prefix string) { f.envPrefix = prefix }

// Title sets the usage section title.
func (f *FlagSet) Title(s string) { f.title = s }

// Authors sets the usage author block.
func (f *FlagSet) Authors(s string) { f.authors = s }

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

// HideEnvs disables environment variable display in help.
func (f *FlagSet) HideEnvs() { f.hideEnvs = true }

// SortedFlags enables or disables sorted help output.
func (f *FlagSet) SortedFlags(enable bool) { f.sortFlags = enable }

// SortedGroups enables or disables sorted group output.
func (f *FlagSet) SortedGroups(enable bool) { f.sortGroups = enable }

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
	f.staticFlagsMap[name] = bf
	f.staticFlagsOrder = append(f.staticFlagsOrder, bf)
}

// OrderedStaticFlags returns all static flags in sorted order.
func (f *FlagSet) OrderedStaticFlags() []*core.BaseFlag {
	var all []*core.BaseFlag
	for _, fl := range f.staticFlagsMap {
		all = append(all, fl)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Name < all[j].Name
	})
	return all
}

func (f *FlagSet) OrderedDynamicGroups() []*dynamic.Group {
	var groups []*dynamic.Group
	for _, g := range f.dynamicGroupsOrder {
		groups = append(groups, g)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].Name() < groups[j].Name()
	})
	return groups
}

// DynamicGroup creates a new dynamic group with the given prefix.
func (f *FlagSet) DynamicGroup(name string) *dynamic.Group {
	if f.dynamicGroupsMap == nil {
		f.dynamicGroupsMap = make(map[string]*dynamic.Group)
	}
	if g, ok := f.dynamicGroupsMap[name]; ok {
		return g
	}
	g := dynamic.NewGroup(f, name)
	f.dynamicGroupsMap[name] = g
	f.dynamicGroupsOrder = append(f.dynamicGroupsOrder, g) // track order
	return g
}

// DynamicGroups returns all dynamic groups in registration order.
func (f *FlagSet) DynamicGroups() []*dynamic.Group {
	return f.dynamicGroupsOrder
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

// AddGroup manually adds a mutual exclusion group.
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

// LookupFlag returns a registered static flag by name.
func (f *FlagSet) LookupFlag(name string) *core.BaseFlag {
	return f.staticFlagsMap[name]
}

// maybeAddBuiltinFlags adds --help and --version if enabled and not already defined.
func (f *FlagSet) maybeAddBuiltinFlags() {
	if f.enableHelp && f.showHelp == nil {
		if _, exists := f.staticFlagsMap["help"]; !exists {
			f.showHelp = f.Bool("help", false, "show help").Short("h").DisableEnv().Value()
		}
	}
	if f.enableVer && f.showVersion == nil && f.versionString != "" {
		if _, exists := f.staticFlagsMap["version"]; !exists {
			f.showVersion = f.Bool("version", false, "show version").DisableEnv().Value()
		}
	}
}
