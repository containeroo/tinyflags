package tinyflags_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmoke_ParseSimpleStringAndEnv(t *testing.T) {
	t.Parallel()

	getEnv := func(key string) string {
		m := map[string]string{
			"APP_NAME": "envname",
		}
		return m[key]
	}

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.SetGetEnvFn(getEnv)
	var output bytes.Buffer
	fs.SetOutput(&output)
	fs.EnvPrefix("APP")

	name := fs.String("name", "default", "your name").Value()
	err := fs.Parse([]string{})

	assert.NoError(t, err)
	assert.Equal(t, "envname", *name)
}

func TestSmoke_ParseArgs(t *testing.T) {
	t.Parallel()

	t.Run("IntP", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		age := fs.IntP("age", "a", 0, "your age").Value()
		err := fs.Parse([]string{"--age=42"})

		assert.NoError(t, err)
		assert.Equal(t, 42, *age)
	})

	t.Run("StingSlice", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		list := fs.StringSlice("list", []string{"a", "b"}, "your list").Value()
		err := fs.Parse([]string{"--list", "alpha", "--list=beta", "--list=gamma,delta"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"alpha", "beta", "gamma", "delta"}, *list)
	})
}

func TestSmoke_HelpRequested(t *testing.T) {
	t.Parallel()
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.Version("v1.0.0")
	buf := new(bytes.Buffer)
	fs.SetOutput(buf)

	err := fs.Parse([]string{"--help"})
	help := `Usage: app [flags]
Flags:
      --version                          show version
  -h, --help                             show help
`
	assert.Equal(t, help, err.Error())
	assert.EqualError(t, err, help)
}

func TestSmoke_RequiredPositional(t *testing.T) {
	t.Parallel()
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.RequirePositional(1)

	err := fs.Parse([]string{})
	assert.Error(t, err)

	err = fs.Parse([]string{"foo"})
	assert.NoError(t, err)
	assert.Equal(t, "foo", fs.Args()[0])
}

func TestSmoke(t *testing.T) {
	t.Parallel()

	t.Run("Parse simple string and env override", func(t *testing.T) {
		t.Parallel()

		getEnv := func(key string) string {
			m := map[string]string{
				"APP_NAME": "envname",
			}
			return m[key]
		}

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.SetGetEnvFn(getEnv)
		fs.EnvPrefix("APP")
		var output bytes.Buffer
		fs.SetOutput(&output)

		name := fs.String("name", "default", "your name").Value()
		err := fs.Parse([]string{})

		require.NoError(t, err)
		assert.Equal(t, "envname", *name)
	})

	t.Run("Parse simple int flag", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		age := fs.IntP("age", "a", 0, "your age").Value()
		err := fs.Parse([]string{"--age=42"})

		require.NoError(t, err)
		assert.Equal(t, 42, *age)
	})

	t.Run("Help requested output", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("v1.0.0")
		var buf bytes.Buffer
		fs.SetOutput(&buf)

		err := fs.Parse([]string{"--help"})

		help := `Usage: app [flags]
Flags:
      --version                          show version
  -h, --help                             show help
`
		assert.EqualError(t, err, help)
		assert.Equal(t, help, err.Error())
	})

	t.Run("Required positional arguments", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.RequirePositional(1)

		err := fs.Parse([]string{})
		assert.EqualError(t, err, "expected at least 1 positional argument, got 0")

		err = fs.Parse([]string{"foo"})
		require.NoError(t, err)
		assert.Equal(t, "foo", fs.Args()[0])
	})

	t.Run("Invalid flag value", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Int("port", 8080, "port to use")

		err := fs.Parse([]string{"--port=notanint"})
		assert.EqualError(t, err, `invalid value for flag --port: strconv.Atoi: parsing "notanint": invalid syntax`)
	})

	t.Run("Get unknown flag", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		_, err := fs.Get("notfound")
		assert.EqualError(t, err, `flag "notfound" not found`)
	})

	t.Run("MustGet panics", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic from MustGet")
			}
		}()

		fs := tinyflags.NewFlagSet("app", tinyflags.PanicOnError)
		fs.MustGet("missing")
	})

	t.Run("Custom title, description, and note in help", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		fs.DescriptionIndent(18)
		fs.Title("Cool Title")
		fs.Description("My Flags:")
		fs.Note("Thanks for using.")
		fs.Int("port", 8080, "the port").Value()
		fs.Version("v1.2.3")
		var out bytes.Buffer
		fs.SetOutput(&out)

		err := fs.Parse([]string{"--help"})
		var help *tinyflags.HelpRequested
		require.ErrorAs(t, err, &help)

		expected := `Usage: test.exe [flags]
Cool Title
My Flags:
      --port PORT  the port (Default: 8080)
      --version    show version
  -h, --help       show help
Thanks for using.
`
		assert.EqualError(t, err, expected)
	})

	t.Run("Custom usage override", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
		fs.Usage = func() {
			fmt.Fprintln(fs.Output(), "Custom Usage") // nolint:errcheck
		}
		var out bytes.Buffer
		fs.SetOutput(&out)

		err := fs.Parse([]string{"--help"})
		assert.EqualError(t, err, "Custom Usage\n")
	})
}
