package tinyflags_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagSet(t *testing.T) {
	t.Parallel()

	t.Run("basic usage", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		port := tf.Int("port", 8080, "port to use")
		err := tf.Parse([]string{"--port=1234"})
		require.NoError(t, err)
		assert.Equal(t, 1234, *port.Value())

		p, err := tf.Get("port")
		require.NoError(t, err)
		assert.Equal(t, 1234, p)
	})

	t.Run("basic usage - custom delimiter", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		port := tf.IntSlice("port", []int{8080}, "port to use").Delimiter(":")
		err := tf.Parse([]string{"--port=1234", "--port=5678", "--port=9012:3456"})
		require.NoError(t, err)
		assert.Equal(t, []int{1234, 5678, 9012, 3456}, *port.Value())

		p, err := tf.Get("port")
		require.NoError(t, err)
		assert.Equal(t, []int{1234, 5678, 9012, 3456}, p)
	})

	t.Run("version flag", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		tf.Version("v1.2.3")

		err := tf.Parse([]string{"--version"})
		var versionErr *tinyflags.VersionRequested
		require.ErrorAs(t, err, &versionErr)
		assert.Equal(t, "v1.2.3", versionErr.Version)
	})

	t.Run("help flag", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		tf.Version("v1.2.3")
		buf := &bytes.Buffer{}
		tf.SetOutput(buf)
		tf.Int("port", 8080, "port to use")
		err := tf.Parse([]string{"--help"})
		var helpErr *tinyflags.HelpRequested
		require.ErrorAs(t, err, &helpErr)
		assert.EqualError(t, err, "Usage: test.exe --port PORT --version -h|--help\n      --port PORT                        port to use (Default: 8080)\n      --version                          show version\n  -h, --help                             show help\n")
	})

	t.Run("disable help", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ContinueOnError)
		tf.DisableHelp()
		err := tf.Parse([]string{"--help"})
		require.Error(t, err)
		assert.EqualError(t, err, "unknown flag: --help")
	})

	t.Run("disable version", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ContinueOnError)
		tf.DisableVersion()
		err := tf.Parse([]string{"--version"})
		require.Error(t, err)
		assert.EqualError(t, err, "unknown flag: --version")
	})

	t.Run("env override", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		tf.EnvPrefix("MYAPP")
		tf.SetGetEnvFn(func(key string) string {
			if key == "MYAPP_PORT" {
				return "5050"
			}
			return ""
		})

		port := tf.Int("port", 8080, "port to use")
		err := tf.Parse([]string{})
		require.NoError(t, err)
		assert.Equal(t, 5050, *port.Value())
	})

	t.Run("positional args", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		tf.Int("port", 8080, "port to use")
		err := tf.Parse([]string{"--port=9000", "arg1", "arg2"})
		require.NoError(t, err)

		assert.Equal(t, []string{"arg1", "arg2"}, tf.Args())

		arg0, ok := tf.Arg(0)
		assert.True(t, ok)
		assert.Equal(t, "arg1", arg0)

		arg1, ok := tf.Arg(1)
		assert.True(t, ok)
		assert.Equal(t, "arg2", arg1)

		_, ok = tf.Arg(2)
		assert.False(t, ok)
	})

	t.Run("require positional", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ContinueOnError)
		tf.RequirePositional(2)
		tf.Int("port", 8080, "port to use")
		err := tf.Parse([]string{"--port=9999", "onlyone"})
		require.Error(t, err)
		assert.EqualError(t, err, "expected at least 2 positional arguments, got 1")
	})

	t.Run("invalid flag type", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ContinueOnError)
		tf.Int("port", 8080, "port to use")
		err := tf.Parse([]string{"--port=notanint"})
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid value for flag --port: strconv.Atoi: parsing \"notanint\": invalid syntax")
	})

	t.Run("get unknown", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		_, err := tf.Get("notfound")
		assert.Error(t, err)
		assert.EqualError(t, err, "flag \"notfound\" not found")
	})

	t.Run("must get panic", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for missing flag")
			}
		}()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.PanicOnError)
		tf.MustGet("missing")
	})

	t.Run("title description note", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		tf.Version("v1.2.3")
		tf.Title("Cool Tool")
		tf.Description("Does things.")
		tf.Note("Thanks for using.")
		tf.Int("port", 8080, "the port")

		var out bytes.Buffer
		tf.SetOutput(&out)

		err := tf.Parse([]string{"--help"})
		var help *tinyflags.HelpRequested
		require.ErrorAs(t, err, &help)

		assert.EqualError(t, help, "Cool Tool\n\nUsage: test.exe --port PORT --version -h|--help\nDoes things.\n\n  --port PORT  the port (Default: 8080)\n  --version    show version\n  -h, --help   show help\n\nThanks for using.\n")
	})

	t.Run("sorted is disabled", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		tf.Version("v1.2.3")
		tf.Sorted(false)

		tf.String("zeta", "x", "zeta flag")
		tf.Bool("gamma", false, "gamma flag")
		tf.Int("beta", 2, "beta flag")

		var out bytes.Buffer
		tf.SetOutput(&out)

		err := tf.Parse([]string{"--help"})
		assert.EqualError(t, err, "Usage: test.exe --zeta ZETA --gamma --beta BETA --version -h|--help\n  --zeta ZETA  zeta flag (Default: x)\n  --gamma      gamma flag\n  --beta BETA  beta flag (Default: 2)\n  --version    show version\n  -h, --help   show help\n")
	})

	t.Run("sorted is enabled", func(t *testing.T) {
		t.Parallel()

		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		tf.Version("v1.2.3")
		tf.Sorted(true)

		tf.Bool("zeta", false, "zeta flag")
		tf.Int("beta", 2, "beta flag").Choices(2, 3, 4)
		tf.String("alpha", "x", "alpha flag")

		var out bytes.Buffer
		tf.SetOutput(&out)

		err := tf.Parse([]string{"--help"})
		assert.EqualError(t, err, `Usage: test.exe --alpha ALPHA --beta <2|3|4> --zeta --version -h|--help
  --alpha ALPHA   alpha flag (Default: x)
  --beta <2|3|4>  beta flag (Default: 2)
  -h, --help      show help
  --version       show version
  --zeta          zeta flag
`)
	})

	t.Run("Custom Usage", func(t *testing.T) {
		tf := tinyflags.NewFlagSet("test.exe", tinyflags.ExitOnError)
		tf.Version("v1.2.3")
		tf.Sorted(true)

		tf.Bool("zeta", false, "zeta flag")
		tf.Int("beta", 2, "beta flag")
		tf.String("alpha", "x", "alpha flag")

		var out bytes.Buffer
		tf.SetOutput(&out)
		tf.Usage = func() {
			fmt.Fprintln(tf.Output(), "Custom Usage") // nolint:errcheck
		}
		err := tf.Parse([]string{"--help"})
		assert.EqualError(t, err, "Custom Usage\n")
	})
}

func TestFlagSet_GetAs(t *testing.T) {
	t.Parallel()

	tf := tinyflags.NewFlagSet("test.exe", tinyflags.PanicOnError)
	tf.Int("count", 3, "how many")
	_ = tf.Parse([]string{"--count=9"})

	count, err := tinyflags.GetAs[int](tf, "count")
	require.NoError(t, err)
	assert.Equal(t, 9, count)
}
