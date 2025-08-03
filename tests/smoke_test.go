package tinyflags_test

import (
	"fmt"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmoke_ParseArgs(t *testing.T) {
	t.Parallel()

	t.Run("smoke counter", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)

		verbose := fs.Counter("verbose", 0, "Enable verbose mode").
			Max(4).
			Short("v").
			Placeholder("NUM").
			Value()

		err := fs.Parse([]string{
			"--verbose",
			"-v",
		})
		require.NoError(t, err)
		assert.Equal(t, 2, *verbose)
	})

	t.Run("smoke counter max reached", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)

		verbose := fs.Counter("verbose", 0, "Enable verbose mode").
			Max(2).
			Value()

		err := fs.Parse([]string{
			"--verbose",
			"--verbose",
			"--verbose",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "must not be greater than 2")
		assert.Equal(t, 2, *verbose)
	})

	t.Run("smoke scalar", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.0.0")
		fs.EnvPrefix("MY_SUPER_APP")

		host := fs.String("host", "Host to use", "").
			Choices("alpha", "beta", "gamma").
			Value()

		err := fs.Parse([]string{
			"--host=beta",
		})
		require.NoError(t, err)

		assert.Equal(t, "beta", *host)
	})

	t.Run("smoke slice", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.0.0")

		host := fs.StringSlice("host", []string{}, "Host to use").
			Choices("alpha", "beta", "gamma").
			Delimiter("|").
			Value()

		err := fs.Parse([]string{
			"--host=alpha",
			"--host=beta|gamma",
		})

		require.NoError(t, err)
		assert.Equal(t, []string{"alpha", "beta", "gamma"}, *host)
	})

	t.Run("smoke dynamic", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.0.0")

		http := fs.DynamicGroup("http")
		addr := http.String("address", "", "API address")

		err := fs.Parse([]string{
			"--http.alpha.address=127.0.0.1",
			"--http.beta.address=10.0.0.1",
		})
		require.NoError(t, err)

		for _, id := range http.Instances() {
			fmt.Println("instance:", id)

			addrVal := addr.MustGet(id)

			fmt.Println("  address:", addrVal)
		}
	})

	t.Run("smoke dyn slice", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.0.0")

		http := fs.DynamicGroup("http")
		addrs := http.StringSlice("addresses", []string{}, "API address")

		err := fs.Parse([]string{
			"--http.alpha.addresses=127.0.0.1",
			"--http.alpha.addresses=10.0.0.2",
			"--http.alpha.addresses=10.0.0.2,10.0.0.3",
		})
		require.NoError(t, err)

		for _, id := range http.Instances() {
			fmt.Println("instance:", id)

			addrVal := addrs.MustGet(id)

			fmt.Println("  addresses:", addrVal)
		}
	})

	t.Run("smoke basic help", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.0.0")
		err := fs.Parse([]string{
			"--help",
		})
		require.Error(t, err)
		assert.EqualError(t, err, `Usage: app [flags]
Flags:
    -h, --help     show help (Default: false)
        --version  show version (Default: false)
`)
	})

	t.Run("smoke custom help", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("custom.exe", tinyflags.ContinueOnError)
		fs.Usage = func() {
			fs.PrintUsage(fs.Output(), tinyflags.PrintNone)
		}
		fs.Version("1.0.0")
		err := fs.Parse([]string{
			"--help",
		})
		require.Error(t, err)
		assert.EqualError(t, err, "Usage: custom.exe\n")
	})

	t.Run("smoke example", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.0.0")
		fs.EnvPrefix("MY_SUPER_APP")

		host := fs.String("host", "Host to use", "").
			Choices("alpha", "beta", "gamma").
			Value()

		list := fs.StringSlice("list", []string{}, "my list").
			Delimiter("|").
			Choices("a", "b", "c", "d").
			Value()

		http := fs.DynamicGroup("http")
		addr := http.String("address", "", "API address").
			Validate(func(s string) error {
				if s == "localhost" {
					return fmt.Errorf("cannot use localhost")
				}
				return nil
			})
		port := http.Int("port", 0, "API port")
		sl := http.StringSlice("list", []string{}, "List of values")

		verbose := fs.Bool("verbose", false, "Enable verbose mode").
			Strict().
			Short("v").
			Value()

		debug := fs.Bool("debug", true, "Enable debug mode").Value()

		err := fs.Parse([]string{
			"--http.alpha.address=127.0.0.1",
			"--http.alpha.port", "8080",
			"--http.beta.address=10.0.0.1",
			"--http.beta.list=a|b|c",
			"--host=beta",
			"--list=a|b|c", "--list", "d",
			"-vtrue",
			"--debug",
		})
		require.NoError(t, err)

		assert.Equal(t, "beta", *host)
		assert.Equal(t, []string{"a", "b", "c", "d"}, *list)
		assert.Equal(t, true, *verbose)
		assert.Equal(t, true, *debug)

		for _, id := range http.Instances() {
			fmt.Println("instance:", id)

			addrVal := addr.MustGet(id)

			// Use Get instead of MustGet for optional fields
			portVal, _ := port.Get(id)
			listVal, _ := sl.Get(id)

			fmt.Println("  address:", addrVal)
			fmt.Println("  port:", portVal)
			fmt.Println("  list:", listVal)
		}
	})

	t.Run("smoke one of group", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)

		debug := fs.Bool("debug", true, "Enable debug mode").OneOfGroup("db").Value()
		noDebug := fs.Bool("no-debug", false, "Disable debug mode").OneOfGroup("db").Value()
		fs.GetOneOfGroup("db").Title("Debug Options").Required()

		err := fs.Parse([]string{
			"--debug",
		})
		require.NoError(t, err)
		assert.Equal(t, true, *debug)
		assert.Equal(t, false, *noDebug)
	})

	t.Run("smoke choice error", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		list := fs.StringSlice("list", []string{}, "my list").
			Delimiter("|").
			Choices("a", "b", "c").
			Value()

		err := fs.Parse([]string{
			"--list=a|b|c", "--list", "d",
		})
		require.EqualError(t, err, "invalid value for flag --list: invalid value \"d\": must be one of: a, b, c.")
		assert.Equal(t, []string{"a", "b", "c"}, *list) // defaults
	})

	t.Run("smoke validator error", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Int("port", 8080, "API port").Validate(func(p int) error {
			if p < 1000 {
				return fmt.Errorf("port must be greater than 1000")
			}
			return nil
		})

		err := fs.Parse([]string{
			"--port", "80",
		})
		require.Error(t, err)
		require.EqualError(t, err, "invalid value for flag --port: port must be greater than 1000.")
	})

	t.Run("smoke dyn validator error", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		http := fs.DynamicGroup("http")
		http.String("ip", "127.0.0.1", "ip address").Choices("127.0.0.1", "10.0.0.1")

		http.String("address", "", "API address").
			Validate(func(s string) error {
				if s == "localhost" {
					return fmt.Errorf("cannot use localhost")
				}
				return nil
			})

		http.StringSlice("list", []string{}, "list of values").
			Choices("a", "b", "c").
			Required()

		http.Int("port", 0, "API port")

		err := fs.Parse([]string{
			"--http.alpha.address=127.0.0.1",
			"--http.alpha.port", "8080",
			"--http.beta.address=localhost",
			"--http.beta.list=a|b|c",
		})
		require.Error(t, err)
		require.EqualError(t, err, "invalid value for dynamic flag --http.beta.address: cannot use localhost")
	})

	t.Run("smoke dyn choices", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		http := fs.DynamicGroup("http")
		addr := http.String("address", "127.0.0.1", "API address").
			Choices("127.0.0.1", "10.0.0.1")

		err := fs.Parse([]string{
			"--http.alpha.address=127.0.0.1",
			"--http.beta.address=10.0.0.1",
		})
		require.NoError(t, err)
		addrs := map[string]string{
			"alpha": "127.0.0.1",
			"beta":  "10.0.0.1",
		}
		assert.Equal(t, addrs, addr.Values())
	})

	t.Run("smoke dyn values", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		http := fs.DynamicGroup("http")
		addr := http.String("address", "", "API address").
			Validate(func(s string) error {
				if s == "localhost" {
					return fmt.Errorf("cannot use localhost")
				}
				return nil
			})
		port := http.Int("port", 0, "API port")
		enabled := http.Bool("enabled", true, "Enable service").Strict()

		verbose := http.Bool("verbose", true, "Enable verbose mode")

		err := fs.Parse([]string{
			"--http.alpha.enabled=true",
			"--http.alpha.address=127.0.0.1",
			"--http.alpha.port", "8080",
			"--http.alpha.verbose",
			"--http.beta.address=10.0.0.1",
			"--http.beta.enabled=false",
		})
		require.NoError(t, err)
		addrs := map[string]string{
			"alpha": "127.0.0.1",
			"beta":  "10.0.0.1",
		}
		assert.Equal(t, addrs, addr.Values())

		ports := map[string]int{
			"alpha": 8080,
		}
		assert.Equal(t, ports, port.Values())

		enableds := map[string]bool{
			"alpha": true,
			"beta":  false,
		}
		assert.Equal(t, enableds, enabled.Values())

		verboses := map[string]bool{
			"alpha": true,
		}
		assert.Equal(t, verboses, verbose.Values())
	})

	t.Run("smoke dyn help", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.SetDynamicUsageIndent(8)
		http := fs.DynamicGroup("http")
		http.String("address", "", "API address").
			Validate(func(s string) error {
				if s == "localhost" {
					return fmt.Errorf("cannot use localhost")
				}
				return nil
			})
		http.Int("port", 0, "API port")
		http.Bool("enabled", true, "Enable service").Strict()
		http.Bool("verbose", true, "Enable verbose mode")

		err := fs.Parse([]string{
			"--http.alpha.enabled=true",
			"--http.alpha.address=127.0.0.1",
			"--http.alpha.port", "8080",
			"--http.alpha.verbose",
			"--http.beta.address=10.0.0.1",
			"--http.beta.enabled=false",
			"--help",
		})
		require.Error(t, err)
		assert.EqualError(t, err, `Usage: app [flags]
Flags:
    -h, --help  show help (Default: false)
        --http.<ID>.address ADDRESS       API address
        --http.<ID>.port PORT             API port (Default: 0)
        --http.<ID>.enabled <true|false>  Enable service (Allowed: true, false) (Default: true)
        --http.<ID>.verbose               Enable verbose mode (Default: true)
`)
	})
}
