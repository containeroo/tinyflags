package tinyflags_test

import (
	"errors"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestExportedUsageErrorHelpers(t *testing.T) {
	t.Parallel()

	cmdErr := &tinyflags.UsageError{
		Err:  &tinyflags.CommandRequired{Command: "app"},
		Help: "Usage: app <command>\n",
	}

	require.True(t, tinyflags.IsCommandRequired(cmdErr))
	require.False(t, tinyflags.IsHelpRequested(cmdErr))
	require.False(t, tinyflags.IsVersionRequested(cmdErr))
	help, ok := tinyflags.HelpText(cmdErr)
	require.True(t, ok)
	assert.Equal(t, "Usage: app <command>\n", help)

	var typed *tinyflags.CommandRequired
	assert.True(t, errors.As(cmdErr, &typed))
	require.NotNil(t, typed)
	assert.Equal(t, "app", typed.Command)

	var usage *tinyflags.UsageError
	assert.True(t, errors.As(cmdErr, &usage))
	require.NotNil(t, usage)
	assert.Equal(t, "Usage: app <command>\n", usage.Help)
}
