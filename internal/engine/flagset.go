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
	allOrNoneGroup     []*core.AllOrNoneGroup           // All all-or-none groups
	positional         []string                         // Remaining non-flag arguments
	requiredPositional int                              // Required positional argument count
	validatePositional func(string) error               // Function to validate positional arguments
	finalizePositional func(string) string              // Function to finalize positional arguments
	envPrefix          string                           // Optional ENV prefix (e.g. "APP_")
	envKeyFunc         EnvKeyFunc                       // Function to derive env keys from prefix+flag name
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
	oneOfVerbose       bool                             // Include conflicting flags in OneOf errors
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
		oneOfVerbose:       true,
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

// Name returns the flag set name.
func (f *FlagSet) Name() string { return f.name }

// EnvPrefix sets the environment variable prefix.
func (f *FlagSet) EnvPrefix(prefix string) { f.envPrefix = prefix }

// SetEnvKeyFunc sets the environment variable naming function.
func (f *FlagSet) SetEnvKeyFunc(fn EnvKeyFunc) { f.envKeyFunc = fn }

// EnvKeyForFlag derives the environment key for a flag name.
func (f *FlagSet) EnvKeyForFlag(name string) string { return f.envKeyFunc(f.envPrefix, name) }

// DefaultDelimiter returns the default slice delimiter.
func (f *FlagSet) DefaultDelimiter() string { return f.defaultDelimiter }

// GlobalDelimiter sets the default slice delimiter.
func (f *FlagSet) GlobalDelimiter(s string) { f.defaultDelimiter = s }

// Globaldelimiter sets the default slice delimiter.
func (f *FlagSet) Globaldelimiter(s string) { f.defaultDelimiter = s }

// BeforeParse sets a hook that can rewrite args before parsing.
func (f *FlagSet) BeforeParse(fn func([]string) ([]string, error)) { f.beforeParse = fn }

// OnUnknownFlag sets the callback for unknown flags.
func (f *FlagSet) OnUnknownFlag(fn func(string) error) { f.unknownFlag = fn }

// Version enables the version flag and sets its output string.
func (f *FlagSet) Version(s string) { f.versionString = s; f.enableVer = true }

// VersionText sets the help text for the version flag.
func (f *FlagSet) VersionText(s string) { f.versionText = s }

// HelpText sets the help text for the help flag.
func (f *FlagSet) HelpText(s string) { f.helpText = s }

// Title sets the title shown above flag listings.
func (f *FlagSet) Title(s string) { f.title = s }

// Authors sets the authors block shown in help output.
func (f *FlagSet) Authors(s string) { f.authors = s }

// Description sets the prose shown before flags.
func (f *FlagSet) Description(s string) { f.desc = s }

// Note sets the prose shown after flags.
func (f *FlagSet) Note(s string) { f.notes = s }

// HideEnvs hides environment variable hints in help output.
func (f *FlagSet) HideEnvs() { f.hideEnvs = true }

// DisableHelp disables the built-in help flag.
func (f *FlagSet) DisableHelp() { f.enableHelp = false }

// DisableVersion disables the built-in version flag.
func (f *FlagSet) DisableVersion() { f.enableVer = false; f.versionString = "" }

// SortedFlags enables or disables sorted static help output.
func (f *FlagSet) SortedFlags(enable bool) { f.sortFlags = enable }

// SortedGroups enables or disables sorted dynamic help output.
func (f *FlagSet) SortedGroups(enable bool) { f.sortGroups = enable }

// SetOneOfGroupVerbose toggles verbose one-of validation errors.
func (f *FlagSet) SetOneOfGroupVerbose(enable bool) { f.oneOfVerbose = enable }

// OneOfGroupVerbose reports whether one-of validation is verbose.
func (f *FlagSet) OneOfGroupVerbose() bool { return f.oneOfVerbose }

// SetOutput sets the writer used for help output.
func (f *FlagSet) SetOutput(w io.Writer) { f.output = w }

// Output returns the configured help output writer.
func (f *FlagSet) Output() io.Writer { return f.output }

// IgnoreInvalidEnv toggles ignoring invalid environment overrides.
func (f *FlagSet) IgnoreInvalidEnv(enable bool) { f.ignoreInvalidEnv = enable }

// SetGetEnvFn replaces the environment lookup function.
func (f *FlagSet) SetGetEnvFn(fn func(string) string) { f.getEnv = fn }

// OverriddenValues returns the changed values after masking.
func (f *FlagSet) OverriddenValues() map[string]any { return f.overriddenValues() }

// --- Positional Arguments ---

// RequirePositional sets the required positional argument count.
func (f *FlagSet) RequirePositional(n int) { f.requiredPositional = n }

// Args returns the parsed positional arguments.
func (f *FlagSet) Args() []string { return f.positional }

// Arg returns the positional argument at index i.
func (f *FlagSet) Arg(i int) (string, bool) {
	if i >= 0 && i < len(f.positional) {
		return f.positional[i], true
	}
	return "", false
}

// SetPositionalValidate sets the positional validation hook.
func (f *FlagSet) SetPositionalValidate(fn func(string) error) { f.validatePositional = fn }

// SetPositionalFinalize sets the positional finalization hook.
func (f *FlagSet) SetPositionalFinalize(fn func(string) string) { f.finalizePositional = fn }

// --- Usage Formatting Configuration ---

// SetDescIndent sets the description indentation.
func (f *FlagSet) SetDescIndent(n int) { f.descIndent = n }

// DescIndent returns the description indentation.
func (f *FlagSet) DescIndent() int { return f.descIndent }

// SetDescWidth sets the description wrap width.
func (f *FlagSet) SetDescWidth(max int) { f.descWidth = max }

// DescWidth returns the description wrap width.
func (f *FlagSet) DescWidth() int { return f.descWidth }

// SetStaticUsageIndent sets the static usage indentation.
func (f *FlagSet) SetStaticUsageIndent(n int) { f.usageStaticIndent = n }

// StaticUsageIndent returns the static usage indentation.
func (f *FlagSet) StaticUsageIndent() int { return f.usageStaticIndent }

// SetStaticUsageColumn sets the static usage description column.
func (f *FlagSet) SetStaticUsageColumn(col int) { f.usageStaticCol = col }

// StaticUsageColumn returns the static usage description column.
func (f *FlagSet) StaticUsageColumn() int { return f.usageStaticCol }

// SetStaticUsageWidth sets the static usage wrap width.
func (f *FlagSet) SetStaticUsageWidth(maxWidth int) { f.usageStaticWidth = maxWidth }

// StaticUsageWidth returns the static usage wrap width.
func (f *FlagSet) StaticUsageWidth() int { return f.usageStaticWidth }

// StaticAutoUsageColumn calculates a static usage column automatically.
func (f *FlagSet) StaticAutoUsageColumn(padding int) int { return f.calcStaticUsageColumn(padding) }

// SetStaticUsageNote sets the note shown after static flags.
func (f *FlagSet) SetStaticUsageNote(s string) { f.usageStaticNote = s }

// StaticUsageNote returns the note shown after static flags.
func (f *FlagSet) StaticUsageNote() string { return f.usageStaticNote }

// SetDynamicUsageIndent sets the dynamic usage indentation.
func (f *FlagSet) SetDynamicUsageIndent(n int) { f.usageDynamicIndent = n }

// DynamicUsageIndent returns the dynamic usage indentation.
func (f *FlagSet) DynamicUsageIndent() int { return f.usageDynamicIndent }

// SetDynamicUsageColumn sets the dynamic usage description column.
func (f *FlagSet) SetDynamicUsageColumn(col int) { f.usageDynamicCol = col }

// DynamicUsageColumn returns the dynamic usage description column.
func (f *FlagSet) DynamicUsageColumn() int { return f.usageDynamicCol }

// SetDynamicUsageWidth sets the dynamic usage wrap width.
func (f *FlagSet) SetDynamicUsageWidth(max int) { f.usageDynamicWidth = max }

// DynamicUsageWidth returns the dynamic usage wrap width.
func (f *FlagSet) DynamicUsageWidth() int { return f.usageDynamicWidth }

// DynamicAutoUsageColumn calculates a dynamic usage column automatically.
func (f *FlagSet) DynamicAutoUsageColumn(padding int) int { return f.calcDynamicUsageColumn(padding) }

// SetDynamicUsageNote sets the note shown after dynamic flags.
func (f *FlagSet) SetDynamicUsageNote(s string) { f.usageDynamicNote = s }

// DynamicUsageNote returns the note shown after dynamic flags.
func (f *FlagSet) DynamicUsageNote() string { return f.usageDynamicNote }

// SetNoteIndent sets the notes indentation.
func (f *FlagSet) SetNoteIndent(n int) { f.noteIndent = n }

// NoteIndent returns the notes indentation.
func (f *FlagSet) NoteIndent() int { return f.noteIndent }

// SetNoteWidth sets the notes wrap width.
func (f *FlagSet) SetNoteWidth(max int) { f.noteWidth = max }

// NoteWidth returns the notes wrap width.
func (f *FlagSet) NoteWidth() int { return f.noteWidth }

// --- Flag & Group Registration ---

// RegisterFlag registers a static flag in the set.
func (f *FlagSet) RegisterFlag(name string, bf *core.BaseFlag) {
	f.staticFlagsMap[name] = bf
	f.staticFlagsOrder = append(f.staticFlagsOrder, bf)
}

// LookupFlag returns a registered static flag by name.
func (f *FlagSet) LookupFlag(name string) *core.BaseFlag {
	return f.staticFlagsMap[name]
}

// OrderedStaticFlags returns static flags sorted by name.
func (f *FlagSet) OrderedStaticFlags() []*core.BaseFlag {
	all := make([]*core.BaseFlag, 0, len(f.staticFlagsMap))
	for _, fl := range f.staticFlagsMap {
		all = append(all, fl)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })
	return all
}

// DynamicGroup returns or creates a dynamic flag group.
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

// overriddenValues returns masked values that changed during parsing.
func (f *FlagSet) overriddenValues() map[string]any {
	out := make(map[string]any)

	for _, fl := range f.staticFlagsMap {
		if fl.Value == nil || !fl.Value.Changed() {
			continue
		}
		val := fl.Value.Get()
		if fl.MaskFn != nil {
			val = fl.MaskFn(val)
		}
		out[fl.Name] = val
	}

	for _, group := range f.dynamicGroups() {
		for field, item := range group.Items() {
			di, ok := item.Value.(core.DynamicItemValues)
			if !ok {
				continue
			}
			for id, val := range di.ValuesAny() {
				if item.Flag != nil && item.Flag.MaskFn != nil {
					val = item.Flag.MaskFn(val)
				}
				key := group.Name() + "." + id + "." + field
				out[key] = val
			}
		}
	}

	return out
}

// DynamicGroups returns dynamic groups in registration order.
func (f *FlagSet) DynamicGroups() []*dynamic.Group {
	return f.dynamicGroupsOrder
}

// OrderedDynamicGroups returns dynamic groups sorted by name.
func (f *FlagSet) OrderedDynamicGroups() []*dynamic.Group {
	groups := make([]*dynamic.Group, len(f.dynamicGroupsOrder))
	copy(groups, f.dynamicGroupsOrder)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name() < groups[j].Name()
	})
	return groups
}

// --- AllOrNone Group Handling ---

// AllOrNoneGroups returns the registered all-or-none groups.
func (f *FlagSet) AllOrNoneGroups() []*core.AllOrNoneGroup { return f.allOrNoneGroup }

// AddAllOrNoneGroup registers an all-or-none group.
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

// OneOfGroups returns the registered one-of groups.
func (f *FlagSet) OneOfGroups() []*core.OneOfGroupGroup {
	return f.oneOfGroup
}

// AddOneOfGroup registers a one-of group.
func (f *FlagSet) AddOneOfGroup(name string, g *core.OneOfGroupGroup) {
	f.oneOfGroup = append(f.oneOfGroup, g)
}

// GetOneOfGroup returns or creates a one-of group.
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

// AttachToOneOfGroup attaches a flag to a one-of group.
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

// maybeAddBuiltinFlags registers built-in help and version flags lazily.
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
