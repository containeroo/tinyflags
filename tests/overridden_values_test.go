package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverriddenValuesStatic(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	env := fs.String("env", "dev", "env").Value()
	tags := fs.StringSlice("tag", nil, "tags").Value()

	err := fs.Parse([]string{"--env=prod", "--tag=a,b"})
	require.NoError(t, err)

	got := fs.OverriddenValues()
	assert.Equal(t, "prod", *env)
	assert.Equal(t, []string{"a", "b"}, *tags)
	assert.Equal(t, map[string]any{
		"env": "prod",
		"tag": []string{"a", "b"},
	}, got)
}

func TestOverriddenValuesDynamic(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	http := fs.DynamicGroup("http")
	http.Int("port", 80, "port")

	err := fs.Parse([]string{"--http.a.port=8080"})
	require.NoError(t, err)

	got := fs.OverriddenValues()
	assert.Equal(t, 8080, got["http.a.port"])
	_, ok := got["http.b.port"]
	assert.False(t, ok)
}
