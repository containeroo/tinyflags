package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentRegressionCases(t *testing.T) {
	t.Parallel()

	t.Run("canonicalStaticKeyNormalizesSeparators", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")

		assert.Equal(t, "APP_DB_USER_NAME", fs.EnvKeyForFlag("db.user/name"))
	})

	t.Run("automaticEnvironmentUsesCanonicalKey", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_DB_USER_NAME" {
				return "alice"
			}
			return ""
		})
		name := fs.String("db.user-name", "default", "name").Value()

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "alice", *name)
		assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceEnvironment, Key: "APP_DB_USER_NAME"}, fs.Origin("db.user-name"))
	})

	t.Run("explicitEnvironmentWorksWithoutPrefix", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.SetGetEnvFn(func(key string) string {
			if key == "TOKEN" {
				return "secret"
			}
			return ""
		})
		token := fs.String("token", "default", "token").Env("TOKEN").Value()

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "secret", *token)
		assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceEnvironment, Key: "TOKEN"}, fs.Origin("token"))
	})

	t.Run("disabledEnvironmentKeepsDefault", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_TOKEN" {
				return "secret"
			}
			return ""
		})
		token := fs.String("token", "default", "token").DisableEnv().Value()

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "default", *token)
		assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceDefault}, fs.Origin("token"))
	})

	t.Run("emptyEnvironmentValueIsIgnored", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(string) string { return "" })
		name := fs.String("name", "default", "name").Value()

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "default", *name)
		assert.Empty(t, fs.Overrides())
	})

	t.Run("invalidEnvironmentCanBeIgnored", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.IgnoreInvalidEnv(true)
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_PORT" {
				return "invalid"
			}
			return ""
		})
		port := fs.Int("port", 8080, "port").Value()

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, 8080, *port)
		assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceDefault}, fs.Origin("port"))
	})

	t.Run("cliSkipsInvalidEnvironmentForSameFlag", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_PORT" {
				return "invalid"
			}
			return ""
		})
		port := fs.Int("port", 8080, "port").Value()

		require.NoError(t, fs.Parse([]string{"--port=9090"}))
		assert.Equal(t, 9090, *port)
		assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceFlag, Key: "--port"}, fs.Origin("port"))
	})

	t.Run("environmentIsReevaluatedOnEveryParse", func(t *testing.T) {
		t.Parallel()

		value := "first"
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_NAME" {
				return value
			}
			return ""
		})
		name := fs.String("name", "default", "name").Value()

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "first", *name)

		value = "second"
		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "second", *name)
		assert.Equal(t, tinyflags.ValueOrigin{Source: tinyflags.ValueSourceEnvironment, Key: "APP_NAME"}, fs.Origin("name"))
	})
}
