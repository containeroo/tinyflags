package tinyflags

import "fmt"

// DynamicFlag is a builder for a dynamic scalar flag.
// Each instance (e.g. `--http.alpha.port`) can hold one value.
type DynamicFlag[T any] struct {
	fs   *FlagSet            // parent flag set
	bf   *baseFlag           // metadata and registration
	item *DynamicItemImpl[T] // actual value store
}

// Get retrieves the value for the given instance ID.
func (d *DynamicFlag[T]) Get(id string) (T, bool) {
	return d.item.Get(id)
}

// MustGet retrieves the value or panics if not set.
func (d *DynamicFlag[T]) MustGet(id string) T {
	val, ok := d.item.Get(id)
	if !ok {
		panic(fmt.Sprintf("value for id %q not set", id))
	}
	return val
}

// Values returns a map of all instance IDs and their values.
func (d *DynamicFlag[T]) Values() map[string]T {
	return d.item.Values()
}

// ValuesAny returns a map[string]any for introspection/debug.
func (d *DynamicFlag[T]) ValuesAny() map[string]any {
	return d.item.ValuesAny()
}

func (d *DynamicFlag[T]) Required() *DynamicFlag[T] {
	d.bf.required = true
	return d
}

func (d *DynamicFlag[T]) Hidden() *DynamicFlag[T] {
	d.bf.hidden = true
	return d
}

func (d *DynamicFlag[T]) Group(name string) *DynamicFlag[T] {
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

func (d *DynamicFlag[T]) Env(key string) *DynamicFlag[T] {
	if d.bf.disableEnv {
		panic("cannot call Env after DisableEnv")
	}
	d.bf.envKey = key
	return d
}

func (d *DynamicFlag[T]) DisableEnv() *DynamicFlag[T] {
	if d.bf.envKey != "" {
		panic("cannot call DisableEnv after Env")
	}
	d.bf.disableEnv = true
	return d
}

func (d *DynamicFlag[T]) Metavar(s string) *DynamicFlag[T] {
	d.bf.metavar = s
	return d
}

func (d *DynamicFlag[T]) Validator(fn func(T) error) *DynamicFlag[T] {
	d.item.SetValidator(fn)
	return d
}
