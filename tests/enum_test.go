package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Mode string

const (
	dev     Mode = "dev"
	staging Mode = "staging"
	prod    Mode = "prod"
)

type LogLevel int

const (
	debug LogLevel = iota
	info
	warn
)

func TestEnumFlag(t *testing.T) {
	t.Parallel()

	t.Run("static enum accepts allowed value", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		mode := tinyflags.Enum(fs, "mode", dev, "deployment mode", dev, staging, prod).Value()

		err := fs.Parse([]string{"--mode=prod"})
		require.NoError(t, err)
		assert.Equal(t, prod, *mode)

		switch *mode {
		case dev, staging, prod:
		default:
			t.Fatalf("unexpected mode %q", *mode)
		}
	})

	t.Run("static string enum rejects unknown value", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		mode := fs.Enum("mode", "dev", "deployment mode", "dev", "staging", "prod").Value()

		err := fs.Parse([]string{"--mode=test"})
		require.EqualError(t, err, "invalid value for flag --mode: must be one of: dev, staging, prod")
		assert.Equal(t, "dev", *mode)
	})

	t.Run("static typed enum var binds pointer", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		var mode Mode
		tinyflags.EnumVar(fs, &mode, "mode", dev, "deployment mode", dev, staging, prod)

		err := fs.Parse([]string{"--mode=staging"})
		require.NoError(t, err)
		assert.Equal(t, staging, mode)
	})

	t.Run("static iota enum accepts allowed value", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		level := tinyflags.Enum(fs, "level", info, "log level", debug, info, warn).Value()

		err := fs.Parse([]string{"--level=2"})
		require.NoError(t, err)
		assert.Equal(t, warn, *level)

		switch *level {
		case debug, info, warn:
		default:
			t.Fatalf("unexpected log level %d", *level)
		}
	})

	t.Run("static iota enum rejects unknown value", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		level := tinyflags.Enum(fs, "level", info, "log level", debug, info, warn).Value()

		err := fs.Parse([]string{"--level=99"})
		require.EqualError(t, err, "invalid value for flag --level: must be one of: 0, 1, 2")
		assert.Equal(t, info, *level)
	})

	t.Run("static named iota enum uses names", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		level := tinyflags.EnumMap(
			fs,
			"level",
			info,
			"log level",
			tinyflags.Choice("debug", debug),
			tinyflags.Choice("info", info),
			tinyflags.Choice("warn", warn),
		).Value()

		err := fs.Parse([]string{"--level=warn"})
		require.NoError(t, err)
		assert.Equal(t, warn, *level)
	})

	t.Run("static named iota enum shows names in help", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		tinyflags.EnumMap(
			fs,
			"level",
			info,
			"log level",
			tinyflags.Choice("debug", debug),
			tinyflags.Choice("info", info),
			tinyflags.Choice("warn", warn),
		)

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "--level <debug|info|warn>")
		assert.Contains(t, err.Error(), "log level (allowed: debug, info, warn) (default: info)")
	})

	t.Run("static named iota enum rejects unknown name", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		level := tinyflags.EnumMap(
			fs,
			"level",
			info,
			"log level",
			tinyflags.Choice("debug", debug),
			tinyflags.Choice("info", info),
			tinyflags.Choice("warn", warn),
		).Value()

		err := fs.Parse([]string{"--level=trace"})
		require.EqualError(t, err, "invalid value for flag --level: must be one of: debug, info, warn")
		assert.Equal(t, info, *level)
	})

	t.Run("dynamic typed enum stores values by id", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		svc := fs.DynamicGroup("svc")
		mode := tinyflags.DynamicEnum(svc, "mode", dev, "deployment mode", dev, staging, prod)

		err := fs.Parse([]string{
			"--svc.api.mode=prod",
			"--svc.worker.mode=staging",
		})
		require.NoError(t, err)
		assert.Equal(t, map[string]Mode{"api": prod, "worker": staging}, mode.Values())
	})

	t.Run("dynamic named iota enum shows names in help", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		svc := fs.DynamicGroup("svc")
		tinyflags.DynamicEnumMap(
			svc,
			"level",
			info,
			"log level",
			tinyflags.Choice("debug", debug),
			tinyflags.Choice("info", info),
			tinyflags.Choice("warn", warn),
		)

		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "--svc.<ID>.level <debug|info|warn>")
		assert.Contains(t, err.Error(), "log level (allowed: debug, info, warn) (default: info)")
	})

	t.Run("dynamic enum rejects unknown value", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		svc := fs.DynamicGroup("svc")
		svc.Enum("mode", "dev", "deployment mode", "dev", "staging", "prod")

		err := fs.Parse([]string{"--svc.api.mode=test"})
		require.EqualError(t, err, "invalid value for flag --svc.api.mode: must be one of: dev, staging, prod")
	})
}
