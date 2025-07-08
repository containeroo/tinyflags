package tinyflags

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Parallel()

	t.Run("Error Message", func(t *testing.T) {
		t.Parallel()
		err := &HelpRequested{Message: "This is a help message"}
		assert.Equal(t, "This is a help message", err.Error(), "Error() should return the correct message")
	})

	t.Run("Version Message", func(t *testing.T) {
		t.Parallel()
		err := &VersionRequested{Version: "v1.2.3"}
		assert.EqualError(t, err, "v1.2.3", "Error() should return the correct version message")
	})

	t.Run("IsHelpRequested returns true", func(t *testing.T) {
		t.Parallel()
		err := RequestHelp("help me")
		assert.True(t, IsHelpRequested(err))
	})

	t.Run("IsHelpRequested returns false", func(t *testing.T) {
		t.Parallel()
		err := errors.New("regular error")
		assert.False(t, IsHelpRequested(err))
	})

	t.Run("IsVersionRequested returns true", func(t *testing.T) {
		t.Parallel()
		err := RequestVersion("1.0.0")
		assert.True(t, IsVersionRequested(err))
	})

	t.Run("IsVersionRequested returns false", func(t *testing.T) {
		t.Parallel()
		err := errors.New("not a version")
		assert.False(t, IsVersionRequested(err))
	})

	t.Run("RequestHelp creates HelpRequested", func(t *testing.T) {
		t.Parallel()
		err := RequestHelp("show help")
		var helpErr *HelpRequested
		assert.True(t, errors.As(err, &helpErr))
		assert.Equal(t, "show help", helpErr.Message)
	})

	t.Run("RequestVersion creates VersionRequested", func(t *testing.T) {
		t.Parallel()
		err := RequestVersion("v9.9.9")
		var verErr *VersionRequested
		assert.True(t, errors.As(err, &verErr))
		assert.Equal(t, "v9.9.9", verErr.Version)
	})
}
