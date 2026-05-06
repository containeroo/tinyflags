package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelpAndVersionErrors(t *testing.T) {
	t.Parallel()

	t.Run("helpSentinel", func(t *testing.T) {
		t.Parallel()

		err := RequestHelp("help me")
		require.True(t, IsHelpRequested(err))
		require.EqualError(t, err, "help me")
	})

	t.Run("versionSentinel", func(t *testing.T) {
		t.Parallel()

		err := RequestVersion("v1.2.3")
		require.True(t, IsVersionRequested(err))
		require.EqualError(t, err, "v1.2.3")
	})
}
