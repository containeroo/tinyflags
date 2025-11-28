package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHideDefaultFromHelp(t *testing.T) {
	t.Parallel()

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
