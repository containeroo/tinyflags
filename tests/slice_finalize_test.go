package tinyflags_test

import (
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlice_Finalize(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)

	calls := 0
	names := fs.StringSlice("name", nil, "User names").
		Finalize(func(s string) string {
			calls++
			return strings.ToUpper(strings.TrimSpace(s))
		}).
		Value()

	err := fs.Parse([]string{
		"--name=alice, bob",
		"--name=carol",
	})
	require.NoError(t, err)

	assert.Equal(t, []string{"ALICE", "BOB", "CAROL"}, *names)
	assert.Equal(t, 3, calls)
}

func TestDynamicSlice_Finalize(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)

	http := fs.DynamicGroup("http")
	hosts := http.StringSlice("hosts", nil, "API hosts").
		Finalize(func(s string) string {
			return strings.ToLower(strings.TrimSpace(s))
		})

	err := fs.Parse([]string{
		"--http.a.hosts=LOCALHOST, api.EXAMPLE.com",
		"--http.b.hosts=example.net",
	})
	require.NoError(t, err)

	assert.Equal(t, []string{"localhost", "api.example.com"}, hosts.MustGet("a"))
	assert.Equal(t, []string{"example.net"}, hosts.MustGet("b"))
}
