package tinyflags_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
