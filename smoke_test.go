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

	t.Run("smoke", func(t *testing.T) {
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
		addr := http.String("address", "API address")
		port := http.Int("port", "API port")
		sl := http.StringSlice("list", "List of values")

		verbose := fs.BoolP("verbose", "v", false, "Enable verbose mode").
			Strict().
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

	t.Run("smoke error", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		list := fs.StringSlice("list", []string{}, "my list").
			Delimiter("|").
			Choices("a", "b", "c").
			Value()

		err := fs.Parse([]string{
			"--list=a|b|c", "--list", "d",
		})
		require.EqualError(t, err, "invalid value for flag --list: got invalid value \"d\": must be one of [a, b, c].")
		assert.Equal(t, []string{"a", "b", "c"}, *list) // defaults
	})
}
