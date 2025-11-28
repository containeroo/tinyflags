package engine

import "testing"

func TestWrapText(t *testing.T) {
	t.Parallel()

	t.Run("emptyString", func(t *testing.T) {
		t.Parallel()
		if got := wrapText("", 10); got != "" {
			t.Fatalf("expected empty, got %q", got)
		}
	})

	t.Run("wrapsAtWidth", func(t *testing.T) {
		t.Parallel()
		text := "a bb ccc dddd"
		got := wrapText(text, 6)
		if got != "a bb\nccc\ndddd" {
			t.Fatalf("unexpected wrap: %q", got)
		}
	})
}
