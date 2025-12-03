package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilderChainingPreservesFlagType(t *testing.T) {
	t.Parallel()

	t.Run("boolShortKeepsChanged", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		debug := fs.Bool("debug", false, "enable debug logs").Short("d")

		err := fs.Parse([]string{"--debug"})
		require.NoError(t, err)
		assert.True(t, debug.Changed())
		assert.True(t, *debug.Value())
	})

	t.Run("scalarShortThenChoices", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		name := fs.String("name", "", "user name").Short("n").Choices("alice", "bob").Value()

		err := fs.Parse([]string{"--name=bob"})
		require.NoError(t, err)
		assert.Equal(t, "bob", *name)
	})

	t.Run("counterShortThenMax", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		verbose := fs.Counter("verbose", 0, "verbosity").Short("v").Max(2).Value()

		err := fs.Parse([]string{"-v", "-v"})
		require.NoError(t, err)
		assert.Equal(t, 2, *verbose)
	})

	t.Run("sliceShortThenChoices", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		hosts := fs.StringSlice("host", []string{}, "hosts").Short("H").Choices("a", "b").Value()

		err := fs.Parse([]string{"--host=a", "--host=b"})
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, *hosts)
	})
}
