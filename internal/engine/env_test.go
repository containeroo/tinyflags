package engine

import (
	"strings"
	"testing"
)

func TestNewReplacerEnvKeyFunc(t *testing.T) {
	t.Parallel()

	t.Run("emptyPrefixReturnsEmpty", func(t *testing.T) {
		t.Parallel()

		fn := NewReplacerEnvKeyFunc(strings.NewReplacer("-", "_"), true)
		if got := fn("", "flag-name"); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})

	t.Run("replacesAndUppercases", func(t *testing.T) {
		t.Parallel()

		fn := NewReplacerEnvKeyFunc(strings.NewReplacer("-", "_", ".", "_"), true)
		if got := fn("app", "flag-name.more"); got != "APP_FLAG_NAME_MORE" {
			t.Fatalf("unexpected env key: %q", got)
		}
	})
}
