package tinyflags

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagItem(t *testing.T) {
	t.Parallel()

	t.Run("NewFlagItem sets default", func(t *testing.T) {
		t.Parallel()
		var x int
		v := NewFlagItem(&x, 42, parseInt, formatInt)

		assert.Equal(t, 42, x)
		assert.False(t, v.changed)
		assert.Equal(t, "42", v.Default())
	})

	t.Run("Set applies parsed value", func(t *testing.T) {
		t.Parallel()
		var x int
		v := NewFlagItem(&x, 10, parseInt, formatInt)

		err := v.Set("123")
		assert.NoError(t, err)
		assert.Equal(t, 123, x)
		assert.True(t, v.changed)
	})

	t.Run("Set fails on parse error", func(t *testing.T) {
		t.Parallel()
		var x int
		v := NewFlagItem(&x, 0, parseInt, formatInt)

		err := v.Set("notanint")
		assert.Error(t, err)
		assert.Equal(t, 0, x)
		assert.False(t, v.changed)
	})

	t.Run("Set fails on validator error", func(t *testing.T) {
		t.Parallel()
		var x int
		v := NewFlagItem(&x, 0, parseInt, formatInt)

		v.SetValidator(func(n int) error {
			if n < 0 {
				return errors.New("must be positive")
			}
			return nil
		}, []int{-1, 0, 1})

		err := v.Set("-3")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid value")
		assert.Contains(t, err.Error(), "must be positive")
		assert.Equal(t, 0, x)
		assert.False(t, v.changed)
	})

	t.Run("SetValidator replaces previous validator", func(t *testing.T) {
		t.Parallel()
		var x int
		v := NewFlagItem(&x, 0, parseInt, formatInt)

		v.SetValidator(func(int) error { return errors.New("fail") }, []int{1})
		v.SetValidator(nil, nil) // should unset validator

		err := v.Set("9")
		assert.NoError(t, err)
		assert.Equal(t, 9, x)
		assert.True(t, v.changed)
	})

	t.Run("Get returns value", func(t *testing.T) {
		t.Parallel()
		var x int
		v := NewFlagItem(&x, 99, parseInt, formatInt)
		assert.Equal(t, 99, v.Get())
	})

	t.Run("IsChanged true after Set", func(t *testing.T) {
		t.Parallel()
		var x int
		v := NewFlagItem(&x, 1, parseInt, formatInt)
		assert.False(t, v.IsChanged())

		_ = v.Set("2")
		assert.True(t, v.IsChanged())
	})
}

// Helpers for int parsing/formatting
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func formatInt(i int) string {
	return fmt.Sprintf("%d", i)
}
