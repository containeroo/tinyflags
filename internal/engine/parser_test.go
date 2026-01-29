package engine

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBeforeParseHookErrors(t *testing.T) {
	t.Parallel()

	fs := NewFlagSet("app", ContinueOnError)
	fs.BeforeParse(func(args []string) ([]string, error) {
		return nil, errors.New("hook failure")
	})

	err := fs.Parse([]string{"--help"})
	require.Error(t, err)
	require.EqualError(t, err, "hook failure")
}
