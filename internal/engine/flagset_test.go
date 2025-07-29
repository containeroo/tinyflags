package engine_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFlagSet(t *testing.T) {
	t.Parallel()

	t.Run("Usage override", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		var buf bytes.Buffer
		fs.SetOutput(&buf)
		fs.Usage = func() {
			fmt.Fprint(fs.Output(), "CUSTOM USAGE\n") // nolint:errcheck
		}
		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.Equal(t, "CUSTOM USAGE\n", err.Error())
	})

	t.Run("Version flag", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		fs.Version("vX.Y.Z")
		err := fs.Parse([]string{"--version"})
		var vr *tinyflags.VersionRequested
		require.True(t, errors.As(err, &vr), "expected VersionRequested")
		assert.Equal(t, "vX.Y.Z", vr.Version)
	})

	t.Run("DisableVersion", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		fs.DisableVersion()
		err := fs.Parse([]string{"--version"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown flag")
	})

	t.Run("Environment controls", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		fs.EnvPrefix("PRE")
		f := fs.String("f", "fallback", "desc")
		fs.String("ignore", "x", "desc").DisableEnv().Value()
		fs.IgnoreInvalidEnv(true)
		fs.SetGetEnvFn(func(k string) string {
			m := map[string]string{
				"PRE_F": "fromenv",
			}
			return m[k]
		})
		val := f.Value()
		require.NoError(t, fs.Parse([]string{}))
		assert.Equal(t, "fromenv", *val)
	})

	t.Run("Help sections", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		var buf bytes.Buffer
		fs.SetOutput(&buf)
		fs.Authors("A & B")
		fs.Title("T")
		fs.Description("D")
		fs.Note("N")

		fs.PrintAuthors(fs.Output())
		fs.PrintTitle(fs.Output())
		fs.PrintDescription(fs.Output(), 0, 40)
		fs.PrintNotes(fs.Output(), 0, 40)

		out := buf.String()
		assert.Contains(t, out, "A & B")
		assert.Contains(t, out, "T")
		assert.Contains(t, out, "D")
		assert.Contains(t, out, "N")
	})

	t.Run("DisableHelp", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		fs.DisableHelp()
		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown flag")
	})

	t.Run("Sorted", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		fs.SortedFlags()
		fs.SortedGroups()
	})

	t.Run("Output setter/getter", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		var b bytes.Buffer
		fs.SetOutput(&b)
		assert.Same(t, &b, fs.Output())
	})

	t.Run("Globaldelimiter", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		fs.Globaldelimiter(";")
		// now a slice flag uses ";" by default
		s := fs.StringSlice("s", nil, "desc").Value()
		require.NoError(t, fs.Parse([]string{"--s=a;b;c"}))
		assert.Equal(t, []string{"a", "b", "c"}, *s)
	})

	t.Run("Positional args", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		fs.RequirePositional(2)
		err := fs.Parse([]string{"one", "two"})
		require.NoError(t, err)
		assert.Equal(t, []string{"one", "two"}, fs.Args())
		a0, ok := fs.Arg(0)
		assert.True(t, ok)
		assert.Equal(t, "one", a0)
		_, ok = fs.Arg(5)
		assert.False(t, ok)
	})

	t.Run("PrintDefaults and Usage", func(t *testing.T) {
		var buf bytes.Buffer
		t.Parallel()
		fs := tinyflags.NewFlagSet("myapp", tinyflags.ContinueOnError)
		fs.SetOutput(&buf)
		fs.String("x", "def", "desc")
		fs.PrintStaticDefaults(fs.Output(), 2, 40, 200)
		fs.PrintUsage(fs.Output(), tinyflags.PrintShort)
		out := buf.String()
		assert.Contains(t, out, "--x")
		assert.Contains(t, out, "Usage:")
	})
}
