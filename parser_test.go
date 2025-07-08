package tinyflags

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("HelpRequested triggers HelpRequested error", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		var helpOut bytes.Buffer
		fs.SetOutput(&helpOut)
		fs.UsagePrintMode(PrintNone)

		err := fs.Parse([]string{"--help"})
		assert.True(t, IsHelpRequested(err))
		assert.EqualError(t, err, "Usage: app\nFlags:\n  -h, --help                             show help\n")
	})

	t.Run("VersionRequested triggers VersionRequested error", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		fs.Version("v1.2.3")

		err := fs.Parse([]string{"--version"})
		assert.True(t, IsVersionRequested(err))
		assert.Equal(t, "v1.2.3", err.Error())
	})

	t.Run("Missing required positional arg returns error", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		fs.RequirePositional(2)

		err := fs.Parse([]string{"one"})
		assert.EqualError(t, err, "expected at least 2 positional arguments, got 1")
	})

	t.Run("Env fallback sets flag if not changed", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		fs.SetGetEnvFn(func(k string) string {
			if k == "USERNAME" {
				return "alice"
			}
			return ""
		})
		f := fs.String("name", "", "").Env("USERNAME")

		err := fs.Parse([]string{})
		assert.NoError(t, err)
		assert.Equal(t, "alice", *f.Value())
	})

	t.Run("Env fallback with invalid value triggers error", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		fs.SetGetEnvFn(func(string) string {
			return "not-an-int"
		})
		fs.IgnoreInvalidEnv(false)
		fs.Int("port", 0, "").Env("PORT")

		err := fs.Parse([]string{})
		assert.EqualError(t, err, "invalid environment value for port: strconv.Atoi: parsing \"not-an-int\": invalid syntax")
	})

	t.Run("Env fallback ignored on parse error if IgnoreInvalidEnv is true", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		fs.SetGetEnvFn(func(string) string {
			return "not-an-int"
		})
		fs.IgnoreInvalidEnv(true)
		fs.Int("port", 1234, "").Env("PORT")

		err := fs.Parse([]string{})
		assert.NoError(t, err)
	})

	t.Run("Mutual exclusion error when two flags in same group are set", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		a := fs.Bool("foo", false, "").Group("G")
		b := fs.Bool("bar", false, "").Group("G")

		err := fs.Parse([]string{"--foo", "--bar"})
		assert.ErrorContains(t, err, "mutually exclusive flags used in group \"G\"")
		assert.True(t, *a.Value())
		assert.True(t, *b.Value())
	})

	t.Run("Required flag missing returns error", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		fs.String("file", "", "").Required()

		err := fs.Parse([]string{})
		assert.ErrorContains(t, err, "flag --file is required")
	})

	t.Run("handleError returns error in ContinueOnError", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		err := fs.handleError(errors.New("bang"))
		assert.ErrorContains(t, err, "bang")
	})

	t.Run("handleError panics in PanicOnError", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", PanicOnError)
		assert.PanicsWithError(t, "boom", func() {
			_ = fs.handleError(errors.New("boom"))
		})
	})
}
