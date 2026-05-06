package engine

import "github.com/containeroo/tinyflags/internal/core"

// resetParseState clears positional args and resets parse lifecycles.
func (f *FlagSet) resetParseState() {
	f.positional = nil
	f.visitParseLifecycles(func(lifecycle core.ParseLifecycle) {
		lifecycle.ResetParseState()
	})
}

// applyDefaultFinalizers finalizes untouched default values after parsing.
func (f *FlagSet) applyDefaultFinalizers() {
	f.visitParseLifecycles(func(lifecycle core.ParseLifecycle) {
		lifecycle.ApplyDefaultFinalize()
	})
}

// visitParseLifecycles visits every registered parse lifecycle value.
func (f *FlagSet) visitParseLifecycles(visit func(core.ParseLifecycle)) {
	for _, fl := range f.staticFlagsMap {
		lifecycle, ok := fl.Value.(core.ParseLifecycle)
		if ok {
			visit(lifecycle)
		}
	}

	for _, group := range f.dynamicGroups() {
		for _, item := range group.Items() {
			lifecycle, ok := item.Value.(core.ParseLifecycle)
			if ok {
				visit(lifecycle)
			}
		}
	}
}
