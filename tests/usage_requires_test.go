package tinyflags_test

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHideRequires_RemovesSuffixFromHelp(t *testing.T) {
	t.Parallel()

	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.String("db", "", "Database connection")

	fs.String("dsn", "", "Connection string").
		Requires("db").
		HideRequires().
		Value()

	err := fs.Parse([]string{"--help"})
	require.Error(t, err)

	assert.NotContains(t, err.Error(), "(Requires: db)")
}
