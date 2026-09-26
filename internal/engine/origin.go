package engine

// ValueSource identifies the kind of input that supplied an effective flag value.
type ValueSource string

const (
	// ValueSourceDefault means the flag kept its default value.
	ValueSourceDefault ValueSource = "Default"
	// ValueSourceFlag means a command-line flag supplied the value.
	ValueSourceFlag ValueSource = "Flag"
	// ValueSourceEnvironment means an environment variable supplied the value.
	ValueSourceEnvironment ValueSource = "Environment"
)

// ValueOrigin identifies the exact input that supplied an effective flag value.
// Key is the flag spelling (for example "-a" or "--address") or environment
// variable name that actually won during parsing. Defaults have an empty Key.
type ValueOrigin struct {
	Source ValueSource
	Key    string
}

// String returns a human-readable origin label suitable for diagnostics and UIs.
func (o ValueOrigin) String() string {
	if o.Source == "" || o.Source == ValueSourceDefault {
		return string(ValueSourceDefault)
	}
	if o.Key == "" {
		return string(o.Source)
	}
	return string(o.Source) + " · " + o.Key
}

// Origin returns the exact origin of name from the most recent parse.
// Names use the same format as OverriddenValues, including group.id.flag for dynamic values.
// Unoverridden values report the default origin.
func (f *FlagSet) Origin(name string) ValueOrigin {
	if origin, ok := f.valueOrigins[name]; ok {
		return origin
	}
	return ValueOrigin{Source: ValueSourceDefault}
}

// OverriddenOrigins returns the exact origin of every value explicitly set by a flag or environment variable.
// Dynamic flags use the key format group.id.flag.
func (f *FlagSet) OverriddenOrigins() map[string]ValueOrigin {
	out := make(map[string]ValueOrigin, len(f.valueOrigins))
	for name, origin := range f.valueOrigins {
		out[name] = origin
	}
	return out
}

// recordOrigin stores the latest successful input that contributed the effective value.
func (f *FlagSet) recordOrigin(name string, source ValueSource, key string) {
	if f.valueOrigins == nil {
		f.valueOrigins = make(map[string]ValueOrigin)
	}
	f.valueOrigins[name] = ValueOrigin{Source: source, Key: key}
}
