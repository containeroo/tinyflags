package tinyflags_test

import (
	"fmt"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
)

func TestSmoke_ParseArgs(t *testing.T) {
	t.Parallel()

	t.Run("smoke", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)

		fs.String("host", "Host to use", "").Choices("alpha", "beta", "gamma")

		http := fs.DynamicGroup("http")
		addr := http.String("address", "API address")
		port := http.Int("port", "API port")

		err := fs.Parse([]string{
			"--http.alpha.address=127.0.0.1",
			"--http.alpha.port", "8080",
			"--http.beta.address=10.0.0.1",
			"--host=beta",
		})
		assert.NoError(t, err)

		for _, id := range http.Instances() {
			fmt.Println("instance:", id)

			addrVal := addr.MustGet(id)

			// Use Get instead of MustGet for optional fields
			portVal, _ := port.Get(id)

			fmt.Println("  address:", addrVal)
			fmt.Println("  port:", portVal)
		}
	})
}
