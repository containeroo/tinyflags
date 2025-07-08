package tinyflags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapText(t *testing.T) {
	t.Parallel()
	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		result := wrapText("", 10)
		assert.Equal(t, "", result)
	})

	t.Run("single short word", func(t *testing.T) {
		t.Parallel()
		result := wrapText("hello", 10)
		assert.Equal(t, "hello", result)
	})

	t.Run("single long word", func(t *testing.T) {
		t.Parallel()
		result := wrapText("supercalifragilisticexpialidocious", 10)
		assert.Equal(t, "supercalifragilisticexpialidocious", result)
	})

	t.Run("two words fit in line", func(t *testing.T) {
		t.Parallel()
		result := wrapText("hello world", 20)
		assert.Equal(t, "hello world", result)
	})

	t.Run("two words exceed width", func(t *testing.T) {
		t.Parallel()
		result := wrapText("hello wonderful", 10)
		assert.Equal(t, "hello\nwonderful", result)
	})

	t.Run("wrap multiple lines", func(t *testing.T) {
		t.Parallel()
		result := wrapText("the quick brown fox jumps over the lazy dog", 15)
		assert.Equal(t, "the quick brown\nfox jumps over\nthe lazy dog", result)
	})

	t.Run("excessive spaces", func(t *testing.T) {
		t.Parallel()
		result := wrapText("   lots   of   space   here   ", 4)
		expected := `lots
of
space
here`
		assert.Equal(t, expected, result)
	})

	t.Run("extra line width", func(t *testing.T) {
		t.Parallel()
		result := wrapText("12345 67890", 11)
		assert.Equal(t, "12345 67890", result)
	})

	t.Run("last word is long", func(t *testing.T) {
		t.Parallel()
		result := wrapText("foo bar supercalifragilisticexpialidocious", 10)
		assert.Equal(t, "foo bar\nsupercalifragilisticexpialidocious", result)
	})
}
