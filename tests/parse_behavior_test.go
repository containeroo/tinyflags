package tinyflags_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSliceStrictDelimiter verifies strict delimiter handling for slice flags.
func TestSliceStrictDelimiter(t *testing.T) {
	t.Parallel()

	t.Run("mixedDelimiterErrors", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		names := fs.StringSlice("name", nil, "names").
			Delimiter("|").
			StrictDelimiter().
			Value()

		err := fs.Parse([]string{"--name=a,b|c"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mixed delimiters")
		assert.Nil(t, *names)
	})
}

// TestSliceAllowEmpty verifies empty-item handling for slice flags.
func TestSliceAllowEmpty(t *testing.T) {
	t.Parallel()

	t.Run("keepsEmptyItems", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		names := fs.StringSlice("name", nil, "names").
			Delimiter(",").
			AllowEmpty().
			Value()

		err := fs.Parse([]string{"--name=a,,b"})
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "", "b"}, *names)
	})

	t.Run("dynamicAllowEmpty", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		g := fs.DynamicGroup("g")
		values := g.StringSlice("v", nil, "vals").
			AllowEmpty()

		err := fs.Parse([]string{"--g.one.v=a,,b"})
		require.NoError(t, err)

		assert.Equal(t, []string{"a", "", "b"}, values.MustGet("one"))
	})
}

// TestUnknownFlagHandler verifies custom handling for unknown flags.
func TestUnknownFlagHandler(t *testing.T) {
	t.Parallel()

	t.Run("ignoredUnknown", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.OnUnknownFlag(func(name string) error {
			return nil
		})
		val := fs.String("val", "", "value").Value()

		err := fs.Parse([]string{"--unknown=1", "--val=ok"})
		require.NoError(t, err)
		assert.Equal(t, "ok", *val)
	})

	t.Run("customError", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.OnUnknownFlag(func(name string) error {
			return errors.New("custom unknown: " + name)
		})

		err := fs.Parse([]string{"--nope"})
		require.Error(t, err)
		assert.EqualError(t, err, "custom unknown: --nope")
	})
}

// TestBeforeParseHook verifies argument preprocessing hooks.
func TestBeforeParseHook(t *testing.T) {
	t.Parallel()

	t.Run("mutatesArgs", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.BeforeParse(func(args []string) ([]string, error) {
			out := make([]string, len(args))
			for i, a := range args {
				out[i] = strings.ReplaceAll(a, "--NAME", "--name")
			}
			return out, nil
		})
		name := fs.String("name", "", "name").Value()

		err := fs.Parse([]string{"--NAME=alice"})
		require.NoError(t, err)
		assert.Equal(t, "alice", *name)
	})

	t.Run("returnsError", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.BeforeParse(func(args []string) ([]string, error) {
			return nil, errors.New("boom")
		})

		err := fs.Parse([]string{"--name=ignored"})
		require.EqualError(t, err, "boom")
	})
}

// TestDynamicFinalizeWithID verifies dynamic finalizers that receive IDs.
func TestDynamicFinalizeWithID(t *testing.T) {
	t.Parallel()

	t.Run("sliceFinalizerUsesID", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		http := fs.DynamicGroup("http")
		hosts := http.StringSlice("host", nil, "host").
			FinalizeWithID(func(id, v string) string { return id + "-" + v })

		err := fs.Parse([]string{
			"--http.alpha.host=a",
			"--http.beta.host=b",
		})
		require.NoError(t, err)

		assert.Equal(t, []string{"alpha-a"}, hosts.MustGet("alpha"))
		assert.Equal(t, []string{"beta-b"}, hosts.MustGet("beta"))
	})

	t.Run("scalarFinalizerUsesID", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		db := fs.DynamicGroup("db")
		user := db.String("user", "", "user").
			FinalizeWithID(func(id, v string) string { return id + ":" + v })

		err := fs.Parse([]string{
			"--db.main.user=admin",
		})
		require.NoError(t, err)

		val, ok := user.Get("main")
		require.True(t, ok)
		assert.Equal(t, "main:admin", val)
	})

	t.Run("sliceStrictDelimiterDynamic", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		http := fs.DynamicGroup("http")
		hosts := http.StringSlice("hosts", nil, "hosts").
			Delimiter(",").
			StrictDelimiter()

		err := fs.Parse([]string{"--http.a.hosts=a|b"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mixed delimiters")

		_, ok := hosts.Get("a")
		assert.False(t, ok)
	})
}

// TestContinueOnErrorParsesAllFlags verifies continued parsing after recoverable errors.
func TestContinueOnErrorParsesAllFlags(t *testing.T) {
	t.Parallel()

	t.Run("unknownFlagDoesNotStopParsing", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		a := fs.String("a", "", "a").Value()
		b := fs.String("b", "", "b").Value()

		err := fs.Parse([]string{"--a=ok", "--unknown", "--b=ok"})
		require.Error(t, err)
		assert.Equal(t, "ok", *a)
		assert.Equal(t, "ok", *b)
		assert.Contains(t, err.Error(), "unknown flag")
	})

	t.Run("missingValueDoesNotStopParsing", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		a := fs.String("a", "", "a").Value()
		b := fs.String("b", "", "b").Value()

		err := fs.Parse([]string{"--a", "--b=ok"})
		require.Error(t, err)
		assert.Equal(t, "", *a)
		assert.Equal(t, "ok", *b)
		assert.Contains(t, err.Error(), "missing value for flag")
	})
}

// TestParseResetsStateBetweenCalls verifies FlagSet state resets on repeated parses.
func TestParseResetsStateBetweenCalls(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.RequirePositional(1)

	name := fs.String("name", "default", "name")
	http := fs.DynamicGroup("http")
	port := http.Int("port", 80, "port")

	err := fs.Parse([]string{"--name=alice", "--http.a.port=8080", "one"})
	require.NoError(t, err)
	assert.Equal(t, "alice", *name.Value())
	assert.True(t, name.Changed())
	assert.Equal(t, 8080, port.MustGet("a"))
	assert.Equal(t, []string{"one"}, fs.Args())
	assert.Equal(t, map[string]any{
		"http.a.port": 8080,
		"name":        "alice",
	}, fs.OverriddenValues())

	err = fs.Parse([]string{"two"})
	require.NoError(t, err)
	assert.Equal(t, "default", *name.Value())
	assert.False(t, name.Changed())
	assert.False(t, port.Has("a"))
	assert.Equal(t, []string{"two"}, fs.Args())
	assert.Empty(t, fs.OverriddenValues())
}

// TestOneOfGroupVerboseToggle verifies verbose one-of error toggling.
func TestOneOfGroupVerboseToggle(t *testing.T) {
	t.Parallel()

	t.Run("verboseDefaultIncludesConflicts", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Bool("debug", false, "debug").OneOfGroup("dbg")
		fs.Bool("no-debug", false, "no-debug").OneOfGroup("dbg")

		err := fs.Parse([]string{"--debug", "--no-debug"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "vs")
		assert.Contains(t, err.Error(), "--debug")
		assert.Contains(t, err.Error(), "--no-debug")
	})

	t.Run("verboseDisabledOmitsConflicts", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.SetOneOfGroupVerbose(false)
		fs.Bool("debug", false, "debug").OneOfGroup("dbg")
		fs.Bool("no-debug", false, "no-debug").OneOfGroup("dbg")

		err := fs.Parse([]string{"--debug", "--no-debug"})
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "vs")
	})
}
