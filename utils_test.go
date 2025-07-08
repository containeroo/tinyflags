package tinyflags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAllowed(t *testing.T) {
	t.Parallel()

	t.Run("with custom format function", func(t *testing.T) {
		t.Parallel()
		input := []int{1, 2, 3}
		format := func(i int) string {
			return "v" + string(rune('A'+i))
		}
		got := formatAllowed(input, format)
		assert.Equal(t, "vB, vC, vD", got)
	})

	t.Run("with nil format function", func(t *testing.T) {
		t.Parallel()
		input := []string{"a", "b", "c"}
		got := formatAllowed(input, nil)
		assert.Equal(t, "[a b c]", got) // default fmt.Sprintf fallback
	})

	t.Run("with empty input", func(t *testing.T) {
		t.Parallel()
		input := []string{}
		got := formatAllowed(input, func(s string) string { return s })
		assert.Equal(t, "", got)
	})
}

func TestPluralSuffix(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "", pluralSuffix(1))
	assert.Equal(t, "s", pluralSuffix(0))
	assert.Equal(t, "s", pluralSuffix(2))
	assert.Equal(t, "s", pluralSuffix(-1))
}
