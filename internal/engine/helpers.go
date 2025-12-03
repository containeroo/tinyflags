package engine

import "github.com/containeroo/tinyflags/internal/core"

// FirstChanged returns the value of the first changed flag in the given order.
// If no flag was changed, it returns defaultValue and false.
func FirstChanged[T any](defaultValue T, flags ...core.Flag[T]) (T, bool) {
	for _, f := range flags {
		if f == nil || !f.Changed() {
			continue
		}
		if v := f.Value(); v != nil {
			return *v, true
		}
		return defaultValue, true
	}
	return defaultValue, false
}
