package tinyflags

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintUsageModes(t *testing.T) {
	t.Parallel()

	t.Run("PrintNone", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("myapp", ContinueOnError)
		var buf bytes.Buffer
		fs.Parse([]string{}) // nolint:errcheck
		fs.PrintUsage(&buf, PrintNone)
		assert.Equal(t, "Usage: myapp\n", buf.String())
	})

	t.Run("PrintFlags", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("myapp", ContinueOnError)
		fs.Bool("debug", false, "debug mode")
		fs.String("hidden", "", "hidden flag").Hidden()
		var buf bytes.Buffer
		fs.Parse([]string{}) // nolint:errcheck
		fs.PrintUsage(&buf, PrintFlags)
		assert.Equal(t, "Usage: myapp [flags]\n", buf.String())
	})

	t.Run("PrintShort", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("myapp", ContinueOnError)
		fs.StringP("foo", "f", "", "foo flag")
		var buf bytes.Buffer
		fs.PrintUsage(&buf, PrintShort)
		out := buf.String()
		assert.Contains(t, out, "-f")
	})

	t.Run("PrintLong", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("myapp", ContinueOnError)
		fs.Version("v1.2.3")
		fs.Sorted(true)
		fs.StringP("foo", "f", "", "foo flag")

		var buf bytes.Buffer
		fs.PrintUsage(&buf, PrintLong)
		out := buf.String()
		assert.Contains(t, out, "--foo")
	})

	t.Run("PrintBoth", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("myapp", ContinueOnError)
		fs.Version("v1.2.3")
		fs.Sorted(true)
		fs.StringP("foo", "f", "", "foo flag")
		fs.Parse([]string{}) // nolint:errcheck
		var buf bytes.Buffer
		fs.PrintUsage(&buf, PrintBoth)
		out := buf.String()
		assert.Contains(t, out, "-f|--foo")
	})
}

func TestPrintTitleAndDescriptionAndNotes(t *testing.T) {
	t.Parallel()
	t.Run("PrintTitle", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("myapp", ContinueOnError)
		fs.Title("MyApp Flags")
		var buf bytes.Buffer
		fs.PrintTitle(&buf)
		assert.Equal(t, "MyApp Flags\n", buf.String())
	})

	t.Run("PrintDescription", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("myapp", ContinueOnError)
		fs.Description("Use this tool for great success.")
		var buf bytes.Buffer
		fs.PrintDescription(&buf, 80)
		assert.Contains(t, buf.String(), "great success")
	})

	t.Run("PrintNotes", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("myapp", ContinueOnError)
		fs.Note("See also: docs.")
		var buf bytes.Buffer
		fs.PrintNotes(&buf, 80)
		assert.Contains(t, buf.String(), "See also")
	})
}

func TestPrintDefaultsVariants(t *testing.T) {
	t.Parallel()
	t.Run("Sorted and Required", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("sorted", ContinueOnError)
		fs.Sorted(true)
		fs.String("zeta", "", "last flag").Required()
		fs.String("alpha", "", "first flag")
		var buf bytes.Buffer
		fs.SetOutput(&buf)
		fs.PrintDefaults()
		out := buf.String()
		assert.True(t, strings.Index(out, "--alpha") < strings.Index(out, "--zeta"))
		assert.Contains(t, out, "Required")
	})

	t.Run("Unsorted and Deprecated", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("unsorted", ContinueOnError)
		fs.Version("v1.2.3")
		fs.Sorted(false)
		fs.String("foo", "", "do not use").Deprecated("use --bar instead")
		var buf bytes.Buffer
		fs.SetOutput(&buf)
		fs.PrintDefaults()
		assert.Contains(t, buf.String(), "DEPRECATED")
	})

	t.Run("EnvPrefix shown", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("env", ContinueOnError)
		fs.Version("v1.2.3")
		fs.EnvPrefix("APP")
		fs.String("hidden", "", "hidden flag").Hidden()
		fs.String("debug", "", "debug mode")
		var buf bytes.Buffer
		fs.SetOutput(&buf)
		fs.Parse([]string{}) // nolint:errcheck
		fs.PrintDefaults()
		assert.Contains(t, buf.String(), "Env: APP_DEBUG")
	})

	t.Run("Bool strict default", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("booltest", ContinueOnError)
		fs.Bool("check", false, "strict flag").Strict()
		var buf bytes.Buffer
		fs.SetOutput(&buf)
		fs.PrintDefaults()
		assert.Contains(t, buf.String(), "Default: false")
	})

	t.Run("Choices and metavar", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("choices", ContinueOnError)
		fs.String("mode", "m", "set mode").Choices("a", "b").Metavar("MODE")
		var buf bytes.Buffer
		fs.SetOutput(&buf)
		fs.PrintDefaults()
		assert.Contains(t, buf.String(), "MODE")
		assert.Contains(t, buf.String(), "Allowed: a, b")
	})

	t.Run("Slice with implicit metavar", func(t *testing.T) {
		t.Parallel()
		fs := NewFlagSet("slice", ContinueOnError)
		fs.StringSlice("list", []string{}, "a list")
		var buf bytes.Buffer
		fs.SetOutput(&buf)
		fs.PrintDefaults()
		assert.Contains(t, buf.String(), "LIST...")
	})
}
