package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMinPositional verifies minimum positional argument enforcement.
func TestMinPositional(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.MinPositional(2)

	require.NoError(t, fs.Parse([]string{"one", "two", "three"}))
	assert.Equal(t, []string{"one", "two", "three"}, fs.Args())

	err := fs.Parse([]string{"one"})
	require.EqualError(t, err, "expected at least 2 positional arguments, got 1")
}

// TestExactPositional verifies exact positional argument enforcement.
func TestExactPositional(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.ExactPositional(2)

	require.NoError(t, fs.Parse([]string{"one", "two"}))

	err := fs.Parse([]string{"one"})
	require.EqualError(t, err, "expected exactly 2 positional arguments, got 1")

	err = fs.Parse([]string{"one", "two", "three"})
	require.EqualError(t, err, "expected exactly 2 positional arguments, got 3")
}

// TestMaxPositional verifies maximum positional argument enforcement.
func TestMaxPositional(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.MaxPositional(2)

	require.NoError(t, fs.Parse([]string{"one", "two"}))
	require.NoError(t, fs.Parse([]string{"one"}))
	require.NoError(t, fs.Parse(nil))

	err := fs.Parse([]string{"one", "two", "three"})
	require.EqualError(t, err, "expected at most 2 positional arguments, got 3")
}

// TestNoPositional verifies positional arguments can be rejected entirely.
func TestNoPositional(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.NoPositional()

	require.NoError(t, fs.Parse(nil))
	require.EqualError(t, fs.Parse([]string{"unexpected"}), "accepts no positional arguments")
}

// TestPositionalRange verifies minimum and maximum constraints compose into a range.
func TestPositionalRange(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.MinPositional(1)
	fs.MaxPositional(2)

	require.NoError(t, fs.Parse([]string{"one"}))
	require.NoError(t, fs.Parse([]string{"one", "two"}))
	require.EqualError(t, fs.Parse(nil), "expected at least 1 positional argument, got 0")
	require.EqualError(t, fs.Parse([]string{"one", "two", "three"}), "expected at most 2 positional arguments, got 3")
}

// TestRequirePositionalCompatibility verifies the deprecated API keeps its minimum semantics.
func TestRequirePositionalCompatibility(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.RequirePositional(2) //nolint:staticcheck // Backward-compatibility coverage.

	require.NoError(t, fs.Parse([]string{"one", "two", "three"}))
	require.EqualError(t, fs.Parse([]string{"one"}), "expected at least 2 positional arguments, got 1")
}

// TestPositionalConstraintsDoNotBlockHelpOrVersion verifies built-ins bypass positional counts.
func TestPositionalConstraintsDoNotBlockHelpOrVersion(t *testing.T) {
	t.Parallel()

	t.Run("help", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.ExactPositional(2)
		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.True(t, tinyflags.IsHelpRequested(err))
	})

	t.Run("version", func(t *testing.T) {
		t.Parallel()
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.MaxPositional(0)
		fs.Version("1.0.0")
		err := fs.Parse([]string{"--version"})
		require.Error(t, err)
		assert.True(t, tinyflags.IsVersionRequested(err))
	})
}
