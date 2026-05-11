package tinyflags_test

import (
	"errors"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExportedCompatibilityAliases verifies exported compatibility aliases.
func TestExportedCompatibilityAliases(t *testing.T) {
	t.Parallel()

	t.Run("globalDelimiterAliasMatchesPreferredMethod", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Globaldelimiter("|")
		val := fs.StringSlice("tag", nil, "tags").Value()

		err := fs.Parse([]string{"--tag=a|b"})
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, *val)
	})

	t.Run("allOrNonePluralAndSingularAccessorsMatch", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.GetAllOrNoneGroup("auth")

		assert.Len(t, fs.AllOrNoneGroup(), 1)
		assert.Len(t, fs.AllOrNoneGroups(), 1)
		assert.Same(t, fs.AllOrNoneGroup()[0], fs.AllOrNoneGroups()[0])
	})
}

// TestExportedHelpVersionSentinels verifies exported help/version sentinels.
func TestExportedHelpVersionSentinels(t *testing.T) {
	t.Parallel()

	helpErr := tinyflags.RequestHelp("show help")
	versionErr := tinyflags.RequestVersion("1.2.3")

	require.True(t, tinyflags.IsHelpRequested(helpErr))
	require.True(t, tinyflags.IsVersionRequested(versionErr))
	require.False(t, tinyflags.IsHelpRequested(versionErr))
	require.False(t, tinyflags.IsVersionRequested(helpErr))

	var helpTyped *tinyflags.HelpRequested
	var versionTyped *tinyflags.VersionRequested
	assert.True(t, errors.As(helpErr, &helpTyped))
	assert.True(t, errors.As(versionErr, &versionTyped))
}

func TestExportedCommandRequiredSentinel(t *testing.T) {
	t.Parallel()

	cmdErr := tinyflags.RequestCommandRequired("app")

	require.True(t, tinyflags.IsCommandRequired(cmdErr))
	require.False(t, tinyflags.IsHelpRequested(cmdErr))
	require.False(t, tinyflags.IsVersionRequested(cmdErr))

	var typed *tinyflags.CommandRequired
	assert.True(t, errors.As(cmdErr, &typed))
	require.NotNil(t, typed)
	assert.Equal(t, "app", typed.Command)
}
