package tinyflags

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseArgsWithFSM(t *testing.T) {
	t.Parallel()

	t.Run("positional arguments", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		pos, err := parseArgsWithFSM(fs, []string{"arg1", "arg2"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"arg1", "arg2"}, pos)
	})

	t.Run("stop parsing on --", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		pos, err := parseArgsWithFSM(fs, []string{"one", "--", "--flag", "rest"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"one", "--flag", "rest"}, pos)
	})

	t.Run("long flag with value", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		tv := &mockValue{}
		fs.flags["name"] = &baseFlag{name: "name", value: tv}
		pos, err := parseArgsWithFSM(fs, []string{"--name=joe"})
		assert.NoError(t, err)
		assert.Equal(t, "joe", tv.val)
		assert.Empty(t, pos)
	})

	t.Run("long flag missing value", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		tv := &mockValue{}
		fs.flags["opt"] = &baseFlag{name: "opt", value: tv}
		_, err := parseArgsWithFSM(fs, []string{"--opt"})
		assert.EqualError(t, err, "missing value for flag: --opt")
	})

	t.Run("unknown long flag", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		_, err := parseArgsWithFSM(fs, []string{"--nope"})
		assert.EqualError(t, err, "unknown flag: --nope")
	})

	t.Run("long flag with invalid value", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		tv := &mockValue{errToReturn: errors.New("bad")}
		fs.flags["fail"] = &baseFlag{name: "fail", value: tv}
		_, err := parseArgsWithFSM(fs, []string{"--fail=no"})
		assert.EqualError(t, err, "invalid value for flag --fail: bad")
	})

	t.Run("short non-strict bool sets true", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		debug := fs.BoolP("debug", "d", false, "").Value()
		args, err := parseArgsWithFSM(fs, []string{"-d"})
		assert.NoError(t, err)
		assert.Len(t, args, 0)
		assert.True(t, *debug)
	})

	t.Run("short flag with value in next arg", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		tv := &mockValue{}
		fs.flags["p"] = &baseFlag{name: "p", short: "p", value: tv}
		_, err := parseArgsWithFSM(fs, []string{"-p", "8080"})
		assert.NoError(t, err)
		assert.Equal(t, "8080", tv.val)
	})

	t.Run("short flag with missing value", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		tv := &mockValue{}
		fs.flags["x"] = &baseFlag{name: "x", short: "x", value: tv}
		_, err := parseArgsWithFSM(fs, []string{"-x"})
		assert.EqualError(t, err, "missing value for flag: -x")
	})

	t.Run("short grouped flags", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		a := fs.BoolP("a", "a", false, "").Value()
		b := fs.BoolP("b", "b", false, "").Value()
		c := fs.BoolP("c", "c", false, "").Value()
		_, err := parseArgsWithFSM(fs, []string{"-abc"})
		assert.NoError(t, err)
		assert.True(t, *a)
		assert.True(t, *b)
		assert.True(t, *c)
	})

	t.Run("short flag with combined value", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		tv := &mockValue{}
		fs.flags["p"] = &baseFlag{name: "p", short: "p", value: tv}
		_, err := parseArgsWithFSM(fs, []string{"-p8080"})
		assert.NoError(t, err)
		assert.Equal(t, "8080", tv.val)
	})

	t.Run("short unknown flag", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		_, err := parseArgsWithFSM(fs, []string{"-z"})
		assert.EqualError(t, err, "unknown short flag: -z")
	})

	t.Run("splitFlagArg splits correctly", func(t *testing.T) {
		t.Parallel()
		k, v, ok := splitFlagArg("foo=bar")
		assert.Equal(t, "foo", k)
		assert.Equal(t, "bar", v)
		assert.True(t, ok)

		k, v, ok = splitFlagArg("baz")
		assert.Equal(t, "baz", k)
		assert.Equal(t, "", v)
		assert.False(t, ok)
	})

	t.Run("trySet wraps error", func(t *testing.T) {
		t.Parallel()
		tv := &mockValue{errToReturn: errors.New("bad")}
		err := trySet(tv, "x", "%s failed: %w", "xyz")
		assert.EqualError(t, err, "xyz failed: bad")
	})

	t.Run("long flag with value in next arg", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		tv := &mockValue{}
		fs.flags["name"] = &baseFlag{name: "name", value: tv}

		pos, err := parseArgsWithFSM(fs, []string{"--name", "Alice"})
		assert.NoError(t, err)
		assert.Equal(t, "Alice", tv.val)
		assert.Empty(t, pos)
	})

	t.Run("long flag with invalid value in next arg", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("app", ContinueOnError)
		tv := &mockValue{errToReturn: errors.New("fail")}
		fs.flags["age"] = &baseFlag{name: "age", value: tv}

		_, err := parseArgsWithFSM(fs, []string{"--age", "oops"})
		assert.EqualError(t, err, "invalid value for flag --age: fail")
	})
}
