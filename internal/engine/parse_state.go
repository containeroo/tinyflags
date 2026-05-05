package engine

import "github.com/containeroo/tinyflags/internal/core"

func (f *FlagSet) resetParseState() {
	f.positional = nil

	for _, fl := range f.staticFlagsMap {
		resetter, ok := fl.Value.(core.ParseStateResetter)
		if ok {
			resetter.ResetParseState()
		}
	}

	for _, group := range f.dynamicGroups() {
		for _, item := range group.Items() {
			resetter, ok := item.Value.(core.ParseStateResetter)
			if ok {
				resetter.ResetParseState()
			}
		}
	}
}

func (f *FlagSet) applyDefaultFinalizers() {
	for _, fl := range f.staticFlagsMap {
		if fl.Value == nil {
			continue
		}
		if finalizer, ok := fl.Value.(core.DefaultFinalizer); ok {
			finalizer.ApplyDefaultFinalize()
		}
	}

	for _, group := range f.dynamicGroups() {
		for _, item := range group.Items() {
			if finalizer, ok := item.Value.(core.DefaultFinalizer); ok {
				finalizer.ApplyDefaultFinalize()
			}
		}
	}
}
