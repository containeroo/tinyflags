package engine

import (
	"testing"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testValue struct {
	set string
}

func (v *testValue) Set(s string) error { v.set = s; return nil }
func (v *testValue) Get() any           { return v.set }
func (v *testValue) Changed() bool      { return v.set != "" }
func (v *testValue) Default() string    { return "" }

func TestRunArgParserFSMHandlesUnknown(t *testing.T) {
	t.Parallel()

	t.Run("unknownHandled", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		fs.OnUnknownFlag(func(name string) error { return nil })
		fs.RegisterFlag("known", &core.BaseFlag{
			Name:  "known",
			Value: &testValue{},
		})

		_, err := runArgParserFSM(fs, []string{"--unknown", "--known=ok"})
		require.NoError(t, err)
	})

	t.Run("doubleDashStopsParsing", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		val := &testValue{}
		fs.RegisterFlag("known", &core.BaseFlag{
			Name:  "known",
			Value: val,
		})

		out, err := runArgParserFSM(fs, []string{"--", "--known=skip", "pos"})
		require.NoError(t, err)
		assert.Equal(t, []string{"--known=skip", "pos"}, out)
		assert.Equal(t, "", val.set)
	})

	t.Run("unknownDynamicGroupShowsFlag", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		_, err := runArgParserFSM(fs, []string{"--target.unknown.http=service"})
		require.Error(t, err)
		assert.EqualError(t, err, "unknown dynamic group \"target\" in flag --target.unknown.http=service")
	})

	t.Run("unknownDynamicFieldShowsFlag", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		g := fs.DynamicGroup("http")
		g.String("addr", "", "addr")

		_, err := runArgParserFSM(fs, []string{"--http.alpha.port=8080"})
	require.Error(t, err)
	assert.EqualError(t, err, "unknown dynamic field \"port\" in flag --http.alpha.port=8080")
	})
}
