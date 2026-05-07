package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommandGlobalsBeforeAndAfter verifies persistent flags parse around subcommands.
func TestCommandGlobalsBeforeAndAfter(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	verbose := root.Globals().Bool("verbose", false, "verbose").Value()

	serve := root.Command("serve", "Run the server")
	port := serve.Int("port", 8080, "port").Value()

	err := root.Parse([]string{"--verbose", "serve", "--port=9000"})
	require.NoError(t, err)
	assert.True(t, *verbose)
	assert.Equal(t, 9000, *port)
	assert.Equal(t, "serve", root.SelectedCommand().Name())

	err = root.Parse([]string{"serve", "--port=7000", "--verbose"})
	require.NoError(t, err)
	assert.True(t, *verbose)
	assert.Equal(t, 7000, *port)
}

// TestNestedCommands verifies nested subcommand selection.
func TestNestedCommands(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	root.Globals().Bool("verbose", false, "verbose").Value()

	admin := root.Command("admin", "Admin tools")
	admin.Globals().Bool("audit", false, "audit").Value()

	users := admin.Command("users", "Manage users")
	name := users.String("name", "", "user name").Value()

	err := root.Parse([]string{"admin", "--audit", "users", "--name=alice", "--verbose"})
	require.NoError(t, err)
	require.NotNil(t, root.SelectedCommand())
	assert.Equal(t, "app admin users", root.SelectedCommand().FullName())
	assert.Equal(t, "alice", *name)
}

// TestSubcommandHelpHidesGlobals verifies subcommand help omits inherited globals.
func TestSubcommandHelpHidesGlobals(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	root.Globals().Bool("verbose", false, "verbose")

	serve := root.Command("serve", "Run the server")
	serve.Int("port", 8080, "port")

	err := root.Parse([]string{"serve", "--help"})
	require.Error(t, err)
	require.True(t, tinyflags.IsHelpRequested(err))
	assert.Contains(t, err.Error(), "--port PORT")
	assert.NotContains(t, err.Error(), "--verbose")
}

// TestCommandHelpListsChildren verifies command help includes child listings.
func TestCommandHelpListsChildren(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	root.Globals().Bool("verbose", false, "verbose")
	root.Command("serve", "Run the server")
	root.Command("build", "Build assets")

	err := root.Parse([]string{"--help"})
	require.Error(t, err)
	require.True(t, tinyflags.IsHelpRequested(err))
	assert.Contains(t, err.Error(), "Commands:")
	assert.Contains(t, err.Error(), "serve")
	assert.Contains(t, err.Error(), "build")
}
