package engine

import "github.com/containeroo/tinyflags/internal/core"

// ValueSource identifies where the effective value for a flag came from.
type ValueSource string

const (
	// ValueSourceDefault means the flag kept its default value.
	ValueSourceDefault ValueSource = "Default"
	// ValueSourceFlag means the value was provided by a command-line flag.
	ValueSourceFlag ValueSource = "Flag"
	// ValueSourceEnvironment means the value was loaded from an environment variable.
	ValueSourceEnvironment ValueSource = "Environment"
)

// Source returns the source of name from the most recent parse.
// Names use the same format as OverriddenValues, including group.id.flag for dynamic values.
func (f *FlagSet) Source(name string) ValueSource {
	if source, ok := f.valueSources[name]; ok {
		return source
	}
	return ValueSourceDefault
}

// OverriddenSources returns the sources of all values explicitly set by flags or environment variables.
// Dynamic flags use the key format group.id.flag.
func (f *FlagSet) OverriddenSources() map[string]ValueSource {
	out := make(map[string]ValueSource, len(f.valueSources))
	for name, source := range f.valueSources {
		out[name] = source
	}
	return out
}

// captureChangedSources records the source for changed values that do not already have one.
// Parsing command-line flags happens before environment loading, so the first source wins.
func (f *FlagSet) captureChangedSources(source ValueSource) {
	if f.valueSources == nil {
		f.valueSources = make(map[string]ValueSource)
	}

	for _, fl := range f.staticFlagsMap {
		if fl.Value == nil || !fl.Value.Changed() {
			continue
		}
		if _, recorded := f.valueSources[fl.Name]; !recorded {
			f.valueSources[fl.Name] = source
		}
	}

	for _, group := range f.dynamicGroups() {
		for field, item := range group.Items() {
			values, ok := item.Value.(core.DynamicItemValues)
			if !ok {
				continue
			}
			for id := range values.ValuesAny() {
				name := group.Name() + "." + id + "." + field
				if _, recorded := f.valueSources[name]; !recorded {
					f.valueSources[name] = source
				}
			}
		}
	}
}
