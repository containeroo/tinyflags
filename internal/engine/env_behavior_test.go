package engine

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReplacerEnvKeyFunc(t *testing.T) {
	t.Parallel()

	t.Run("emptyPrefixReturnsEmpty", func(t *testing.T) {
		t.Parallel()

		fn := NewReplacerEnvKeyFunc(strings.NewReplacer("-", "_"), true)
		assert.Equal(t, "", fn("", "flag-name"))
	})

	t.Run("replacesAndUppercases", func(t *testing.T) {
		t.Parallel()

		fn := NewReplacerEnvKeyFunc(strings.NewReplacer("-", "_", ".", "_"), true)
		assert.Equal(t, "APP_FLAG_NAME_MORE", fn("app", "flag-name.more"))
	})
}
