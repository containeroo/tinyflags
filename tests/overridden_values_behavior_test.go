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

// TestOverriddenOriginsStatic reports the exact CLI or environment input that supplied static values.
func TestOverriddenOriginsStatic(t *testing.T) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	fs.String("name", "default", "name").Short("n")
	fs.String("region", "local", "region").Env("DEPLOY_REGION")
	fs.String("mode", "development", "mode")

	t.Setenv("APP_NAME", "environment")
	t.Setenv("DEPLOY_REGION", "production")

	err := fs.Parse([]string{"-n", "flag"})
	require.NoError(t, err)

	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "-n"}, fs.Origin("name"))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceEnvironment, Key: "DEPLOY_REGION"}, fs.Origin("region"))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceDefault}, fs.Origin("mode"))
	assert.Equal(t, "Flag · -n", fs.Origin("name").String())
	assert.Equal(t, "Environment · DEPLOY_REGION", fs.Origin("region").String())
	assert.Equal(t, "Default", fs.Origin("mode").String())
	assert.Equal(t, map[string]tinyflags.ValueOrigin{
		"name":   {Source: tinyflags.ValueSourceFlag, Key: "-n"},
		"region": {Source: tinyflags.ValueSourceEnvironment, Key: "DEPLOY_REGION"},
	}, fs.OverriddenOrigins())
}

// TestOriginReportsExactLongFlagWithoutValue verifies provenance never includes a potentially sensitive flag value.
func TestOriginReportsExactLongFlagWithoutValue(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("token", "", "token")

	require.NoError(t, fs.Parse([]string{"--token=secret"}))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--token"}, fs.Origin("token"))
	assert.NotContains(t, fs.Origin("token").String(), "secret")
}

// TestOriginTracksLastSuccessfulFlagSpelling verifies repeated scalar flags report the spelling that supplied the final value.
func TestOriginTracksLastSuccessfulFlagSpelling(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	name := fs.String("name", "default", "name").Short("n").Value()

	require.NoError(t, fs.Parse([]string{"-n", "first", "--name=second"}))
	assert.Equal(t, "second", *name)
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--name"}, fs.Origin("name"))
}

// TestOverriddenOriginsDynamic reports exact origins for individual dynamic flag IDs.
func TestOverriddenOriginsDynamic(t *testing.T) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	http := fs.DynamicGroup("http")
	http.Int("port", 80, "port")

	t.Setenv("APP_HTTP_API_PORT", "9090")

	err := fs.Parse([]string{"--http.web.port=8080"})
	require.NoError(t, err)

	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--http.web.port"}, fs.Origin("http.web.port"))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceEnvironment, Key: "APP_HTTP_API_PORT"}, fs.Origin("http.api.port"))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceDefault}, fs.Origin("http.worker.port"))
	assert.Equal(t, map[string]tinyflags.ValueOrigin{
		"http.api.port": {Source: tinyflags.ValueSourceEnvironment, Key: "APP_HTTP_API_PORT"},
		"http.web.port": {Source: tinyflags.ValueSourceFlag, Key: "--http.web.port"},
	}, fs.OverriddenOrigins())
}

// TestOriginCLIWinsOverEnvironment verifies the origin follows the same precedence as the effective value.
func TestOriginCLIWinsOverEnvironment(t *testing.T) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	name := fs.String("name", "default", "name").Short("n").Value()

	t.Setenv("APP_NAME", "environment")
	require.NoError(t, fs.Parse([]string{"-n", "flag"}))

	assert.Equal(t, "flag", *name)
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "-n"}, fs.Origin("name"))
}

// TestOverriddenOriginsResetBetweenParses verifies provenance belongs only to the most recent parse.
func TestOverriddenOriginsResetBetweenParses(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("name", "default", "name")

	require.NoError(t, fs.Parse([]string{"--name=alice"}))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--name"}, fs.Origin("name"))

	require.NoError(t, fs.Parse(nil))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceDefault}, fs.Origin("name"))
	assert.Empty(t, fs.OverriddenOrigins())
}

// TestOverriddenOriginsReturnsCopy verifies callers cannot mutate parser provenance state.
func TestOverriddenOriginsReturnsCopy(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("name", "default", "name")

	require.NoError(t, fs.Parse([]string{"--name=alice"}))
	origins := fs.OverriddenOrigins()
	origins["name"] = tinyflags.ValueOrigin{Source: tinyflags.ValueSourceEnvironment, Key: "APP_NAME"}

	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--name"}, fs.Origin("name"))
}
