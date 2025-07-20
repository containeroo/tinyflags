package engine

import (
	"strconv"
	"time"

	"github.com/containeroo/tinyflags/internal/scalar"
)

func (f *FlagSet) StringVarP(ptr *string, name, short string, def string, usage string) *scalar.ScalarFlag[string] {
	return defineScalar(f, ptr, name, short, usage, def, func(s string) (string, error) { return s, nil }, func(s string) string { return s })
}

func (f *FlagSet) DurationVarP(ptr *time.Duration, name, short string, def time.Duration, usage string) *scalar.ScalarFlag[time.Duration] {
	return defineScalar(f, ptr, name, short, usage, def, time.ParseDuration, time.Duration.String)
}

func (f *FlagSet) IntVarP(ptr *int, name, short string, def int, usage string) *scalar.ScalarFlag[int] {
	return defineScalar(f, ptr, name, short, usage, def, strconv.Atoi, strconv.Itoa)
}
