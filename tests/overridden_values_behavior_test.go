package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOverriddenValuesStatic verifies overridden values for static flags.
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

// TestOverriddenValuesDynamic verifies overridden values for dynamic flags.
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

// TestOverriddenValuesMaskFn verifies masking in overridden values.
func TestOverriddenValuesMaskFn(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	secret := fs.String("secret", "", "secret").
		OverriddenValueMaskFn(tinyflags.MaskFirstLast).
		Value()

	err := fs.Parse([]string{"--secret=opensesame"})
	require.NoError(t, err)

	got := fs.OverriddenValues()
	assert.Equal(t, "opensesame", *secret)
	assert.Equal(t, "o********e", got["secret"])
}

// TestMaskPostgresURL verifies Postgres URL masking in overridden values.
func TestMaskPostgresURL(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	dsn := fs.String("dsn", "", "dsn").
		OverriddenValueMaskFn(tinyflags.MaskPostgresURL).
		Value()

	err := fs.Parse([]string{"--dsn=postgres://user:pass@localhost:5432/app"})
	require.NoError(t, err)

	got := fs.OverriddenValues()
	assert.Equal(t, "postgres://user:pass@localhost:5432/app", *dsn)
	assert.Equal(t, "postgres://****:****@localhost:5432/app", got["dsn"])
}

// TestOverriddenSourcesStatic reports whether static values came from flags or environment variables.
func TestOverriddenSourcesStatic(t *testing.T) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	fs.String("name", "default", "name").Short("n")
	fs.String("region", "local", "region")
	fs.String("mode", "development", "mode")

	t.Setenv("APP_NAME", "environment")
	t.Setenv("APP_REGION", "production")

	err := fs.Parse([]string{"-n", "flag"})
	require.NoError(t, err)

	assert.Equal(t, tinyflags.ValueSourceFlag, fs.Source("name"))
	assert.Equal(t, tinyflags.ValueSourceEnvironment, fs.Source("region"))
	assert.Equal(t, tinyflags.ValueSourceDefault, fs.Source("mode"))
	assert.Equal(t, map[string]tinyflags.ValueSource{
		"name":   tinyflags.ValueSourceFlag,
		"region": tinyflags.ValueSourceEnvironment,
	}, fs.OverriddenSources())
}

// TestOverriddenSourcesDynamic reports sources for individual dynamic flag IDs.
func TestOverriddenSourcesDynamic(t *testing.T) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	http := fs.DynamicGroup("http")
	http.Int("port", 80, "port")

	t.Setenv("APP_HTTP_API_PORT", "9090")

	err := fs.Parse([]string{"--http.web.port=8080"})
	require.NoError(t, err)

	assert.Equal(t, tinyflags.ValueSourceFlag, fs.Source("http.web.port"))
	assert.Equal(t, tinyflags.ValueSourceEnvironment, fs.Source("http.api.port"))
	assert.Equal(t, tinyflags.ValueSourceDefault, fs.Source("http.worker.port"))
	assert.Equal(t, map[string]tinyflags.ValueSource{
		"http.api.port": tinyflags.ValueSourceEnvironment,
		"http.web.port": tinyflags.ValueSourceFlag,
	}, fs.OverriddenSources())
}

// TestOverriddenSourcesResetBetweenParses verifies provenance belongs only to the most recent parse.
func TestOverriddenSourcesResetBetweenParses(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("name", "default", "name")

	require.NoError(t, fs.Parse([]string{"--name=alice"}))
	assert.Equal(t, tinyflags.ValueSourceFlag, fs.Source("name"))

	require.NoError(t, fs.Parse(nil))
	assert.Equal(t, tinyflags.ValueSourceDefault, fs.Source("name"))
	assert.Empty(t, fs.OverriddenSources())
}

// TestOverriddenSourcesReturnsCopy verifies callers cannot mutate parser provenance state.
func TestOverriddenSourcesReturnsCopy(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("name", "default", "name")

	require.NoError(t, fs.Parse([]string{"--name=alice"}))
	sources := fs.OverriddenSources()
	sources["name"] = tinyflags.ValueSourceEnvironment

	assert.Equal(t, tinyflags.ValueSourceFlag, fs.Source("name"))
}
