package engine

import (
	"testing"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDynamicEnvironmentRegressionCases(t *testing.T) {
	t.Parallel()

	t.Run("longerFieldNameWinsSuffixAmbiguity", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		fs.EnvPrefix("APP")
		fs.getEnvVars = func() []string { return []string{"APP_HTTP_API_PRIMARY_ADDR=10.0.0.1"} }
		http := fs.DynamicGroup("http")
		addr := http.String("addr", "default", "address")
		primaryAddr := http.String("primary-addr", "default", "primary address")

		require.NoError(t, fs.Parse(nil))
		assert.False(t, addr.Has("api_primary"))
		assert.Equal(t, "10.0.0.1", primaryAddr.MustGet("api"))
	})

	t.Run("environmentIdentifierIsNormalizedToLowercase", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		fs.EnvPrefix("APP")
		fs.getEnvVars = func() []string { return []string{"APP_HTTP_BACKEND_ONE_ADDR=10.0.0.2"} }
		http := fs.DynamicGroup("http")
		addr := http.String("addr", "default", "address")

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "10.0.0.2", addr.MustGet("backend_one"))
		assert.Equal(t, ValueOrigin{Source: ValueSourceEnvironment, Key: "APP_HTTP_BACKEND_ONE_ADDR"}, fs.Origin("http.backend_one.addr"))
	})

	t.Run("emptyDynamicEnvironmentValueIsIgnored", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		fs.EnvPrefix("APP")
		fs.getEnvVars = func() []string { return []string{"APP_HTTP_API_ADDR="} }
		http := fs.DynamicGroup("http")
		addr := http.String("addr", "default", "address")

		require.NoError(t, fs.Parse(nil))
		assert.False(t, addr.Has("api"))
		assert.Equal(t, ValueOrigin{Source: ValueSourceDefault}, fs.Origin("http.api.addr"))
	})

	t.Run("unrelatedEnvironmentEntryIsIgnored", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		fs.EnvPrefix("APP")
		fs.getEnvVars = func() []string { return []string{"APP_OTHER_API_ADDR=10.0.0.3"} }
		http := fs.DynamicGroup("http")
		addr := http.String("addr", "default", "address")

		require.NoError(t, fs.Parse(nil))
		assert.False(t, addr.Has("api"))
	})

	t.Run("dynamicEnvironmentIsReevaluatedAcrossParses", func(t *testing.T) {
		t.Parallel()

		value := "10.0.0.4"
		fs := NewFlagSet("app", ContinueOnError)
		fs.EnvPrefix("APP")
		fs.getEnvVars = func() []string { return []string{"APP_HTTP_API_ADDR=" + value} }
		http := fs.DynamicGroup("http")
		addr := http.String("addr", "default", "address")

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "10.0.0.4", addr.MustGet("api"))

		value = "10.0.0.5"
		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "10.0.0.5", addr.MustGet("api"))
	})

	t.Run("canonicalDynamicKeyNormalizesGroupAndField", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "MY_APP_HTTP_API_PRIMARY_ADDR", core.DynamicEnvKey("my-app", "http", "API", "primary-addr"))
	})
}
