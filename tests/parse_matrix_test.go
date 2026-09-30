package tinyflags_test

import (
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseValuePrecedenceMatrix verifies default, env, and CLI precedence.
func TestParseValuePrecedenceMatrix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		wantValue  string
		wantSource map[string]any
	}{
		{
			name:       "defaultOnly",
			wantValue:  "default",
			wantSource: map[string]any{},
		},
		{
			name:      "envOverridesDefault",
			env:       map[string]string{"APP_NAME": "from-env"},
			wantValue: "from-env",
			wantSource: map[string]any{
				"name": "from-env",
			},
		},
		{
			name:      "argsOverrideEnv",
			args:      []string{"--name=from-arg"},
			env:       map[string]string{"APP_NAME": "from-env"},
			wantValue: "from-arg",
			wantSource: map[string]any{
				"name": "from-arg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
			fs.EnvPrefix("APP")
			fs.SetGetEnvFn(func(key string) string { return tt.env[key] })
			name := fs.String("name", "default", "name").Value()

			err := fs.Parse(tt.args)
			require.NoError(t, err)
			assert.Equal(t, tt.wantValue, *name)
			assert.Equal(t, tt.wantSource, fs.Overrides().Values())
		})
	}
}

// TestNotEmpty verifies presence and non-empty value checks remain independent.
func TestNotEmpty(t *testing.T) {
	t.Parallel()

	t.Run("optionalUnset", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("token", "", "token").NotEmpty()

		require.NoError(t, fs.Parse(nil))
	})

	t.Run("requiredStillAllowsExplicitEmpty", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("token", "", "token").Required()

		require.NoError(t, fs.Parse([]string{"--token="}))
	})

	t.Run("staticRejectsEmpty", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("token", "", "token").NotEmpty()

		err := fs.Parse([]string{"--token="})
		require.ErrorContains(t, err, "flag --token must not be empty")
	})

	t.Run("staticChecksFinalizedValue", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("token", "", "token").Finalize(strings.TrimSpace).NotEmpty()

		err := fs.Parse([]string{"--token=   "})
		require.ErrorContains(t, err, "flag --token must not be empty")
	})

	t.Run("staticAcceptsValue", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		token := fs.String("token", "", "token").Required().NotEmpty().Value()

		require.NoError(t, fs.Parse([]string{"--token=abc"}))
		assert.Equal(t, "abc", *token)
	})

	t.Run("zeroNumberIsEmpty", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Int("count", 1, "count").NotEmpty()

		err := fs.Parse([]string{"--count=0"})
		require.ErrorContains(t, err, "flag --count must not be empty")
	})

	t.Run("dynamicRejectsEmpty", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		service := fs.DynamicGroup("service")
		service.String("addr", "", "address").NotEmpty()

		err := fs.Parse([]string{"--service.api.addr="})
		require.ErrorContains(t, err, "flag --service.api.addr must not be empty")
	})

	t.Run("dynamicUnsetFieldIsOptional", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		service := fs.DynamicGroup("service")
		service.String("addr", "", "address").NotEmpty()
		service.Int("port", 0, "port")

		require.NoError(t, fs.Parse([]string{"--service.api.port=8080"}))
	})
}

// TestParseConstraintMatrix verifies grouped parse constraint behavior.
func TestParseConstraintMatrix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		build     func(*tinyflags.FlagSet)
		args      []string
		wantErr   string
		wantNoErr bool
	}{
		{
			name: "requiredSatisfied",
			build: func(fs *tinyflags.FlagSet) {
				fs.String("token", "", "token").Required()
			},
			args:      []string{"--token=abc"},
			wantNoErr: true,
		},
		{
			name: "requiredMissing",
			build: func(fs *tinyflags.FlagSet) {
				fs.String("token", "", "token").Required()
			},
			wantErr: "flag --token is required",
		},
		{
			name: "requiresSatisfied",
			build: func(fs *tinyflags.FlagSet) {
				fs.String("db", "", "db")
				fs.String("dsn", "", "dsn").Requires("db")
			},
			args:      []string{"--db=main", "--dsn=postgres"},
			wantNoErr: true,
		},
		{
			name: "requiresMissingDependency",
			build: func(fs *tinyflags.FlagSet) {
				fs.String("db", "", "db")
				fs.String("dsn", "", "dsn").Requires("db")
			},
			args:    []string{"--dsn=postgres"},
			wantErr: "--dsn requires --db",
		},
		{
			name: "oneOfConflict",
			build: func(fs *tinyflags.FlagSet) {
				fs.Bool("debug", false, "debug").OneOfGroup("mode")
				fs.Bool("quiet", false, "quiet").OneOfGroup("mode")
			},
			args:    []string{"--debug", "--quiet"},
			wantErr: "only one of the flags in group",
		},
		{
			name: "allOrNonePartial",
			build: func(fs *tinyflags.FlagSet) {
				fs.String("user", "", "user").AllOrNone("auth")
				fs.String("pass", "", "pass").AllOrNone("auth")
			},
			args:    []string{"--user=alice"},
			wantErr: "must be set together",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
			tt.build(fs)

			err := fs.Parse(tt.args)
			if tt.wantNoErr {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
