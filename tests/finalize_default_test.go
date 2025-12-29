package tinyflags_test

import (
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinalizeDefaultValueScalar(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	val := fs.String("path", " /config ", "config path").
		Finalize(strings.TrimSpace).
		FinalizeDefaultValue().
		Value()

	err := fs.Parse([]string{})
	require.NoError(t, err)

	assert.Equal(t, "/config", *val)
}

func TestFinalizeDefaultValueScalarSetValue(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	calls := 0
	val := fs.String("name", " default ", "name").
		Finalize(func(s string) string {
			calls++
			return strings.TrimSpace(s)
		}).
		FinalizeDefaultValue().
		Value()

	err := fs.Parse([]string{"--name=  bob  "})
	require.NoError(t, err)

	assert.Equal(t, "bob", *val)
	assert.Equal(t, 1, calls)
}

func TestFinalizeDefaultValueSlice(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	val := fs.StringSlice("names", []string{"alice", " bob "}, "names").
		Finalize(strings.TrimSpace).
		FinalizeDefaultValue().
		Value()

	err := fs.Parse([]string{})
	require.NoError(t, err)

	assert.Equal(t, []string{"alice", "bob"}, *val)
}

func TestFinalizeDefaultValueDynamic(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	http := fs.DynamicGroup("http")
	host := http.String("host", " LOCAL ", "host").
		Finalize(strings.TrimSpace).
		FinalizeDefaultValue()

	err := fs.Parse([]string{})
	require.NoError(t, err)

	val, ok := host.Get("a")
	assert.False(t, ok)
	assert.Equal(t, "LOCAL", val)
}
