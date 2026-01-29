package engine

import (
	"testing"

	"github.com/containeroo/tinyflags/internal/core"
)

type testValue struct {
	set string
}

func (v *testValue) Set(s string) error { v.set = s; return nil }
func (v *testValue) Get() any           { return v.set }
func (v *testValue) Changed() bool      { return v.set != "" }
func (v *testValue) Default() string    { return "" }

func TestParseArgsWithFSMHandlesUnknown(t *testing.T) {
	t.Parallel()

	t.Run("unknownHandled", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		fs.OnUnknownFlag(func(name string) error { return nil })
		fs.RegisterFlag("known", &core.BaseFlag{
			Name:  "known",
			Value: &testValue{},
		})

		_, err := parseArgsWithFSM(fs, []string{"--unknown", "--known=ok"})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("doubleDashStopsParsing", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		val := &testValue{}
		fs.RegisterFlag("known", &core.BaseFlag{
			Name:  "known",
			Value: val,
		})

		out, err := parseArgsWithFSM(fs, []string{"--", "--known=skip", "pos"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out) != 2 || out[0] != "--known=skip" || out[1] != "pos" {
			t.Fatalf("unexpected positional output: %v", out)
		}
		if val.set != "" {
			t.Fatalf("flag should not be set when after --")
		}
	})

	t.Run("unknownDynamicGroupShowsFlag", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		_, err := parseArgsWithFSM(fs, []string{"--target.unknown.http=service"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if got := err.Error(); got != "unknown dynamic group \"target\" in flag --target.unknown.http=service" {
			t.Fatalf("unexpected error: %q", got)
		}
	})

	t.Run("unknownDynamicFieldShowsFlag", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		g := fs.DynamicGroup("http")
		g.String("addr", "", "addr")

		_, err := parseArgsWithFSM(fs, []string{"--http.alpha.port=8080"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if got := err.Error(); got != "unknown dynamic field \"port\" in flag --http.alpha.port=8080" {
			t.Fatalf("unexpected error: %q", got)
		}
	})
}
