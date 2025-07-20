package dynamic

import "fmt"

// BoolFlag represents a dynamic boolean flag (e.g. --svc.foo.enabled=true).
type BoolFlag struct {
	item       *DynamicScalarValue[bool] // value container for each instance ID
	strictMode bool                      // whether explicit true/false is required
}

// Strict enables strict mode (disables --flag shorthand).
func (b *BoolFlag) Strict() *BoolFlag {
	b.strictMode = true
	return b
}

// Get returns the parsed boolean for a given ID, or false if not set.
func (b *BoolFlag) Get(id string) (bool, bool) {
	return b.item.Get(id)
}

// MustGet returns the value for the ID or panics if not set.
func (b *BoolFlag) MustGet(id string) bool {
	val, ok := b.item.Get(id)
	if !ok {
		panic(fmt.Sprintf("value for id %q not set", id))
	}
	return val
}

// Values returns all boolean values by instance ID.
func (b *BoolFlag) Values() map[string]bool {
	return b.item.Values()
}

// ValuesAny returns values in a map[string]any form.
func (b *BoolFlag) ValuesAny() map[string]any {
	return b.item.ValuesAny()
}
