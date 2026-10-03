package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/require"
)

// TestFlagRegistrationRejectsDuplicates verifies registration cannot create ambiguous lookups.
func TestFlagRegistrationRejectsDuplicates(t *testing.T) {
	t.Parallel()

	t.Run("long name", func(t *testing.T) {
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("name", "", "first")

		require.PanicsWithValue(t, "tinyflags: duplicate flag --name", func() {
			fs.String("name", "", "second")
		})
	})

	t.Run("short alias", func(t *testing.T) {
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("name", "", "first").Short("n")

		require.PanicsWithValue(t, "tinyflags: duplicate short flag -n", func() {
			fs.String("number", "", "second").Short("n")
		})
	})

	t.Run("dynamic field", func(t *testing.T) {
		group := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError).DynamicGroup("server")
		group.String("host", "", "first")

		require.PanicsWithValue(t, "tinyflags: duplicate dynamic field host in group server", func() {
			group.String("host", "", "second")
		})
	})
}
