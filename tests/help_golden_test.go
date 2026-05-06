package tinyflags_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/require"
)

// TestHelpOutputGolden_Static verifies static help output against a golden file.
func TestHelpOutputGolden_Static(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	fs.SetStaticUsageIndent(2)
	fs.SetStaticUsageColumn(24)
	fs.Title("Flags:")
	fs.Description("Static help output.")

	fs.String("name", "demo", "Service name")
	fs.Int("port", 8080, "API port").Short("p")
	fs.String("mode", "dev", "Mode").Choices("dev", "prod")

	err := fs.Parse([]string{"--help"})
	require.Error(t, err)
	require.True(t, tinyflags.IsHelpRequested(err))

	assertGolden(t, "testdata/help_static.golden", err.Error())
}

// TestHelpOutputGolden_Dynamic verifies dynamic help output against a golden file.
func TestHelpOutputGolden_Dynamic(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.SetDynamicUsageIndent(2)
	fs.SetDynamicUsageColumn(34)
	fs.Title("Flags:")

	http := fs.DynamicGroup("http")
	http.Title("HTTP Targets")
	http.Description("Configure per-target HTTP settings.")
	http.Note("Every target may define its own address and port.")
	http.Placeholder("<TARGET>")
	http.String("addr", "", "Target address")
	http.Int("port", 80, "Target port")

	err := fs.Parse([]string{"--help"})
	require.Error(t, err)
	require.True(t, tinyflags.IsHelpRequested(err))

	assertGolden(t, "testdata/help_dynamic.golden", err.Error())
}

// assertGolden compares rendered output against a golden file.
func assertGolden(t *testing.T, relativePath string, actual string) {
	t.Helper()

	path := filepath.Join(relativePath)
	expectedBytes, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, string(expectedBytes), actual)
}
