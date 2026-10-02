package tinyflags_test

import (
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelpSectionRegressionCases(t *testing.T) {
	t.Parallel()

	t.Run("unlistedSectionsFollowExplicitOrder", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("database", "", "database").Section("Database")
		fs.String("server", "", "server").Section("Server")
		fs.String("services", "", "services").Section("Services")
		fs.SectionOrder("Server")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		require.Contains(t, out, "Server:")
		require.Contains(t, out, "Database:")
		require.Contains(t, out, "Services:")
		assert.Less(t, strings.Index(out, "Server:"), strings.Index(out, "Database:"))
		assert.Less(t, strings.Index(out, "Database:"), strings.Index(out, "Services:"))
	})

	t.Run("unknownExplicitSectionIsIgnored", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("server", "", "server").Section("Server")
		fs.SectionOrder("Missing", "Server")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		assert.NotContains(t, out, "Missing:")
		assert.Equal(t, 1, strings.Count(out, "Server:"))
	})

	t.Run("duplicateExplicitSectionsRenderOnce", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("server", "", "server").Section("Server")
		fs.String("database", "", "database").Section("Database")
		fs.SectionOrder("Server", "Server", "Database", "Server")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		assert.Equal(t, 1, strings.Count(out, "Server:"))
		assert.Equal(t, 1, strings.Count(out, "Database:"))
	})

	t.Run("hiddenOnlySectionIsOmitted", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("secret", "", "secret").Section("Hidden").Hidden()
		fs.String("server", "", "server").Section("Server")
		fs.SectionOrder("Hidden", "Server")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		assert.NotContains(t, out, "Hidden:")
		assert.NotContains(t, out, "--secret")
		assert.Contains(t, out, "Server:")
	})

	t.Run("registrationOrderIsPreservedWithinSection", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("third", "", "third").Section("Server")
		fs.String("first", "", "first").Section("Database")
		fs.String("second", "", "second").Section("Server")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		require.Contains(t, out, "--third")
		require.Contains(t, out, "--second")
		assert.Less(t, strings.Index(out, "--third"), strings.Index(out, "--second"))
	})

	t.Run("noSectionsKeepsFlatRegistrationOrder", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("first", "", "first")
		fs.String("second", "", "second")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		assert.NotContains(t, out, "first:")
		assert.Less(t, strings.Index(out, "--first"), strings.Index(out, "--second"))
		assert.Less(t, strings.Index(out, "--second"), strings.Index(out, "--help"))
	})

	t.Run("helpOptionsCanSetSectionOrder", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("server", "", "server").Section("Server")
		fs.String("database", "", "database").Section("Database")
		fs.Help().SectionOrder("Database", "Server")

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		out := err.Error()

		assert.Less(t, strings.Index(out, "Database:"), strings.Index(out, "Server:"))
	})
}
