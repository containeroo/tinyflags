package engine

import (
	"errors"
	"testing"
)

func TestBeforeParseHookErrors(t *testing.T) {
	t.Parallel()

	fs := NewFlagSet("app", ContinueOnError)
	fs.BeforeParse(func(args []string) ([]string, error) {
		return nil, errors.New("hook failure")
	})

	err := fs.Parse([]string{"--help"})
	if err == nil || err.Error() != "hook failure" {
		t.Fatalf("expected hook failure, got %v", err)
	}
}
