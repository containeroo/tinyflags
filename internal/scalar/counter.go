package scalar

import (
	"fmt"
	"strconv"

	"github.com/containeroo/tinyflags/internal/core"
)

func (c *CounterValue) Increment() error {
	next := *c.ptr + 1
	if c.validate != nil {
		if err := c.validate(next); err != nil {
			return err
		}
	}
	*c.ptr = next
	c.value = next
	c.changed = true
	return nil
}

// CounterValue is a scalar int that increments on each occurrence.
type CounterValue struct {
	*ScalarValue[int]
}

// NewCounterValue returns a new counter that increments on each Set().
func NewCounterValue(ptr *int, def int) *CounterValue {
	return &CounterValue{
		ScalarValue: NewScalarValue(
			ptr,
			def,
			func(s string) (int, error) {
				// Allow explicit --flag=N
				if s == "" {
					return *ptr + 1, nil
				}
				n, err := strconv.Atoi(s)
				if err != nil {
					return 0, fmt.Errorf("invalid counter value: %w", err)
				}
				return n, nil
			},
			strconv.Itoa,
		),
	}
}

// Set overrides the default Set to increment on empty string.
func (c *CounterValue) Set(s string) error {
	if s == "" {
		*c.ptr++
		c.value = *c.ptr
		c.changed = true
		return nil
	}
	return c.ScalarValue.Set(s)
}

// CounterFlag provides fluent builder methods for counter flags.
type CounterFlag struct {
	*ScalarFlag[int]
	val *CounterValue
}

// NewCounter creates a new counter flag.
func NewCounter(r core.Registry, ptr *int, name, usage string, def int) *CounterFlag {
	val := NewCounterValue(ptr, def)
	flag := RegisterScalar(r, name, usage, val, ptr)
	return &CounterFlag{ScalarFlag: flag, val: val}
}

// Max sets a maximum allowed value for the counter.
func (c *CounterFlag) Max(n int) *CounterFlag {
	c.val.setValidate(func(v int) error {
		if v > n {
			return fmt.Errorf("must not be greater than %d", n)
		}
		return nil
	})
	return c
}
