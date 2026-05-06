package tinyflags_test

import (
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
			assert.Equal(t, tt.wantSource, fs.OverriddenValues())
		})
	}
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
