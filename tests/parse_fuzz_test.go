package tinyflags_test

import (
	"strings"
	"testing"

	"github.com/containeroo/tinyflags"
)

func FuzzParseStaticArgs(f *testing.F) {
	seeds := [][]string{
		{"--name=alice", "--count=3", "--debug"},
		{"-d", "--count", "7", "positional"},
		{"--name=a,b", "--unknown"},
		{"--", "--name=literal", "tail"},
	}
	for _, seed := range seeds {
		f.Add(strings.Join(seed, "\n"))
	}

	f.Fuzz(func(t *testing.T, joined string) {
		args := fuzzArgs(joined)

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		fs.String("name", "", "name")
		fs.Int("count", 0, "count")
		fs.Bool("debug", false, "debug").Short("d")
		fs.StringSlice("tag", nil, "tags")
		fs.RequirePositional(0)

		_ = fs.Parse(args)
	})
}

func FuzzParseDynamicArgs(f *testing.F) {
	seeds := [][]string{
		{"--http.alpha.addr=127.0.0.1", "--http.alpha.port=8080"},
		{"--http.beta.tags=a,b,c", "--http.beta.enabled"},
		{"--http.alpha.addr", "localhost"},
		{"--http.alpha.unknown=value"},
	}
	for _, seed := range seeds {
		f.Add(strings.Join(seed, "\n"))
	}

	f.Fuzz(func(t *testing.T, joined string) {
		args := fuzzArgs(joined)

		fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
		http := fs.DynamicGroup("http")
		http.String("addr", "", "addr")
		http.Int("port", 80, "port")
		http.StringSlice("tags", nil, "tags")
		http.Bool("enabled", false, "enabled")

		_ = fs.Parse(args)
	})
}

func fuzzArgs(joined string) []string {
	if joined == "" {
		return nil
	}

	parts := strings.Split(joined, "\n")
	args := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(part) > 64 {
			part = part[:64]
		}
		args = append(args, part)
	}
	return args
}
