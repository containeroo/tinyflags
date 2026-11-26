package tinyflags_test

import (
	"context"
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
)

// TestBooleanExplicitValues verifies that only bare boolean flags imply true.
func TestBooleanExplicitValues(t *testing.T) {
	t.Parallel()

	t.Run("explicitFalseOverridesEnvironment", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_ENABLED" {
				return "true"
			}
			return ""
		})
		enabled := fs.Bool("enabled", false, "Enable feature")

		if err := fs.Parse([]string{"--enabled=false"}); err != nil {
			t.Fatal(err)
		}
		if *enabled.Value() || !enabled.Changed() {
			t.Fatal("expected an explicitly set false value")
		}
		if got := fs.Origin("enabled"); got.Source != tinyflags.ValueSourceFlag || got.Key != "--enabled" {
			t.Fatalf("unexpected origin: %v", got)
		}
	})

	t.Run("explicitTrue", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		enabled := fs.Bool("enabled", false, "Enable feature").DisableEnv().Value()

		if err := fs.Parse([]string{"--enabled=true"}); err != nil {
			t.Fatal(err)
		}
		if !*enabled {
			t.Fatal("expected true")
		}
	})

	t.Run("bareLongFlag", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		enabled := fs.Bool("enabled", false, "Enable feature").DisableEnv().Value()

		if err := fs.Parse([]string{"--enabled"}); err != nil {
			t.Fatal(err)
		}
		if !*enabled {
			t.Fatal("expected true")
		}
	})

	t.Run("bareShortFlag", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		enabled := fs.Bool("enabled", false, "Enable feature").Short("e").DisableEnv().Value()

		if err := fs.Parse([]string{"-e"}); err != nil {
			t.Fatal(err)
		}
		if !*enabled {
			t.Fatal("expected true")
		}
	})

	t.Run("invalidExplicitValue", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Bool("enabled", false, "Enable feature").DisableEnv()

		err := fs.Parse([]string{"--enabled=invalid"})
		if err == nil || !strings.Contains(err.Error(), "invalid value for flag --enabled") {
			t.Fatalf("expected invalid boolean error, got %v", err)
		}
	})

	t.Run("emptyExplicitValue", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Bool("enabled", false, "Enable feature").DisableEnv()

		err := fs.Parse([]string{"--enabled="})
		if err == nil || !strings.Contains(err.Error(), "invalid value for flag --enabled") {
			t.Fatalf("expected invalid boolean error, got %v", err)
		}
	})

	t.Run("strictRequiresValue", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.Bool("enabled", false, "Enable feature").Strict().DisableEnv()

		err := fs.Parse([]string{"--enabled"})
		if err == nil || !strings.Contains(err.Error(), "missing value for flag --enabled") {
			t.Fatalf("expected missing value error, got %v", err)
		}
	})

	t.Run("dynamicExplicitFalse", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		group := fs.DynamicGroup("service")
		enabled := group.Bool("enabled", true, "Enable service")

		if err := fs.Parse([]string{"--service.api.enabled=false"}); err != nil {
			t.Fatal(err)
		}
		if value, ok := enabled.Get("api"); value || !ok {
			t.Fatal("expected an explicitly set false value")
		}
		if got := fs.Origin("service.api.enabled"); got.Source != tinyflags.ValueSourceFlag || got.Key != "--service.api.enabled" {
			t.Fatalf("unexpected origin: %v", got)
		}
	})

	t.Run("dynamicBareFlag", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		enabled := fs.DynamicGroup("service").Bool("enabled", false, "Enable service")

		if err := fs.Parse([]string{"--service.api.enabled"}); err != nil {
			t.Fatal(err)
		}
		if value, ok := enabled.Get("api"); !value || !ok {
			t.Fatal("expected an explicitly set true value")
		}
	})

	t.Run("dynamicInvalidExplicitValue", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.DynamicGroup("service").Bool("enabled", false, "Enable service")

		err := fs.Parse([]string{"--service.api.enabled=invalid"})
		if err == nil || !strings.Contains(err.Error(), "invalid value for flag --service.api.enabled") {
			t.Fatalf("expected invalid boolean error, got %v", err)
		}
	})

	t.Run("commandRunnerOverridesEnvironment", func(t *testing.T) {
		t.Parallel()

		root := tinyflags.NewCommand("warden", tinyflags.ContinueOnError)
		root.EnvPrefix("WARDEN_")
		serve := root.Command("serve", "Start server")
		serve.SetGetEnvFn(func(key string) string {
			if key == "WARDEN__OBSERVE_ONLY" {
				return "true"
			}
			return ""
		})
		observeOnly := serve.Bool("observe-only", false, "Observe only").Value()

		var executed bool
		var observed bool
		serve.Run(func(value bool) {
			executed = true
			observed = value
		}, observeOnly)

		runner, err := root.ParseRunner([]string{"serve"})
		if err != nil {
			t.Fatal(err)
		}
		if err := runner.Run(context.Background()); err != nil {
			t.Fatal(err)
		}
		if !executed || !observed {
			t.Fatal("expected environment to enable observe-only")
		}

		executed = false
		runner, err = root.ParseRunner([]string{"serve", "--observe-only=false"})
		if err != nil {
			t.Fatal(err)
		}
		if err := runner.Run(context.Background()); err != nil {
			t.Fatal(err)
		}
		if !executed || observed {
			t.Fatal("expected CLI to disable observe-only")
		}
	})
}
