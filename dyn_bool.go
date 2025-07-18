package tinyflags

import "fmt"

// DynamicBoolFlag defines a builder for dynamic boolean flags (e.g. --http.alpha.debug=true).
type DynamicBoolFlag struct {
	fs     *FlagSet               // parent FlagSet
	bf     *baseFlag              // registration metadata
	item   *DynamicItemImpl[bool] // values by instance ID
	strict bool                   // requires explicit --flag=true/false
}

// Strict marks this dynamic bool flag as requiring an explicit value.
func (d *DynamicBoolFlag) Strict() *DynamicBoolFlag {
	d.strict = true
	return d
}

// Validator sets a custom validation function for the boolean value.
func (d *DynamicBoolFlag) Validator(fn func(bool) error) *DynamicBoolFlag {
	d.item.SetValidator(fn)
	return d
}

// Required marks the dynamic flag as required (for at least one instance).
func (d *DynamicBoolFlag) Required() *DynamicBoolFlag {
	d.bf.required = true
	return d
}

// Hidden hides the flag from help output.
func (d *DynamicBoolFlag) Hidden() *DynamicBoolFlag {
	d.bf.hidden = true
	return d
}

// Group adds the dynamic flag to a mutual exclusion group.
func (d *DynamicBoolFlag) Group(name string) *DynamicBoolFlag {
	if name == "" {
		return d
	}
	for _, g := range d.fs.groups {
		if g.name == name {
			g.flags = append(g.flags, d.bf)
			d.bf.group = g
			return d
		}
	}
	group := &mutualGroup{name: name, flags: []*baseFlag{d.bf}}
	d.fs.groups = append(d.fs.groups, group)
	d.bf.group = group
	return d
}

// Env manually sets the environment variable name for this flag.
func (d *DynamicBoolFlag) Env(key string) *DynamicBoolFlag {
	if d.bf.disableEnv {
		panic("cannot call Env after DisableEnv")
	}
	d.bf.envKey = key
	return d
}

// DisableEnv disables environment variable resolution for this flag.
func (d *DynamicBoolFlag) DisableEnv() *DynamicBoolFlag {
	if d.bf.envKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	d.bf.disableEnv = true
	return d
}

// Metavar sets the metavar used in help output (e.g. BOOL).
func (d *DynamicBoolFlag) Metavar(s string) *DynamicBoolFlag {
	d.bf.metavar = s
	return d
}

// Get returns the value for the given instance ID.
func (d *DynamicBoolFlag) Get(id string) (bool, bool) {
	return d.item.Get(id)
}

// MustGet returns the value or panics if not set.
func (d *DynamicBoolFlag) MustGet(id string) bool {
	val, ok := d.item.Get(id)
	if !ok {
		panic(fmt.Sprintf("value for id %q not set", id))
	}
	return val
}

// Values returns all instanceID → bool values.
func (d *DynamicBoolFlag) Values() map[string]bool {
	return d.item.Values()
}

// ValuesAny returns instanceID → any values for debugging.
func (d *DynamicBoolFlag) ValuesAny() map[string]any {
	return d.item.ValuesAny()
}
