package tinyflags

import "time"

// DurationP defines a duration flag with the specified name, shorthand, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) DurationP(name, short string, def time.Duration, usage string) *Flag[time.Duration] {
	ptr := new(time.Duration)
	val := NewFlagItem(ptr, def, time.ParseDuration, time.Duration.String)
	return addScalar(f, name, short, usage, val, ptr)
}

// Duration defines a duration flag with the specified name, default value, and usage string.
// Returns the flag for chaining. Retrieve the parsed value with .Value().
func (f *FlagSet) Duration(name string, def time.Duration, usage string) *Flag[time.Duration] {
	return f.DurationP(name, "", def, usage)
}

// DurationVarP defines a duration flag with the specified name, shorthand, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) DurationVarP(ptr *time.Duration, name, short string, def time.Duration, usage string) *Flag[time.Duration] {
	val := NewFlagItem(ptr, def, time.ParseDuration, time.Duration.String)
	return addScalar(f, name, short, usage, val, ptr)
}

// DurationVar defines a duration flag with the specified name, default value, and usage string.
// The parsed value is stored in the provided pointer. Returns the flag for chaining.
// Retrieve the parsed value with .Value().
func (f *FlagSet) DurationVar(ptr *time.Duration, name string, def time.Duration, usage string) *Flag[time.Duration] {
	return f.DurationVarP(ptr, name, "", def, usage)
}
