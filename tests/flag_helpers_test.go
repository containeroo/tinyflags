package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirstChanged(t *testing.T) {
	t.Parallel()

	t.Run("returnsFirstChanged", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		a := fs.Bool("a", false, "flag a")
		b := fs.Bool("b", false, "flag b")

		err := fs.Parse([]string{"--b=true"})
		require.NoError(t, err)

		val, ok := tinyflags.FirstChanged(false, a, b)
		assert.True(t, ok)
		assert.True(t, val)
	})

	t.Run("returnsDefaultWhenNoneChanged", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		a := fs.Bool("a", false, "flag a")
		b := fs.Bool("b", true, "flag b")

		err := fs.Parse([]string{})
		require.NoError(t, err)

		val, ok := tinyflags.FirstChanged(false, a, b)
		assert.False(t, ok)
		assert.False(t, val)
	})
}
