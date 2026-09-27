package tinyflags_test

import (
	"math"
	"testing"
	"time"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPositive verifies validation of values that must be greater than zero.
func TestPositive(t *testing.T) {
	t.Parallel()

	validateInt := tinyflags.Positive[int]()
	assert.NoError(t, validateInt(1))
	assert.EqualError(t, validateInt(0), "must be greater than zero")
	assert.EqualError(t, validateInt(-1), "must be greater than zero")

	validateDuration := tinyflags.Positive[time.Duration]()
	assert.NoError(t, validateDuration(time.Second))
	assert.EqualError(t, validateDuration(0), "must be greater than zero")

	validateFloat := tinyflags.Positive[float64]()
	assert.EqualError(t, validateFloat(math.NaN()), "must be greater than zero")
}

// TestNonNegative verifies validation of values that may be zero but not negative.
func TestNonNegative(t *testing.T) {
	t.Parallel()

	validate := tinyflags.NonNegative[int]()
	assert.NoError(t, validate(1))
	assert.NoError(t, validate(0))
	assert.EqualError(t, validate(-1), "must not be negative")

	validateFloat := tinyflags.NonNegative[float64]()
	assert.EqualError(t, validateFloat(math.NaN()), "must not be negative")
}

// TestAtLeast verifies inclusive lower-bound validation.
func TestAtLeast(t *testing.T) {
	t.Parallel()

	validate := tinyflags.AtLeast(3)
	assert.NoError(t, validate(3))
	assert.NoError(t, validate(4))
	assert.EqualError(t, validate(2), "must be at least 3")

	validateDuration := tinyflags.AtLeast(time.Second)
	assert.NoError(t, validateDuration(time.Second))
	assert.EqualError(t, validateDuration(500*time.Millisecond), "must be at least 1s")
}

// TestAtMost verifies inclusive upper-bound validation.
func TestAtMost(t *testing.T) {
	t.Parallel()

	validate := tinyflags.AtMost(3)
	assert.NoError(t, validate(3))
	assert.NoError(t, validate(2))
	assert.EqualError(t, validate(4), "must be at most 3")
}

// TestBetween verifies inclusive range validation and invalid range handling.
func TestBetween(t *testing.T) {
	t.Parallel()

	t.Run("inclusiveRange", func(t *testing.T) {
		t.Parallel()

		validate := tinyflags.Between(1, 3)
		assert.NoError(t, validate(1))
		assert.NoError(t, validate(2))
		assert.NoError(t, validate(3))
		assert.EqualError(t, validate(0), "must be between 1 and 3 (inclusive)")
		assert.EqualError(t, validate(4), "must be between 1 and 3 (inclusive)")
	})

	t.Run("orderedStrings", func(t *testing.T) {
		t.Parallel()

		validate := tinyflags.Between("b", "d")
		assert.NoError(t, validate("c"))
		assert.EqualError(t, validate("a"), "must be between b and d (inclusive)")
	})

	t.Run("invalidRangePanics", func(t *testing.T) {
		t.Parallel()

		assert.PanicsWithValue(
			t,
			"tinyflags: Between minimum must be less than or equal to maximum",
			func() { tinyflags.Between(3, 1) },
		)
	})

	t.Run("nanBoundsPanic", func(t *testing.T) {
		t.Parallel()

		assert.PanicsWithValue(
			t,
			"tinyflags: Between minimum must be less than or equal to maximum",
			func() { tinyflags.Between(math.NaN(), 1.0) },
		)
	})
}

// TestNotBlank verifies validation of empty and whitespace-only strings.
func TestNotBlank(t *testing.T) {
	t.Parallel()

	validate := tinyflags.NotBlank()
	assert.NoError(t, validate("tinyflags"))
	assert.NoError(t, validate(" tinyflags "))
	assert.EqualError(t, validate(""), "must not be blank")
	assert.EqualError(t, validate(" \t\n"), "must not be blank")
}

// TestBuiltInValidatorIntegration verifies validator errors retain normal flag context.
func TestBuiltInValidatorIntegration(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.Duration("timeout", time.Second, "Request timeout").
		Validate(tinyflags.Positive[time.Duration]())
	fs.Int("port", 8080, "Listen port").
		Validate(tinyflags.Between(1, 65535))
	fs.String("name", "app", "Application name").
		Validate(tinyflags.NotBlank())

	err := fs.Parse([]string{"--timeout=-1s", "--port=70000", "--name=   "})
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid value for flag --timeout: must be greater than zero")
	assert.ErrorContains(t, err, "invalid value for flag --port: must be between 1 and 65535 (inclusive)")
	assert.ErrorContains(t, err, "invalid value for flag --name: must not be blank")
}

// TestMinLength verifies inclusive minimum string length validation.
func TestMinLength(t *testing.T) {
	t.Parallel()

	validate := tinyflags.MinLength(3)
	assert.NoError(t, validate("abc"))
	assert.NoError(t, validate("abcd"))
	assert.EqualError(t, validate("ab"), "must be at least 3 characters")
}

// TestMaxLength verifies inclusive maximum string length validation.
func TestMaxLength(t *testing.T) {
	t.Parallel()

	validate := tinyflags.MaxLength(3)
	assert.NoError(t, validate("abc"))
	assert.NoError(t, validate("ab"))
	assert.EqualError(t, validate("abcd"), "must be at most 3 characters")
}

// TestLengthBetween verifies inclusive string length range validation.
func TestLengthBetween(t *testing.T) {
	t.Parallel()

	validate := tinyflags.LengthBetween(2, 4)
	assert.NoError(t, validate("ab"))
	assert.NoError(t, validate("abc"))
	assert.NoError(t, validate("abcd"))
	assert.EqualError(t, validate("a"), "must be between 2 and 4 characters (inclusive)")
	assert.EqualError(t, validate("abcde"), "must be between 2 and 4 characters (inclusive)")
}

// TestOptional verifies that zero values skip the wrapped validator.
func TestOptional(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		validate := tinyflags.Optional(tinyflags.MinLength(3))
		assert.NoError(t, validate(""))
		assert.NoError(t, validate("abc"))
		assert.EqualError(t, validate("ab"), "must be at least 3 characters")
	})

	t.Run("integer", func(t *testing.T) {
		t.Parallel()

		validate := tinyflags.Optional(tinyflags.AtLeast(3))
		assert.NoError(t, validate(0))
		assert.NoError(t, validate(3))
		assert.EqualError(t, validate(2), "must be at least 3")
	})

	t.Run("duration", func(t *testing.T) {
		t.Parallel()

		validate := tinyflags.Optional(tinyflags.Positive[time.Duration]())
		assert.NoError(t, validate(0))
		assert.NoError(t, validate(time.Second))
		assert.EqualError(t, validate(-time.Second), "must be greater than zero")
	})
}

// TestOptionalMinLengthIntegration verifies optional string validation retains normal flag context.
func TestOptionalMinLengthIntegration(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("oidc-session-secret", "", "OIDC session secret").
		Validate(tinyflags.Optional(tinyflags.MinLength(32)))

	assert.NoError(t, fs.Parse([]string{"--oidc-session-secret="}))

	err := fs.Parse([]string{"--oidc-session-secret=too-short"})
	require.Error(t, err)
	assert.EqualError(t, err, "invalid value for flag --oidc-session-secret: must be at least 32 characters")
}
