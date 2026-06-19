package engine

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewReplacerEnvKeyFunc verifies env key normalization behavior.
func TestNewReplacerEnvKeyFunc(t *testing.T) {
	t.Parallel()

	t.Run("emptyPrefixReturnsEmpty", func(t *testing.T) {
		t.Parallel()

		fn := NewReplacerEnvKeyFunc(strings.NewReplacer("-", "_"), true)
		assert.Equal(t, "", fn("", "flag-name"))
	})

	t.Run("replacesAndUppercases", func(t *testing.T) {
		t.Parallel()

		fn := NewReplacerEnvKeyFunc(strings.NewReplacer("-", "_", ".", "_"), true)
		assert.Equal(t, "APP_FLAG_NAME_MORE", fn("app", "flag-name.more"))
	})
}

// TestParseEnvBehavior verifies static and dynamic ENV loading rules.
func TestParseEnvBehavior(t *testing.T) {
	t.Run("static automatic env requires prefix", func(t *testing.T) {
		fs := NewFlagSet("app", ContinueOnError)
		fs.SetGetEnvFn(func(key string) string {
			values := map[string]string{
				"NAME":     "bare",
				"APP_NAME": "prefixed",
			}
			return values[key]
		})

		var value string
		fs.StringVar(&value, "name", "default", "desc")

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "default", value)

		fs.EnvPrefix("APP")
		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "prefixed", value)
	})

	t.Run("static explicit env works without prefix", func(t *testing.T) {
		fs := NewFlagSet("app", ContinueOnError)
		fs.SetGetEnvFn(func(key string) string {
			if key == "NAME" {
				return "explicit"
			}
			return ""
		})

		var value string
		fs.StringVar(&value, "name", "default", "desc").Env("NAME")

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "explicit", value)
	})

	t.Run("dynamic env requires prefix", func(t *testing.T) {
		fs := NewFlagSet("app", ContinueOnError)
		fs.getEnvVars = func() []string {
			return []string{"SVC_API_ADDR=bare", "APP_SVC_API_ADDR=prefixed"}
		}

		svc := fs.DynamicGroup("svc")
		addr := svc.String("addr", "default", "desc")

		require.NoError(t, fs.Parse(nil))
		assert.False(t, addr.Has("api"))

		fs.EnvPrefix("APP")
		require.NoError(t, fs.Parse(nil))
		got, ok := addr.Get("api")
		assert.True(t, ok)
		assert.Equal(t, "prefixed", got)
	})

	t.Run("dynamic env does not override cli", func(t *testing.T) {
		fs := NewFlagSet("app", ContinueOnError)
		fs.EnvPrefix("APP")
		fs.getEnvVars = func() []string {
			return []string{"APP_SVC_API_ADDR=env"}
		}

		svc := fs.DynamicGroup("svc")
		addr := svc.String("addr", "default", "desc")

		require.NoError(t, fs.Parse([]string{"--svc.api.addr=cli"}))
		assert.Equal(t, "cli", addr.MustGet("api"))
	})

	t.Run("dynamic invalid env can be ignored", func(t *testing.T) {
		fs := NewFlagSet("app", ContinueOnError)
		fs.EnvPrefix("APP")
		fs.IgnoreInvalidEnv(true)
		fs.getEnvVars = func() []string {
			return []string{"APP_SVC_API_PORT=not-an-int"}
		}

		svc := fs.DynamicGroup("svc")
		port := svc.Int("port", 80, "desc")

		require.NoError(t, fs.Parse(nil))
		assert.False(t, port.Has("api"))
	})

	t.Run("dynamic help shows canonical env key", func(t *testing.T) {
		fs := NewFlagSet("app", ContinueOnError)
		fs.EnvPrefix("APP")
		fs.DynamicGroup("svc").String("addr", "default", "desc")

		help := fs.RenderHelpText()
		assert.Contains(t, help, "APP_SVC_<ID>_ADDR")
		assert.NotContains(t, help, "APP_ADDR")
	})
}
