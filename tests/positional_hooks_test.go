package tinyflags_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPositionalValidateAll(t *testing.T) {
	t.Parallel()

	t.Run("runs_after_individual_hooks_and_before_finalizer", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		var calls []string
		fs.SetPositionalValidate(func(value string) error {
			calls = append(calls, "validate:"+value)
			return nil
		})
		fs.SetPositionalFinalize(func(value string) string {
			calls = append(calls, "finalize:"+value)
			return strings.ToUpper(value)
		})
		fs.SetPositionalValidateAll(func(args []string) error {
			calls = append(calls, "validate all")
			assert.Equal(t, []string{"A", "B"}, args)
			return nil
		})
		fs.SetPositionalFinalizeAll(func(args []string) ([]string, error) {
			calls = append(calls, "finalize all")
			return append(args, "C"), nil
		})
		fs.Validate(func() error {
			calls = append(calls, "cross-flag")
			assert.Equal(t, []string{"A", "B", "C"}, fs.Args())
			return nil
		})

		require.NoError(t, fs.Parse([]string{"a", "b"}))
		assert.Equal(t, []string{
			"validate:a", "validate:b", "finalize:a", "finalize:b",
			"validate all", "finalize all", "cross-flag",
		}, calls)
	})

	t.Run("rejects_invalid_collection", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		expected := errors.New("expected pairs")
		fs.SetPositionalValidateAll(func(args []string) error {
			if len(args)%2 != 0 {
				return expected
			}
			return nil
		})
		fs.SetPositionalFinalizeAll(func(args []string) ([]string, error) {
			t.Fatal("finalizer must not run after validation fails")
			return args, nil
		})

		require.ErrorIs(t, fs.Parse([]string{"one"}), expected)
	})

	t.Run("cannot_mutate_stored_arguments", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.SetPositionalValidateAll(func(args []string) error {
			args[0] = "changed"
			return nil
		})

		require.NoError(t, fs.Parse([]string{"original"}))
		assert.Equal(t, []string{"original"}, fs.Args())
	})
}

func TestPositionalFinalizeAll(t *testing.T) {
	t.Parallel()

	t.Run("deduplicates_after_input_count_validation", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.ExactPositional(3)
		fs.SetPositionalFinalizeAll(func(args []string) ([]string, error) {
			return slices.Compact(args), nil
		})

		require.NoError(t, fs.Parse([]string{"one", "one", "two"}))
		assert.Equal(t, []string{"one", "two"}, fs.Args())
		value, ok := fs.Arg(1)
		assert.True(t, ok)
		assert.Equal(t, "two", value)
		_, ok = fs.Arg(2)
		assert.False(t, ok)
	})

	t.Run("expands_beyond_input_count", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.MaxPositional(1)
		fs.SetPositionalFinalizeAll(func(args []string) ([]string, error) {
			return strings.Split(args[0], ","), nil
		})

		require.NoError(t, fs.Parse([]string{"a,b"}))
		assert.Equal(t, []string{"a", "b"}, fs.Args())
	})

	t.Run("allows_nil_result", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.MinPositional(1)
		fs.SetPositionalFinalizeAll(func([]string) ([]string, error) {
			return nil, nil
		})

		require.NoError(t, fs.Parse([]string{"one"}))
		assert.Nil(t, fs.Args())
	})

	t.Run("does_not_store_failed_result", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		expected := errors.New("invalid input")
		fs.SetPositionalFinalizeAll(func(args []string) ([]string, error) {
			args[0] = "changed"
			return args, expected
		})
		fs.Validate(func() error {
			t.Fatal("cross-flag validator must not run after finalization fails")
			return nil
		})

		require.ErrorIs(t, fs.Parse([]string{"original"}), expected)
		assert.Equal(t, []string{"original"}, fs.Args())
	})

	t.Run("honors_panic_mode", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.PanicOnError)
		expected := errors.New("invalid input")
		fs.SetPositionalFinalizeAll(func([]string) ([]string, error) {
			return nil, expected
		})

		assert.PanicsWithValue(t, expected, func() {
			_ = fs.Parse([]string{"one"})
		})
	})
}

func TestPositionalCollectionHookLifecycle(t *testing.T) {
	t.Parallel()

	t.Run("runs_on_empty_input", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.NoPositional()
		validated := false
		fs.SetPositionalValidateAll(func(args []string) error {
			validated = true
			assert.Empty(t, args)
			return nil
		})
		fs.SetPositionalFinalizeAll(func(args []string) ([]string, error) {
			assert.True(t, validated)
			assert.Empty(t, args)
			return []string{"default"}, nil
		})

		require.NoError(t, fs.Parse(nil))
		assert.Equal(t, []string{"default"}, fs.Args())
	})

	t.Run("skips_on_invalid_count", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.MinPositional(1)
		fs.SetPositionalValidateAll(func([]string) error {
			t.Fatal("validator must not run after count failure")
			return nil
		})
		fs.SetPositionalFinalizeAll(func([]string) ([]string, error) {
			t.Fatal("finalizer must not run after count failure")
			return nil, nil
		})

		require.EqualError(t, fs.Parse(nil), "expected at least 1 positional argument, got 0")
	})

	t.Run("skips_on_individual_validation_error", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		expected := errors.New("invalid argument")
		fs.SetPositionalValidate(func(string) error { return expected })
		fs.SetPositionalValidateAll(func([]string) error {
			t.Fatal("collection validator must not run")
			return nil
		})
		fs.SetPositionalFinalizeAll(func([]string) ([]string, error) {
			t.Fatal("collection finalizer must not run")
			return nil, nil
		})

		require.ErrorIs(t, fs.Parse([]string{"one"}), expected)
	})

	t.Run("skips_help_and_version", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Version("1.0.0")
		fs.SetPositionalValidateAll(func([]string) error {
			t.Fatal("validator must not run for built-in requests")
			return nil
		})
		fs.SetPositionalFinalizeAll(func([]string) ([]string, error) {
			t.Fatal("finalizer must not run for built-in requests")
			return nil, nil
		})

		assert.True(t, tinyflags.IsHelpRequested(fs.Parse([]string{"--help"})))
		assert.True(t, tinyflags.IsVersionRequested(fs.Parse([]string{"--version"})))
	})

	t.Run("reparse_and_clear_hooks", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		calls := 0
		fs.SetPositionalValidateAll(func([]string) error {
			calls++
			return nil
		})
		fs.SetPositionalFinalizeAll(func(args []string) ([]string, error) {
			return append(args, "extra"), nil
		})

		require.NoError(t, fs.Parse([]string{"first"}))
		assert.Equal(t, []string{"first", "extra"}, fs.Args())
		require.NoError(t, fs.Parse([]string{"second"}))
		assert.Equal(t, []string{"second", "extra"}, fs.Args())
		assert.Equal(t, 2, calls)

		fs.SetPositionalValidateAll(nil)
		fs.SetPositionalFinalizeAll(nil)
		require.NoError(t, fs.Parse([]string{"third"}))
		assert.Equal(t, []string{"third"}, fs.Args())
		assert.Equal(t, 2, calls)
	})
}

func TestCommandPositionalCollectionHooks(t *testing.T) {
	t.Parallel()

	t.Run("keeps_command_scopes_separate", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
		child := root.Command("child", "Child command").MinPositional(1)
		root.SetPositionalValidateAll(func(args []string) error {
			assert.Equal(t, []string{"root-arg"}, args)
			return nil
		})
		child.SetPositionalFinalizeAll(func(args []string) ([]string, error) {
			return slices.Compact(args), nil
		})

		require.NoError(t, root.Parse([]string{"root-arg", "child", "a", "a"}))
		assert.Equal(t, []string{"root-arg"}, root.Args())
		assert.Equal(t, []string{"a"}, child.Args())
	})

	t.Run("returns_finalizer_error", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
		child := root.Command("child", "Child command")
		expected := errors.New("invalid address")
		child.SetPositionalFinalizeAll(func([]string) ([]string, error) {
			return nil, expected
		})

		require.ErrorIs(t, root.Parse([]string{"child", "bad"}), expected)
	})
}
