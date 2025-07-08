package tinyflags

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagSet(t *testing.T) {
	t.Parallel()

	t.Run("name", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("myapp", ExitOnError)
		assert.Equal(t, "myapp", fs.Name())
	})

	t.Run("env prefix", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.EnvPrefix("MYAPP_")
		assert.Equal(t, "MYAPP_", fs.envPrefix)
	})

	t.Run("version", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.Version("v1.2.3")
		assert.Equal(t, "v1.2.3", fs.versionString)
		assert.True(t, fs.enableVer)
	})

	t.Run("title", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.Title("Usage:")
		assert.Equal(t, "Usage:", fs.title)
	})

	t.Run("description", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.Description("This is a test app.")
		assert.Equal(t, "This is a test app.", fs.desc)
	})

	t.Run("note", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.Note("Footer note")
		assert.Equal(t, "Footer note", fs.notes)
	})

	t.Run("disable help", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.DisableHelp()
		assert.False(t, fs.enableHelp)
	})

	t.Run("disable version", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.DisableVersion()
		assert.False(t, fs.enableVer)
		assert.Equal(t, "", fs.versionString)
	})

	t.Run("sorted", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.Sorted(true)
		assert.True(t, fs.sortFlags)
	})

	t.Run("set output and output", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		buf := new(bytes.Buffer)
		fs.SetOutput(buf)
		assert.Equal(t, buf, fs.Output())
	})

	t.Run("ignore invalid env", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.IgnoreInvalidEnv(true)
		assert.True(t, fs.ignoreInvalidEnv)
	})

	t.Run("set get env fn", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.SetGetEnvFn(func(key string) string {
			if key == "X" {
				return "42"
			}
			return ""
		})
		assert.Equal(t, "42", fs.getEnv("X"))
	})

	t.Run("globaldelimiter", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.Globaldelimiter(";")
		assert.Equal(t, ";", fs.defaultDelimiter)
	})

	t.Run("require positional", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.RequirePositional(2)
		assert.Equal(t, 2, fs.requiredPositional)
	})

	t.Run("args and arg", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.positional = []string{"one", "two"}

		assert.Equal(t, []string{"one", "two"}, fs.Args())

		val, ok := fs.Arg(1)
		assert.True(t, ok)
		assert.Equal(t, "two", val)

		_, ok = fs.Arg(5)
		assert.False(t, ok)
	})

	t.Run("usage print mode", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.UsagePrintMode(PrintBoth)
		assert.Equal(t, PrintBoth, fs.usagePrintMode)
	})

	t.Run("description max len", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.DescriptionMaxLen(50)
		assert.Equal(t, 50, fs.descMaxLen)
	})

	t.Run("description indent", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		fs.DescriptionIndent(10)
		assert.Equal(t, 10, fs.descIndent)
	})

	t.Run("get and must get", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		ptr := fs.String("name", "default", "name of user").Value()

		v, err := fs.Get("name")
		assert.NoError(t, err)
		assert.Equal(t, *ptr, v)

		assert.Panics(t, func() { fs.MustGet("notfound") })
	})

	t.Run("MustGet success", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		ptr := fs.String("name", "default", "name of user").Value()

		v := fs.MustGet("name")
		assert.Equal(t, *ptr, v)
	})

	t.Run("generic", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		name := fs.String("name", "bob", "the name").Value()

		v, err := GetAs[string](fs, "name")
		assert.NoError(t, err)
		assert.Equal(t, *name, v)

		_, err = GetAs[int](fs, "unknown")
		assert.Error(t, err)
	})
}
