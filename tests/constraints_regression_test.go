package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstraintRegressionCases(t *testing.T) {
	t.Parallel()

	t.Run("requiredFlagCanBeSatisfiedByEnvironment", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_TOKEN" {
				return "secret"
			}
			return ""
		})
		token := fs.String("token", "", "token").Required().Value()

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, "secret", *token)
	})

	t.Run("requiresCanBeSatisfiedAcrossCliAndEnvironment", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_DATABASE" {
				return "main"
			}
			return ""
		})
		fs.String("database", "", "database")
		dsn := fs.String("dsn", "", "dsn").Requires("database").Value()

		require.NoError(t, fs.Parse([]string{"--dsn=postgres"}))
		assert.Equal(t, "postgres", *dsn)
	})

	t.Run("allOrNoneCanBeSatisfiedAcrossCliAndEnvironment", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_USER" {
				return "alice"
			}
			return ""
		})
		user := fs.String("user", "", "user").AllOrNone("auth").Value()
		password := fs.String("password", "", "password").AllOrNone("auth").Value()

		require.NoError(t, fs.Parse([]string{"--password=secret"}))
		assert.Equal(t, "alice", *user)
		assert.Equal(t, "secret", *password)
	})

	t.Run("requiredOneOfCountsExplicitFalseAsSelected", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_DEBUG" {
				return "false"
			}
			return ""
		})
		debug := fs.Bool("debug", true, "debug").Strict().OneOfGroup("mode").Value()
		fs.Bool("quiet", false, "quiet").Strict().OneOfGroup("mode")
		fs.GetOneOfGroup("mode").Required()

		require.NoError(t, fs.Parse(nil))
		assert.False(t, *debug)
	})

	t.Run("dynamicRequiredFieldIsCheckedForExistingInstance", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		service := fs.DynamicGroup("service")
		service.String("addr", "", "address").Required()
		service.Int("port", 80, "port")

		err := fs.Parse([]string{"--service.api.port=8080"})
		require.Error(t, err)
		assert.EqualError(t, err, "flag --service.api.addr is required")
	})

	t.Run("dynamicRequiredFieldIsCheckedPerInstance", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		service := fs.DynamicGroup("service")
		service.String("addr", "", "address").Required()
		service.Int("port", 80, "port")

		err := fs.Parse([]string{
			"--service.alpha.addr=10.0.0.1",
			"--service.beta.port=8080",
		})
		require.Error(t, err)
		assert.EqualError(t, err, "flag --service.beta.addr is required")
	})

	t.Run("helpBypassesRequiredAndRelationshipChecks", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("required", "", "required").Required()
		fs.String("user", "", "user").AllOrNone("auth")
		fs.String("password", "", "password").AllOrNone("auth")
		fs.Bool("left", false, "left").OneOfGroup("choice")
		fs.Bool("right", false, "right").OneOfGroup("choice")
		fs.GetOneOfGroup("choice").Required()

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.True(t, tinyflags.IsHelpRequested(err))
		assert.Contains(t, err.Error(), "--required REQUIRED")
	})
}
