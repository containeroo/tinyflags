package engine

import (
	"cmp"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/containeroo/tinyflags/internal/dynamic"
)

// FlagSet manages the definition, parsing, and usage output of command-line flags.
type FlagSet struct {
	name               string                           // Application or command name (used in usage)
	errorHandling      ErrorHandling                    // Behavior when parsing fails
	staticFlagsMap     map[string]*core.BaseFlag        // All registered static flags by name
	staticFlagsOrder   []*core.BaseFlag                 // Static flags in registration order
	dynamicGroupsMap   map[string]*dynamic.Group        // All dynamic groups by name
	dynamicGroupsOrder []*dynamic.Group                 // Dynamic groups in registration order
	oneOfGroup         []*core.OneOfGroupGroup          // All oneOfGroup groups
	allOrNoneGroup     []*core.AllOrNoneGroup           // All AllOrNoneGroup groups
	positional         []string                         // Remaining non-flag arguments
	requiredPositional int                              // Required positional argument count
	validatePositional func(string) error               // Function to validate positional arguments
	finalizePositional func(string) string              // Function to finalize positional arguments
	envPrefix          string                           // Optional ENV prefix (e.g. "APP_")
	envKeyFunc         EnvKeyFunc                       // Function to derive env keys from prefix+flag name  <-- NEW
	getEnv             func(string) string              // Function used to read ENV vars (default: os.Getenv)
	hideEnvs           bool                             // Globally hide environment key hints
	ignoreInvalidEnv   bool                             // Whether to ignore unknown ENV overrides
	defaultDelimiter   string                           // Global slice delimiter (default: ",")
	title              string                           // Title shown in usage output
	desc               string                           // Prolog before flags
	notes              string                           // Epilog after flags
	versionString      string                           // Version string for --version
	usagePrintMode     FlagPrintMode                    // Usage print mode (short|long|both|flags|none)
	output             io.Writer                        // Destination for help output
	enableHelp         bool                             // Whether built-in --help is enabled
	enableVer          bool                             // Whether built-in --version is enabled
	showHelp           *bool                            // Parsed value of --help
	helpText           string                           // Custom help text
	showVersion        *bool                            // Parsed value of --version
	versionText        string                           // Custom version text
	Usage              func()                           // Custom usage function (optional)
	sortFlags          bool                             // Enable static flag sorting
	sortGroups         bool                             // Enable dynamic group sorting
	authors            string                           // Optional authors block
	beforeParse        func([]string) ([]string, error) // Hook to preprocess args
	unknownFlag        func(string) error               // Handler for unknown flags

	// Indentation and width config for description
	descIndent int
	descWidth  int

	// Indentation and width config for static flags
	usageStaticNote   string
	usageStaticIndent int
	usageStaticCol    int
	usageStaticWidth  int

	// Indentation and width config for dynamic flags
	usageDynamicNote   string
	usageDynamicIndent int
	usageDynamicCol    int
	usageDynamicWidth  int

	// Indentation and width config for notes
	noteIndent int
	noteWidth  int
}

// NewFlagSet creates a new FlagSet with the given name and error handling policy.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	fs := &FlagSet{
		name:               name,
		errorHandling:      errorHandling,
		staticFlagsMap:     make(map[string]*core.BaseFlag),
		getEnv:             os.Getenv,
		envKeyFunc:         NewReplacerEnvKeyFunc(strings.NewReplacer("-", "_", ".", "_", "/", "_"), true),
		enableHelp:         true,
		enableVer:          true,
		defaultDelimiter:   ",",
		output:             os.Stdout,
		usagePrintMode:     PrintFlags,
		descIndent:         0,
		descWidth:          400,
		usageStaticIndent:  4,
		usageDynamicIndent: 4,
		noteIndent:         0,
		noteWidth:          400,
		title:              "Flags:",
	}

	fs.Usage = func() {
		out := fs.Output()
		if fs.usageStaticCol == 0 {
			fs.usageStaticCol = fs.StaticAutoUsageColumn(1)
		}
		if fs.usageDynamicCol == 0 {
			fs.usageDynamicCol = fs.DynamicAutoUsageColumn(1)
		}

		fs.PrintUsage(out, fs.usagePrintMode)
		fs.PrintTitle(out)
		fs.PrintAuthors(out)
		fs.PrintDescription(out, fs.descIndent, fs.descWidth)
		fs.PrintStaticDefaults(out, fs.usageStaticIndent, fs.usageStaticCol, fs.usageStaticWidth)
		fs.PrintDynamicDefaults(out, fs.usageDynamicIndent, fs.usageDynamicCol, fs.usageDynamicWidth)
		fs.PrintNotes(out, fs.noteIndent, fs.noteWidth)
	}

	return fs
}

// --- Metadata Configuration ---

func (f *FlagSet) Name() string                                    { return f.name }
func (f *FlagSet) EnvPrefix(prefix string)                         { f.envPrefix = prefix }
func (f *FlagSet) SetEnvKeyFunc(fn EnvKeyFunc)                     { f.envKeyFunc = fn }
func (f *FlagSet) EnvKeyForFlag(name string) string                { return f.envKeyFunc(f.envPrefix, name) }
func (f *FlagSet) DefaultDelimiter() string                        { return f.defaultDelimiter }
func (f *FlagSet) Globaldelimiter(s string)                        { f.defaultDelimiter = s }
func (f *FlagSet) BeforeParse(fn func([]string) ([]string, error)) { f.beforeParse = fn }
func (f *FlagSet) OnUnknownFlag(fn func(string) error)             { f.unknownFlag = fn }
func (f *FlagSet) Version(s string)                                { f.versionString = s; f.enableVer = true }
func (f *FlagSet) VersionText(s string)                            { f.versionText = s }
func (f *FlagSet) HelpText(s string)                               { f.helpText = s }
func (f *FlagSet) Title(s string)                                  { f.title = s }
func (f *FlagSet) Authors(s string)                                { f.authors = s }
func (f *FlagSet) Description(s string)                            { f.desc = s }
func (f *FlagSet) Note(s string)                                   { f.notes = s }
func (f *FlagSet) HideEnvs()                                       { f.hideEnvs = true }
func (f *FlagSet) DisableHelp()                                    { f.enableHelp = false }
func (f *FlagSet) DisableVersion()                                 { f.enableVer = false; f.versionString = "" }
func (f *FlagSet) SortedFlags(enable bool)                         { f.sortFlags = enable }
func (f *FlagSet) SortedGroups(enable bool)                        { f.sortGroups = enable }
func (f *FlagSet) SetOutput(w io.Writer)                           { f.output = w }
func (f *FlagSet) Output() io.Writer                               { return f.output }
func (f *FlagSet) IgnoreInvalidEnv(enable bool)                    { f.ignoreInvalidEnv = enable }
func (f *FlagSet) SetGetEnvFn(fn func(string) string)              { f.getEnv = fn }

// --- Positional Arguments ---

func (f *FlagSet) RequirePositional(n int) { f.requiredPositional = n }
func (f *FlagSet) Args() []string          { return f.positional }
func (f *FlagSet) Arg(i int) (string, bool) {
	if i >= 0 && i < len(f.positional) {
		return f.positional[i], true
	}
	return "", false
}
func (f *FlagSet) SetPositionalValidate(fn func(string) error)  { f.validatePositional = fn }
func (f *FlagSet) SetPositionalFinalize(fn func(string) string) { f.finalizePositional = fn }

// --- Usage Formatting Configuration ---

func (f *FlagSet) SetDescIndent(n int)  { f.descIndent = n }
func (f *FlagSet) DescIndent() int      { return f.descIndent }
func (f *FlagSet) SetDescWidth(max int) { f.descWidth = max }
func (f *FlagSet) DescWidth() int       { return f.descWidth }

func (f *FlagSet) SetStaticUsageIndent(n int)            { f.usageStaticIndent = n }
func (f *FlagSet) StaticUsageIndent() int                { return f.usageStaticIndent }
func (f *FlagSet) SetStaticUsageColumn(col int)          { f.usageStaticCol = col }
func (f *FlagSet) StaticUsageColumn() int                { return f.usageStaticCol }
func (f *FlagSet) SetStaticUsageWidth(maxWidth int)      { f.usageStaticWidth = maxWidth }
func (f *FlagSet) StaticUsageWidth() int                 { return f.usageStaticWidth }
func (f *FlagSet) StaticAutoUsageColumn(padding int) int { return f.calcStaticUsageColumn(padding) }
func (f *FlagSet) SetStaticUsageNote(s string)           { f.usageStaticNote = s }
func (f *FlagSet) StaticUsageNote() string               { return f.usageStaticNote }

func (f *FlagSet) SetDynamicUsageIndent(n int)            { f.usageDynamicIndent = n }
func (f *FlagSet) DynamicUsageIndent() int                { return f.usageDynamicIndent }
func (f *FlagSet) SetDynamicUsageColumn(col int)          { f.usageDynamicCol = col }
func (f *FlagSet) DynamicUsageColumn() int                { return f.usageDynamicCol }
func (f *FlagSet) SetDynamicUsageWidth(max int)           { f.usageDynamicWidth = max }
func (f *FlagSet) DynamicUsageWidth() int                 { return f.usageDynamicWidth }
func (f *FlagSet) DynamicAutoUsageColumn(padding int) int { return f.calcDynamicUsageColumn(padding) }
func (f *FlagSet) SetDynamicUsageNote(s string)           { f.usageDynamicNote = s }
func (f *FlagSet) DynamicUsageNote() string               { return f.usageDynamicNote }

func (f *FlagSet) SetNoteIndent(n int)  { f.noteIndent = n }
func (f *FlagSet) NoteIndent() int      { return f.noteIndent }
func (f *FlagSet) SetNoteWidth(max int) { f.noteWidth = max }
func (f *FlagSet) NoteWidth() int       { return f.noteWidth }

// --- Flag & Group Registration ---

func (f *FlagSet) RegisterFlag(name string, bf *core.BaseFlag) {
	f.staticFlagsMap[name] = bf
	f.staticFlagsOrder = append(f.staticFlagsOrder, bf)
}

func (f *FlagSet) LookupFlag(name string) *core.BaseFlag {
	return f.staticFlagsMap[name]
}

func (f *FlagSet) OrderedStaticFlags() []*core.BaseFlag {
	all := make([]*core.BaseFlag, 0, len(f.staticFlagsMap))
	for _, fl := range f.staticFlagsMap {
		all = append(all, fl)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })
	return all
}

func (f *FlagSet) DynamicGroup(name string) *dynamic.Group {
	if f.dynamicGroupsMap == nil {
		f.dynamicGroupsMap = make(map[string]*dynamic.Group)
	}
	if g, ok := f.dynamicGroupsMap[name]; ok {
		return g
	}
	g := dynamic.NewGroup(f, name)
	f.dynamicGroupsMap[name] = g
	f.dynamicGroupsOrder = append(f.dynamicGroupsOrder, g)
	return g
}

func (f *FlagSet) DynamicGroups() []*dynamic.Group {
	return f.dynamicGroupsOrder
}

func (f *FlagSet) OrderedDynamicGroups() []*dynamic.Group {
	groups := make([]*dynamic.Group, len(f.dynamicGroupsOrder))
	copy(groups, f.dynamicGroupsOrder)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name() < groups[j].Name()
	})
	return groups
}

// --- AllOrNone Group Handling ---

func (f *FlagSet) AllOrNoneGroups() []*core.AllOrNoneGroup { return f.allOrNoneGroup }

func (f *FlagSet) AddAllOrNoneGroup(name string, g *core.AllOrNoneGroup) {
	f.allOrNoneGroup = append(f.allOrNoneGroup, g)
}

// GetAllOrNoneGroup returns a group that requires all flags to be set together.
// It creates the group if it doesn't exist.
func (f *FlagSet) GetAllOrNoneGroup(name string) *core.AllOrNoneGroup {
	for _, g := range f.allOrNoneGroup {
		if g.Name == name {
			return g
		}
	}
	g := &core.AllOrNoneGroup{Name: name}
	f.allOrNoneGroup = append(f.allOrNoneGroup, g)
	return g
}

// AttachToAllOrNoneGroup attaches a flag to a require-together group.
func (f *FlagSet) AttachToAllOrNoneGroup(bf *core.BaseFlag, group string) {
	g := f.GetAllOrNoneGroup(group)
	g.Flags = append(g.Flags, bf)
}

// AttachGroupToAllOrNone nests one AllOrNone group into another.
func (f *FlagSet) AttachGroupToAllOrNone(parent string, child string) {
	p := f.GetAllOrNoneGroup(parent)
	c := f.GetAllOrNoneGroup(child)
	p.AddGroup(c)
}

// --- Mutual Group Handling ---

func (f *FlagSet) OneOfGroups() []*core.OneOfGroupGroup {
	return f.oneOfGroup
}

func (f *FlagSet) AddOneOfGroup(name string, g *core.OneOfGroupGroup) {
	f.oneOfGroup = append(f.oneOfGroup, g)
}

func (f *FlagSet) GetOneOfGroup(name string) *core.OneOfGroupGroup {
	for _, g := range f.oneOfGroup {
		if g.Name == name {
			return g
		}
	}
	group := &core.OneOfGroupGroup{Name: name}
	f.oneOfGroup = append(f.oneOfGroup, group)
	return group
}

func (f *FlagSet) AttachToOneOfGroup(bf *core.BaseFlag, group string) {
	g := f.GetOneOfGroup(group)
	g.Flags = append(g.Flags, bf)
	bf.OneOfGroup = g
}

// AttachGroupToOneOf adds an AllOrNone group as a single OneOf choice.
func (f *FlagSet) AttachGroupToOneOf(group string, aon string) {
	g := f.GetOneOfGroup(group)
	ag := f.GetAllOrNoneGroup(aon)
	g.AddGroup(ag)
}

// --- Builtin Flags ---

func (f *FlagSet) maybeAddBuiltinFlags() {
	if f.enableHelp && f.showHelp == nil {
		if _, exists := f.staticFlagsMap["help"]; !exists {
			f.showHelp = f.Bool("help", false, cmp.Or(f.helpText, "Show help")).Short("h").DisableEnv().Value()
		}
	}
	if f.enableVer && f.showVersion == nil && f.versionString != "" {
		if _, exists := f.staticFlagsMap["version"]; !exists {
			f.showVersion = f.Bool("version", false, cmp.Or(f.versionText, "Show version")).DisableEnv().Value()
		}
	}
}
