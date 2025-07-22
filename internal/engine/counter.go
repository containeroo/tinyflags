package engine

import "github.com/containeroo/tinyflags/internal/scalar"

func (f *FlagSet) CounterVar(ptr *int, name, usage string, def int) *scalar.CounterFlag {
	return scalar.NewCounter(f, ptr, name, usage, def)
}
