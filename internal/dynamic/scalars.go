package dynamic

import (
	"strconv"
	"time"
)

// String registers a dynamic string flag under this group.
func (g *Group) String(field, usage string) *ScalarFlag[string] {
	return defineDynamicScalar(g,
		field,
		func(s string) (string, error) { return s, nil },
		func(s string) string { return s },
	)
}

// Int registers a dynamic int flag under this group.
func (g *Group) Int(field, usage string) *ScalarFlag[int] {
	return defineDynamicScalar(g, field, strconv.Atoi, strconv.Itoa)
}

// Duration registers a dynamic time.Duration flag under this group.
func (g *Group) Duration(field, usage string) *ScalarFlag[time.Duration] {
	return defineDynamicScalar(g, field, time.ParseDuration, time.Duration.String)
}

func (g *Group) Float64(field, usage string) *ScalarFlag[float64] {
	return defineDynamicScalar(g, field,
		func(s string) (float64, error) { return strconv.ParseFloat(s, 64) },
		func(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) },
	)
}

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
