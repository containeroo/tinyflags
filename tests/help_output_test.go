package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHideDefaultFromHelp verifies hidden defaults in help output.
func TestHideDefaultFromHelp(t *testing.T) {
	t.Parallel()

	t.Run("counterDefaultHidden", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Counter("verbose", 0, "Enable verbose mode")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "Default:")
	})

	t.Run("staticHideDefault", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("name", "alice", "User name").HideDefault()

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "Default: alice")
	})

	t.Run("dynamicHideDefault", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		g := fs.DynamicGroup("g")
		g.String("mode", "prod", "mode").HideDefault()

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "Default: prod")
	})
}

// TestHelpSections verifies help section headings.
func TestHelpSections(t *testing.T) {
	t.Parallel()

	t.Run("sectionHeadersPrinted", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("general", "", "general flag").Section("General")
		fs.String("net", "", "net flag").Section("Network")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()
		assert.Contains(t, out, "General:")
		assert.Contains(t, out, "Network:")
	})
}

// TestDynamicGroupFooter verifies dynamic group footer notes in help output.
func TestDynamicGroupFooter(t *testing.T) {
	t.Parallel()

	t.Run("notePrinted", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.SetDynamicUsageIndent(8)
		g := fs.DynamicGroup("srv")
		g.Note("footer note for group")
		g.String("addr", "", "address")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()
		assert.Contains(t, out, "footer")
		assert.Contains(t, out, "note")
		assert.Contains(t, out, "group")
	})
}

// TestNotesPreserveNewlines verifies help notes preserve explicit newlines.
func TestNotesPreserveNewlines(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.Note("first line\nsecond line")

	err := fs.Parse([]string{"--help"})
	require.Error(t, err)
	out := err.Error()
	assert.Contains(t, out, "first line\nsecond line")
}

// TestCounterPlaceholderHidden verifies counters do not render value placeholders.
func TestCounterPlaceholderHidden(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.Counter("verbose", 0, "Enable verbose logging").
		Short("v")

	err := fs.Parse([]string{"--help"})
	require.Error(t, err)
	out := err.Error()
	assert.Contains(t, out, "-v, --verbose")
	assert.NotContains(t, out, "VERBOSE")
}
