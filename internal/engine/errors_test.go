package engine

import "testing"

func TestHelpAndVersionErrors(t *testing.T) {
	t.Parallel()

	t.Run("helpSentinel", func(t *testing.T) {
		t.Parallel()

		err := RequestHelp("help me")
		if !IsHelpRequested(err) {
			t.Fatalf("expected help sentinel")
		}
		if err.Error() != "help me" {
			t.Fatalf("unexpected message: %q", err.Error())
		}
	})

	t.Run("versionSentinel", func(t *testing.T) {
		t.Parallel()

		err := RequestVersion("v1.2.3")
		if !IsVersionRequested(err) {
			t.Fatalf("expected version sentinel")
		}
		if err.Error() != "v1.2.3" {
			t.Fatalf("unexpected version: %q", err.Error())
		}
	})
}
