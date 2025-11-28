package engine

import (
	"strconv"
	"testing"

	"github.com/containeroo/tinyflags/internal/utils"
)

func TestRegisterStaticScalarInt(t *testing.T) {
	t.Parallel()

	t.Run("defaultIsApplied", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)

		var value int
		RegisterStaticScalar(fs, &value, "num", "usage", 7, strconv.Atoi, strconv.Itoa)

		if value != 7 {
			t.Fatalf("expected default 7, got %d", value)
		}
	})

	t.Run("parsesInput", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)

		num := new(int)
		RegisterStaticScalar(fs, num, "num", "usage", 0, strconv.Atoi, strconv.Itoa)

		if err := fs.Parse([]string{"--num=5"}); err != nil {
			t.Fatalf("parse failed: %v", err)
		}

		if *num != 5 {
			t.Fatalf("expected parsed value 5, got %d", *num)
		}
	})
}

func TestRegisterStaticSliceString(t *testing.T) {
	t.Parallel()

	t.Run("defaultIsApplied", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)

		var names []string
		RegisterStaticSlice(fs, &names, "names", "usage", []string{"a", "b"}, utils.ParseString, utils.FormatString, ";")

		if got := len(names); got != 2 || names[0] != "a" || names[1] != "b" {
			t.Fatalf("expected default [a b], got %v", names)
		}
	})

	t.Run("parsesInputWithDelimiter", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)

		names := new([]string)
		RegisterStaticSlice(fs, names, "names", "usage", nil, utils.ParseString, utils.FormatString, ";")

		if err := fs.Parse([]string{"--names=a; b"}); err != nil {
			t.Fatalf("parse failed: %v", err)
		}

		want := []string{"a", "b"}
		got := *names
		if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})
}
