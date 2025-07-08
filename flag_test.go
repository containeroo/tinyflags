package tinyflags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagEnv(t *testing.T) {
	t.Parallel()

	t.Run("Env and DisableEnv", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("name", "", "user name").Env("USERNAME")
		assert.Equal(t, "USERNAME", f.bf.envKey)

		assert.Panics(t, func() {
			fs.String("x", "", "").DisableEnv().Env("FAIL")
		})
		assert.Panics(t, func() {
			fs.String("y", "", "").Env("OK").DisableEnv()
		})
	})

	t.Run("Group assigns flag to named group", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f1 := fs.String("foo", "", "").Group("g1")
		f2 := fs.String("bar", "", "").Group("g1")
		f3 := fs.String("baz", "", "") // not in group

		assert.Len(t, fs.groups, 1)
		assert.Equal(t, "g1", fs.groups[0].name)
		assert.Contains(t, fs.groups[0].flags, f1.bf)
		assert.Contains(t, fs.groups[0].flags, f2.bf)
		assert.NotEqual(t, f3.bf.group, fs.groups[0])
	})
}

func TestFlagGroup(t *testing.T) {
	t.Parallel()

	t.Run("Group avoids duplicates", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("dup", "", "").Group("x").Group("x")

		assert.Len(t, fs.groups, 1)
		assert.Len(t, fs.groups[0].flags, 1)
		assert.Equal(t, f.bf.group, fs.groups[0])
	})

	t.Run("Group empty name returns same flag", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("noop", "", "").Group("")
		assert.Nil(t, f.bf.group)
	})
}

func TestFlagDeprecated(t *testing.T) {
	t.Parallel()

	t.Run("Deprecated sets message", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("old", "", "").Deprecated("use --new instead")
		assert.Equal(t, "use --new instead", f.bf.deprecated)
	})
}

func TestFlagRequired(t *testing.T) {
	t.Parallel()

	t.Run("Required sets required true", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("name", "", "").Required()
		assert.True(t, f.bf.required)
	})
}

func TestFlagMetavar(t *testing.T) {
	t.Parallel()

	t.Run("Metavar overrides default placeholder", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("file", "", "").Metavar("PATH")
		assert.Equal(t, "PATH", f.bf.metavar)
	})
}

func TestFlagChoices(t *testing.T) {
	t.Parallel()

	t.Run("Choices restricts allowed values", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("mode", "", "").Choices("fast", "slow")

		assert.ElementsMatch(t, []string{"fast", "slow"}, f.bf.allowed)
		assert.NotNil(t, f.bf.value.(*FlagItem[string]).validator)
	})

	t.Run("Choices on non-FlagItem does not panic", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		var fake Flag[any]
		fake.fs = fs
		fake.bf = &baseFlag{name: "x", value: &mockValue{}}
		assert.NotPanics(t, func() {
			fake.Choices("abc")
		})
	})

	t.Run("Choices validator enforces allowed values", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("color", "", "").Choices("red", "green", "blue")

		v := f.bf.value.(*FlagItem[string])

		assert.NoError(t, v.validator("red"), "allowed value should pass")
		assert.NoError(t, v.validator("green"), "allowed value should pass")
		assert.Error(t, v.validator("yellow"), "disallowed value should fail")
	})
}

func TestFlagValidator(t *testing.T) {
	t.Parallel()

	t.Run("Validator installs custom validation", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.Int("port", 0, "").Validator(func(v int) error {
			if v < 0 || v > 65535 {
				return assert.AnError
			}
			return nil
		})
		assert.NotNil(t, f.bf.value.(*FlagItem[int]).validator)
	})

	t.Run("Validator on non-FlagItem does not panic", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		var fake Flag[any]
		fake.fs = fs
		fake.bf = &baseFlag{name: "x", value: &mockValue{}}
		assert.NotPanics(t, func() {
			fake.Validator(func(any) error { return nil })
		})
	})
}

func TestFlagHidden(t *testing.T) {
	t.Parallel()

	t.Run("Hidden hides flag from help", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("secret", "", "").Hidden()
		assert.True(t, f.bf.hidden)
	})
}

func TestFlagValue(t *testing.T) {
	t.Parallel()

	t.Run("Value returns parsed pointer", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.String("x", "abc", "")
		ptr := f.Value()
		assert.Equal(t, "abc", *ptr)
	})
}

func TestFlagDelimiter(t *testing.T) {
	t.Parallel()

	t.Run("Delimiter sets value when HasDelimiter is implemented", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ExitOnError)
		f := fs.StringSlice("tags", []string{"a", "b"}, "").Delimiter("|")
		v := f.bf.value.(HasDelimiter)
		assert.Implements(t, (*HasDelimiter)(nil), v)
	})

	t.Run("Delimiter no-op if value doesn't implement HasDelimiter", func(t *testing.T) {
		t.Parallel()
		var fake SliceFlag[int]
		fake.bf = &baseFlag{name: "x", value: &mockValue{}}
		assert.NotPanics(t, func() {
			fake.Delimiter(";")
		})
	})
}
