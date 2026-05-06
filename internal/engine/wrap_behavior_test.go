package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWrapText verifies usage text wrapping behavior.
func TestWrapText(t *testing.T) {
	t.Parallel()

	t.Run("emptyString", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "", wrapText("", 10))
	})

	t.Run("wrapsAtWidth", func(t *testing.T) {
		t.Parallel()
		text := "a bb ccc dddd"
		got := wrapText(text, 6)
		assert.Equal(t, "a bb\nccc\ndddd", got)
	})

	t.Run("preservesNewlines", func(t *testing.T) {
		t.Parallel()
		text := "first line\nsecond line"
		assert.Equal(t, text, wrapText(text, 20))
	})
}
