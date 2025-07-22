package dynamic

import (
	"strconv"
	"time"
)

// String adds a string field.
func (g *Group) String(field, usage string) *ScalarFlag[string] {
	return defineDynamicScalar(g,
		field,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)
}

// Int adds an int field.
func (g *Group) Int(field, usage string) *ScalarFlag[int] {
	return defineDynamicScalar(g, field, strconv.Atoi, strconv.Itoa)
}

// Duration add a time.Duration field.
func (g *Group) Duration(field, usage string) *ScalarFlag[time.Duration] {
	return defineDynamicScalar(g, field, time.ParseDuration, time.Duration.String)
}

// Float64 adds a float64 field.
func (g *Group) Float64(field, usage string) *ScalarFlag[float64] {
	return defineDynamicScalar(g, field,
		func(s string) (float64, error) { return strconv.ParseFloat(s, 64) },
		func(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) },
	)
}

// Float32 adds a float32 field.
func (g *Group) Float32(field, usage string) *ScalarFlag[float32] {
	return defineDynamicScalar(g, field,
		func(s string) (float32, error) {
			v, err := strconv.ParseFloat(s, 32)
			return float32(v), err
		},
		func(f float32) string {
			return strconv.FormatFloat(float64(f), 'f', -1, 32)
		},
	)
}
