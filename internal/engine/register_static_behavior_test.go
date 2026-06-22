package engine

import (
	"strconv"
	"testing"

	"github.com/containeroo/tinyflags/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegisterStaticScalarInt verifies static scalar registration behavior.
func TestRegisterStaticScalarInt(t *testing.T) {
	t.Parallel()

	t.Run("defaultIsApplied", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)

		var value int
		RegisterStaticScalar(fs, &value, "num", "usage", 7, strconv.Atoi, strconv.Itoa)

		assert.Equal(t, 7, value)
	})

	t.Run("parsesInput", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)

		num := new(int)
		RegisterStaticScalar(fs, num, "num", "usage", 0, strconv.Atoi, strconv.Itoa)

		require.NoError(t, fs.Parse([]string{"--num=5"}))
		assert.Equal(t, 5, *num)
	})
}

// TestRegisterStaticSliceString verifies static slice registration behavior.
func TestRegisterStaticSliceString(t *testing.T) {
	t.Parallel()

	t.Run("defaultIsApplied", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)

		var names []string
		RegisterStaticSlice(fs, &names, "names", "usage", []string{"a", "b"}, utils.ParseString, utils.FormatString, ";", false)

		assert.Equal(t, []string{"a", "b"}, names)
	})

	t.Run("parsesInputWithDelimiter", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)

		names := new([]string)
		RegisterStaticSlice(fs, names, "names", "usage", nil, utils.ParseString, utils.FormatString, ";", true)

		require.NoError(t, fs.Parse([]string{"--names=a; b"}))

		want := []string{"a", "b"}
		assert.Equal(t, want, *names)
	})
}
