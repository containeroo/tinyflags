package tinyflags_test

import (
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelpBuiltinSections(t *testing.T) {
	t.Parallel()

	t.Run("unnamedSectionRendersLastByDefault", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.2.3")
		fs.String("plain", "", "plain flag")
		fs.String("server", "", "server flag").Section("Server")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		require.Contains(t, out, "Server:")
		require.Contains(t, out, "--plain")
		require.Contains(t, out, "--help")
		require.Contains(t, out, "--version")
		assert.NotContains(t, out, "\n:\n")
		assert.Less(t, strings.Index(out, "Server:"), strings.Index(out, "--plain"))
		assert.Less(t, strings.Index(out, "--plain"), strings.Index(out, "--help"))
		assert.Less(t, strings.Index(out, "--help"), strings.Index(out, "--version"))
	})

	t.Run("sectionOrderCanPositionUnnamedSectionFirst", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.2.3")
		fs.String("plain", "", "plain flag")
		fs.String("server", "", "server flag").Section("Server")
		fs.SectionOrder("", "Server")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		require.Contains(t, out, "Server:")
		assert.Less(t, strings.Index(out, "--plain"), strings.Index(out, "Server:"))
		assert.Less(t, strings.Index(out, "--help"), strings.Index(out, "Server:"))
		assert.Less(t, strings.Index(out, "--version"), strings.Index(out, "Server:"))
	})

	t.Run("builtinSectionGroupsHelpAndVersion", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.2.3")
		fs.String("plain", "", "plain flag")
		fs.String("server", "", "server flag").Section("Server")
		fs.BuiltinSection("General")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		assert.Equal(t, 1, strings.Count(out, "General:"))
		assert.Less(t, strings.Index(out, "Server:"), strings.Index(out, "General:"))
		assert.Less(t, strings.Index(out, "General:"), strings.Index(out, "--help"))
		assert.Less(t, strings.Index(out, "--help"), strings.Index(out, "--version"))
		assert.Less(t, strings.Index(out, "--version"), strings.Index(out, "--plain"))
	})

	t.Run("helpOptionsCanSetBuiltinSection", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.2.3")
		fs.String("server", "", "server flag").Section("Server")
		fs.Help().BuiltinSection("CLI")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		assert.Equal(t, 1, strings.Count(out, "CLI:"))
		assert.Contains(t, out, "--help")
		assert.Contains(t, out, "--version")
	})

	t.Run("builtinSectionCanChangeAfterInitialHelpRender", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.2.3")
		fs.String("server", "", "server flag").Section("Server")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "General:")

		fs.Help().BuiltinSection("General")

		err = fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		assert.Equal(t, 1, strings.Count(out, "General:"))
		assert.Less(t, strings.Index(out, "General:"), strings.Index(out, "--help"))
		assert.Less(t, strings.Index(out, "--help"), strings.Index(out, "--version"))
	})

	t.Run("emptyBuiltinSectionMovesBuiltinsBackToUnnamed", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.2.3")
		fs.String("server", "", "server flag").Section("Server")
		fs.BuiltinSection("General")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "General:")

		fs.BuiltinSection("")

		err = fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		assert.NotContains(t, out, "General:")
		assert.Less(t, strings.Index(out, "Server:"), strings.Index(out, "--help"))
		assert.Less(t, strings.Index(out, "--help"), strings.Index(out, "--version"))
	})
}
