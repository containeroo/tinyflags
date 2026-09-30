package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOverridesStatic verifies values and origins for static overrides.
func TestOverridesStatic(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	env := fs.String("env", "dev", "env").Value()
	tags := fs.StringSlice("tag", nil, "tags").Value()

	err := fs.Parse([]string{"--env=prod", "--tag=a,b"})
	require.NoError(t, err)

	got := fs.Overrides()
	assert.Equal(t, "prod", *env)
	assert.Equal(t, []string{"a", "b"}, *tags)
	assert.Equal(t, tinyflags.Overrides{
		"env": {
			Value:  "prod",
			Origin: tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--env"},
		},
		"tag": {
			Value:  []string{"a", "b"},
			Origin: tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--tag"},
		},
	}, got)
	assert.Equal(t, map[string]any{
		"env": "prod",
		"tag": []string{"a", "b"},
	}, got.Values())
}

// TestOverridesDynamic verifies values and origins for dynamic overrides.
func TestOverridesDynamic(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	http := fs.DynamicGroup("http")
	http.Int("port", 80, "port")

	err := fs.Parse([]string{"--http.a.port=8080"})
	require.NoError(t, err)

	got := fs.Overrides()
	assert.Equal(t, 8080, got["http.a.port"].Value)
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--http.a.port"}, got["http.a.port"].Origin)
	_, ok := got["http.b.port"]
	assert.False(t, ok)
}

// TestOverridesMaskFn verifies masking applies only to the reporting override value.
func TestOverridesMaskFn(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	secret := fs.String("secret", "", "secret").
		OverriddenValueMaskFn(tinyflags.MaskFirstLast).
		Value()

	err := fs.Parse([]string{"--secret=opensesame"})
	require.NoError(t, err)

	got := fs.Overrides()
	assert.Equal(t, "opensesame", *secret)
	assert.Equal(t, "o********e", got["secret"].Value)
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--secret"}, got["secret"].Origin)
}

// TestMaskPostgresURL verifies Postgres URL masking in override values.
func TestMaskPostgresURL(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	dsn := fs.String("dsn", "", "dsn").
		OverriddenValueMaskFn(tinyflags.MaskPostgresURL).
		Value()

	err := fs.Parse([]string{"--dsn=postgres://user:pass@localhost:5432/app"})
	require.NoError(t, err)

	got := fs.Overrides()
	assert.Equal(t, "postgres://user:pass@localhost:5432/app", *dsn)
	assert.Equal(t, "postgres://****:****@localhost:5432/app", got["dsn"].Value)
}

// TestOverrideOriginsStatic reports the exact CLI or environment input that supplied static values.
func TestOverrideOriginsStatic(t *testing.T) {
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
	}, fs.Overrides().Origins())
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

// TestOverrideOriginsDynamic reports exact origins for individual dynamic flag IDs.
func TestOverrideOriginsDynamic(t *testing.T) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	http := fs.DynamicGroup("http")
	http.Int("port", 80, "port")

	t.Setenv("APP_HTTP_API_PORT", "9090")

	err := fs.Parse([]string{"--http.web.port=8080"})
	require.NoError(t, err)

	overrides := fs.Overrides()
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--http.web.port"}, fs.Origin("http.web.port"))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceEnvironment, Key: "APP_HTTP_API_PORT"}, fs.Origin("http.api.port"))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceDefault}, fs.Origin("http.worker.port"))
	assert.Equal(t, map[string]tinyflags.ValueOrigin{
		"http.api.port": {Source: tinyflags.ValueSourceEnvironment, Key: "APP_HTTP_API_PORT"},
		"http.web.port": {Source: tinyflags.ValueSourceFlag, Key: "--http.web.port"},
	}, overrides.Origins())
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

// TestOverridesResetBetweenParses verifies override metadata belongs only to the most recent parse.
func TestOverridesResetBetweenParses(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("name", "default", "name")

	require.NoError(t, fs.Parse([]string{"--name=alice"}))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--name"}, fs.Origin("name"))
	assert.Len(t, fs.Overrides(), 1)

	require.NoError(t, fs.Parse(nil))
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceDefault}, fs.Origin("name"))
	assert.Empty(t, fs.Overrides())
}

// TestOverridesReturnsCopy verifies callers cannot mutate parser override state.
func TestOverridesReturnsCopy(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("name", "default", "name")

	require.NoError(t, fs.Parse([]string{"--name=alice"}))
	overrides := fs.Overrides()
	overrides["name"] = tinyflags.Override{
		Value:  "changed",
		Origin: tinyflags.ValueOrigin{Source: tinyflags.ValueSourceEnvironment, Key: "APP_NAME"},
	}

	got := fs.Overrides()["name"]
	assert.Equal(t, "alice", got.Value)
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--name"}, got.Origin)
}

// TestOverrideProjectionsReturnCopies verifies Values and Origins do not share their result maps.
func TestOverrideProjectionsReturnCopies(t *testing.T) {
	t.Parallel()

	overrides := tinyflags.Overrides{
		"name": {
			Value:  "alice",
			Origin: tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--name"},
		},
	}

	values := overrides.Values()
	origins := overrides.Origins()
	values["name"] = "bob"
	origins["name"] = tinyflags.ValueOrigin{Source: tinyflags.ValueSourceEnvironment, Key: "APP_NAME"}

	assert.Equal(t, "alice", overrides["name"].Value)
	assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--name"}, overrides["name"].Origin)
}
