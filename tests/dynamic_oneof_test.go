package tinyflags_test

import (
	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDynamicOneOfPerInstance(t *testing.T) {
	for _, test := range []struct {
		name     string
		args     []string
		required bool
		fail     bool
	}{
		{name: "mixed sources across instances", required: true, args: []string{"--update.a.path=/A", "--update.a.value=literal", "--update.b.path=/B", "--update.b.valueFromEnv=TOKEN"}},
		{name: "both on same instance", required: true, args: []string{"--update.a.path=/A", "--update.a.value=literal", "--update.a.valueFromEnv=TOKEN"}, fail: true},
		{name: "neither on active instance", required: true, args: []string{"--update.a.path=/A"}, fail: true},
		{name: "missing source on another instance", required: true, args: []string{"--update.a.path=/A", "--update.a.value=literal", "--update.b.path=/B"}, fail: true},
		{name: "optional missing source", args: []string{"--update.a.path=/A"}},
		{name: "optional still exclusive", args: []string{"--update.a.value=x", "--update.a.valueFromEnv=TOKEN"}, fail: true},
		{name: "empty literal is present", required: true, args: []string{"--update.a.path=/A", "--update.a.value="}},
		{name: "defaults are absent", required: true, args: []string{"--update.a.path=/A"}, fail: true},
		{name: "no active instances", required: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
			group := fs.DynamicGroup("update")
			group.String("path", "", "XML selector")
			value := group.String("value", "default", "literal")
			fromEnv := group.String("valueFromEnv", "TOKEN", "source")
			constraint := group.OneOf(value, fromEnv)
			if test.required {
				constraint.Required()
			}
			err := fs.Parse(test.args)
			if test.fail {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "--update.")
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDynamicOneOfMixedFieldTypesAndReuse(t *testing.T) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	group := fs.DynamicGroup("input")
	group.String("path", "", "selector")
	number := group.Int("number", 5, "number")
	boolean := group.Bool("enabled", true, "enabled").Strict()
	list := group.StringSlice("list", nil, "list")
	group.OneOf(number, boolean, list).Required()
	require.NoError(t, fs.Parse([]string{"--input.a.number=0", "--input.b.enabled=false", "--input.c.list=one,two"}))
	require.Error(t, fs.Parse([]string{"--input.a.path=/A"}))
	require.NoError(t, fs.Parse([]string{"--input.a.list=last"}))
	require.NoError(t, fs.Parse(nil))
}

func TestDynamicOneOfEnvironmentAndCLI(t *testing.T) {
	t.Setenv("ONEOF_TEST_UPDATE_A_VALUE", "environment")
	t.Setenv("ONEOF_TEST_UPDATE_B_VALUEFROMENV", "TOKEN")
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("ONEOF_TEST")
	group := fs.DynamicGroup("update")
	value := group.String("value", "", "literal")
	fromEnv := group.String("valueFromEnv", "", "source")
	group.OneOf(value, fromEnv).Required()
	require.NoError(t, fs.Parse([]string{"--update.a.value=cli"}))
	literal, present := value.Get("a")
	require.True(t, present)
	assert.Equal(t, "cli", literal)
	require.Error(t, fs.Parse([]string{"--update.a.valueFromEnv=OTHER"}))
}

func TestDynamicOneOfRejectsInvalidHandles(t *testing.T) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	group := fs.DynamicGroup("update")
	value := group.String("value", "", "literal")
	foreign := fs.DynamicGroup("other").String("value", "", "other")
	var nilField tinyflags.DynamicField
	typedNil := nilFieldPointer(value)
	assert.Panics(t, func() { group.OneOf() })
	assert.Panics(t, func() { group.OneOf(nilField) })
	assert.Panics(t, func() { group.OneOf(typedNil) })
	assert.Panics(t, func() { group.OneOf(foreign) })
	assert.Panics(t, func() { group.OneOf(value, value) })
	// Invalid registration must not leave a partially registered constraint.
	require.NoError(t, fs.Parse([]string{"--update.value.value=valid"}))
}

func TestDynamicOneOfHelpAndIndependentGlobalGroups(t *testing.T) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("first", "", "first").OneOfGroup("global")
	fs.String("second", "", "second").OneOfGroup("global")
	group := fs.DynamicGroup("update")
	value := group.String("value", "", "literal")
	fromEnv := group.String("valueFromEnv", "", "source")
	group.OneOf(value, fromEnv).Required()
	require.NoError(t, fs.Parse([]string{"--first=one", "--update.a.value=x", "--update.b.valueFromEnv=TOKEN"}))
	require.Error(t, fs.Parse([]string{"--first=one", "--second=two"}))
	err := fs.Parse([]string{"--help"})
	require.True(t, tinyflags.IsHelpRequested(err))
	assert.Contains(t, err.Error(), "per instance")
	assert.Contains(t, err.Error(), "one of, required")
}

func nilFieldPointer[T any](field *T) *T { return nil }
