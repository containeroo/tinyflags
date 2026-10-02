package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandRoutingRegressionCases(t *testing.T) {
	t.Parallel()

	t.Run("combinedShortFlagsCanBelongToDifferentScopes", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
		verbose := root.Globals().Bool("verbose", false, "verbose").Short("v").Value()
		serve := root.Command("serve", "serve")
		debug := serve.Bool("debug", false, "debug").Short("d").Value()

		require.NoError(t, root.Parse([]string{"serve", "-vd"}))
		assert.True(t, *verbose)
		assert.True(t, *debug)
	})

	t.Run("localLongFlagShadowsAncestorGlobal", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
		globalMode := root.Globals().String("mode", "global", "global mode").Value()
		serve := root.Command("serve", "serve")
		localMode := serve.String("mode", "local", "local mode").Value()

		require.NoError(t, root.Parse([]string{"serve", "--mode=child"}))
		assert.Equal(t, "global", *globalMode)
		assert.Equal(t, "child", *localMode)
	})

	t.Run("localShortFlagShadowsAncestorGlobal", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
		globalMode := root.Globals().Bool("global-mode", false, "global").Short("m").Value()
		serve := root.Command("serve", "serve")
		localMode := serve.Bool("local-mode", false, "local").Short("m").Value()

		require.NoError(t, root.Parse([]string{"serve", "-m"}))
		assert.False(t, *globalMode)
		assert.True(t, *localMode)
	})

	t.Run("flagValueMatchingSubcommandNameIsConsumedAsValue", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
		mode := root.Globals().String("mode", "default", "mode").Value()
		serve := root.Command("serve", "serve")
		port := serve.Int("port", 8080, "port").Value()

		require.NoError(t, root.Parse([]string{"--mode", "serve"}))
		assert.Equal(t, "serve", *mode)
		assert.Equal(t, "app", root.SelectedCommand().FullName())
		assert.Equal(t, 8080, *port)
	})

	t.Run("doubleDashPreservesFlagLikePositionals", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
		serve := root.Command("serve", "serve")

		require.NoError(t, root.Parse([]string{"serve", "--", "--literal", "tail"}))
		assert.Equal(t, "app serve", root.SelectedCommand().FullName())
		assert.Equal(t, []string{"--literal", "tail"}, serve.Args())
	})

	t.Run("selectedCommandResetsAcrossParses", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
		root.Command("serve", "serve")

		require.NoError(t, root.Parse([]string{"serve"}))
		assert.Equal(t, "app serve", root.SelectedCommand().FullName())

		require.NoError(t, root.Parse(nil))
		assert.Equal(t, "app", root.SelectedCommand().FullName())
	})

	t.Run("nestedHelpTargetsCurrentCommand", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
		admin := root.Command("admin", "admin")
		users := admin.Command("users", "users")
		users.String("name", "", "name")

		err := root.Parse([]string{"admin", "users", "--help"})
		require.Error(t, err)
		require.True(t, tinyflags.IsHelpRequested(err))
		assert.Contains(t, err.Error(), "Usage: app admin users")
		assert.Contains(t, err.Error(), "--name NAME")
		assert.NotContains(t, err.Error(), "Commands:")
	})
}
