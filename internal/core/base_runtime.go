package core

// EnvKeyLookup derives an environment key from a prefix and flag name.
type EnvKeyLookup = func(prefix, flagName string) string

// IsChanged reports whether the flag's underlying value has changed.
func (f *BaseFlag) IsChanged() bool {
	return f != nil && f.Value != nil && f.Value.Changed()
}

// MissingRequired reports whether a required flag is currently unset.
func (f *BaseFlag) MissingRequired() bool {
	return f != nil && f.Required && !f.IsChanged()
}

// LookupEnvKey returns the environment key to check and whether env loading applies.
func (f *BaseFlag) LookupEnvKey(prefix string, lookup EnvKeyLookup) (string, bool) {
	if f == nil || f.Value == nil || f.DisableEnv || f.IsChanged() {
		return "", false
	}

	envKey := f.EnvKey
	if envKey == "" && lookup != nil {
		envKey = lookup(prefix, f.Name)
	}
	if envKey == "" {
		return "", false
	}
	return envKey, true
}

// FirstMissingRequirement returns the first required flag name that is not satisfied.
func (f *BaseFlag) FirstMissingRequirement(flags map[string]*BaseFlag) (string, bool) {
	if f == nil || !f.IsChanged() {
		return "", false
	}
	for _, req := range f.Requires {
		required, ok := flags[req]
		if !ok || !required.IsChanged() {
			return req, true
		}
	}
	return "", false
}
