package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/require"
)

// TestFlagNotBlank rejects explicitly supplied blank strings without rejecting defaults.
func TestFlagNotBlank(t *testing.T) {
	t.Parallel()

	t.Run("static", func(t *testing.T) {
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("name", "", "Application name").NotBlank()

		require.NoError(t, fs.Parse(nil), "unset defaults remain valid")
		require.NoError(t, fs.Parse([]string{"--name= service "}))
		require.EqualError(t, fs.Parse([]string{"--name= \t\n"}), "flag --name must not be blank")
	})

	t.Run("dynamic", func(t *testing.T) {
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.DynamicGroup("server").String("host", "", "Server host").NotBlank()

		require.NoError(t, fs.Parse([]string{"--server.alpha.host=api.example"}))
		require.EqualError(t, fs.Parse([]string{"--server.alpha.host=  "}), "flag --server.alpha.host must not be blank")
	})
}
