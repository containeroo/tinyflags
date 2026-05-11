package tinyflags_test

import (
	"bytes"
	"context"
	"errors"
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

func TestCommandHelpTextAndWriteHelp(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	root.Globals().Bool("verbose", false, "verbose")
	serve := root.Command("serve", "Run the server")
	serve.Int("port", 8080, "port")

	helpText := root.HelpText()
	assert.Contains(t, helpText, "Usage: app")
	assert.Contains(t, helpText, "Commands:")
	assert.Contains(t, helpText, "serve")

	var buf bytes.Buffer
	require.NoError(t, root.WriteHelp(&buf))
	assert.Equal(t, helpText, buf.String())
}

func TestCommandHelpTextDoesNotMutateParsedPositionals(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	pin := root.Command("pin", "Pin a command.")

	err := root.Parse([]string{"pin", "--", "kubectl", "get", "pods"})
	require.NoError(t, err)
	assert.Equal(t, []string{"kubectl", "get", "pods"}, pin.Args())

	helpText := root.HelpText()
	assert.Contains(t, helpText, "Usage: app pin")
	assert.Equal(t, []string{"kubectl", "get", "pods"}, pin.Args())
}

func TestRequireCommandRoot(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError).RequireCommand()
	root.Globals().Bool("verbose", false, "verbose")
	root.Command("serve", "Run the server")

	err := root.Parse([]string{"--verbose"})
	require.Error(t, err)
	assert.True(t, tinyflags.IsCommandRequired(err))
	assert.EqualError(t, err, `command "app" requires a subcommand`)
	help, ok := tinyflags.HelpText(err)
	require.True(t, ok)
	assert.Contains(t, help, `Usage: app [flags] <command>`)
	assert.Contains(t, help, "Commands:")

	var typed *tinyflags.CommandRequired
	assert.True(t, errors.As(err, &typed))
	require.NotNil(t, typed)
	assert.Equal(t, "app", typed.Command)

	var usageErr *tinyflags.UsageError
	assert.True(t, errors.As(err, &usageErr))
	require.NotNil(t, usageErr)
	assert.Equal(t, help, usageErr.Help)
}

func TestRequireCommandNested(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	admin := root.Command("admin", "Admin tools").RequireCommand()
	admin.Globals().Bool("audit", false, "audit")
	admin.Command("users", "Manage users")

	err := root.Parse([]string{"admin", "--audit"})
	require.Error(t, err)
	assert.True(t, tinyflags.IsCommandRequired(err))
	assert.EqualError(t, err, `command "app admin" requires a subcommand`)
	help, ok := tinyflags.HelpText(err)
	require.True(t, ok)
	assert.Contains(t, help, `Usage: app admin [flags] <command>`)
	assert.Contains(t, help, "users")
}

// TestParseRunnerInvokesSelectedHandler verifies handler registration receives parsed values.
func TestParseRunnerInvokesSelectedHandler(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	verbose := root.Globals().Bool("verbose", false, "verbose").Value()
	root.Run(func(_ context.Context, gotVerbose bool) error {
		assert.Equal(t, *verbose, gotVerbose)
		return nil
	}, verbose)

	serve := root.Command("serve", "Run the server")
	port := serve.Int("port", 8080, "port").Value()
	serve.Run(func(_ context.Context, gotVerbose bool, gotPort int) error {
		assert.Equal(t, *verbose, gotVerbose)
		assert.Equal(t, 9000, gotPort)
		return nil
	}, verbose, port)

	runner, err := root.ParseRunner([]string{"--verbose", "serve", "--port=9000"})
	require.NoError(t, err)
	require.NotNil(t, runner)
	require.NoError(t, runner.Run(context.Background()))
}

// TestParseRunnerInjectsContext verifies one registered handler receives the execution context.
func TestParseRunnerInjectsContext(t *testing.T) {
	t.Parallel()

	type contextKey struct{}

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	root.Run(func(ctx context.Context) error {
		assert.Equal(t, "trace-123", ctx.Value(contextKey{}))
		return nil
	})

	runner, err := root.ParseRunner(nil)
	require.NoError(t, err)
	require.NotNil(t, runner)
	ctx := context.WithValue(context.Background(), contextKey{}, "trace-123")
	require.NoError(t, runner.Run(ctx))
}

// TestParseRunnerUsesFrozenValues verifies later flag mutations do not change one parsed runner.
func TestParseRunnerUsesFrozenValues(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	serve := root.Command("serve", "Run the server")
	port := serve.Int("port", 8080, "port").Value()

	var calledPort int
	serve.Run(func(_ context.Context, gotPort int) error {
		calledPort = gotPort
		return nil
	}, port)

	runner, err := root.ParseRunner([]string{"serve", "--port=9000"})
	require.NoError(t, err)
	*port = 7000
	require.NoError(t, runner.Run(context.Background()))
	assert.Equal(t, 9000, calledPort)
}

// TestParseRunnerRequiresHandler verifies missing handlers fail clearly after parsing.
func TestParseRunnerRequiresHandler(t *testing.T) {
	t.Parallel()

	root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	root.Command("serve", "Run the server")

	runner, err := root.ParseRunner([]string{"serve"})
	require.Error(t, err)
	assert.Nil(t, runner)
	assert.Contains(t, err.Error(), `no command runner registered for command "app serve"`)
}
