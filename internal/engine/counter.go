package engine

import "github.com/containeroo/tinyflags/internal/scalar"

func (f *FlagSet) CounterVar(ptr *int, name string, def int, usage string) *scalar.CounterFlag {
	return scalar.NewCounter(f, ptr, name, def, usage)
}
