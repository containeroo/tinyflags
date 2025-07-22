package tinyflags

import "github.com/containeroo/tinyflags/internal/scalar"

func (f *FlagSet) Counter(name, usage string) *scalar.CounterFlag {
	return f.impl.CounterVar(new(int), name, usage, 0)
}

func (f *FlagSet) CounterVar(ptr *int, name, usage string) *scalar.CounterFlag {
	return f.impl.CounterVar(ptr, name, usage, 0)
}
