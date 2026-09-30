package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBytesHumanReadableInput(t *testing.T) {
	t.Parallel()

	t.Run("cli", func(t *testing.T) {
		t.Parallel()

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		size := fs.Bytes("size", 0, "size").Value()

		err := fs.Parse([]string{"--size=32MiB"})
		require.NoError(t, err)
		assert.Equal(t, uint64(32<<20), *size)
	})

	t.Run("environment", func(t *testing.T) {
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		size := fs.Bytes("size", 0, "size").Value()
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_SIZE" {
				return "32MiB"
			}
			return ""
		})

		err := fs.Parse(nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(32<<20), *size)
	})

	t.Run("cliWinsOverEnvironment", func(t *testing.T) {
		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.EnvPrefix("APP")
		size := fs.Bytes("size", 0, "size").Value()
		fs.SetGetEnvFn(func(key string) string {
			if key == "APP_SIZE" {
				return "64MiB"
			}
			return ""
		})

		err := fs.Parse([]string{"--size=32MiB"})
		require.NoError(t, err)
		assert.Equal(t, uint64(32<<20), *size)
	})
}

func TestBytesSliceHumanReadableInput(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	sizes := fs.BytesSlice("size", nil, "sizes").Value()

	err := fs.Parse([]string{"--size=32MiB,64MB,1024"})
	require.NoError(t, err)
	assert.Equal(t, []uint64{32 << 20, 64_000_000, 1024}, *sizes)
}
