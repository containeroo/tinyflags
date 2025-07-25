package scalar

import (
	"strconv"

	"github.com/containeroo/tinyflags/internal/core"
)

// BoolValue holds the internal state of a boolean flag and whether it is strict.
type BoolValue struct {
	*ScalarValue[bool]
	Strict *bool
}

// Strict marks this boolean flag as requiring an explicit value.
func (b *BoolFlag) Strict() *BoolFlag {
	*b.val.Strict = true
	return b
}

// IsStrictBool reports whether the flag requires an explicit value (--flag=true/false).
func (b *BoolValue) IsStrictBool() bool {
	return b.Strict != nil && *b.Strict
}

// BoolFlag provides fluent builder methods for boolean flags,
// including support for .Strict() to require explicit values.
type BoolFlag struct {
	*ScalarFlag[bool] // embeds core builder methods like Env(), Required(), etc.
	val               *BoolValue
}

// NewBoolValue returns a BoolValue with parse/format logic and default value.
func NewBoolValue(ptr *bool, def bool) *BoolValue {
	strict := new(bool)
	return &BoolValue{
		ScalarValue: NewScalarValue(
			ptr,
			def,
			strconv.ParseBool,
			strconv.FormatBool,
		),
		Strict: strict,
	}
}

// NewBool creates a new BoolFlag with full builder support.
func NewBool(r core.Registry, ptr *bool, name string, def bool, usage string) *BoolFlag {
	val := NewBoolValue(ptr, def)

	flag := RegisterScalar(r, name, usage, val, ptr)
	return &BoolFlag{ScalarFlag: flag, val: val}
}
