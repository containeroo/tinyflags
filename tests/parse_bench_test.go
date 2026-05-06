package tinyflags_test

import (
	"fmt"
	"testing"

	"github.com/containeroo/tinyflags"
)

// BenchmarkParseStaticFlags benchmarks parsing common static flags.
func BenchmarkParseStaticFlags(b *testing.B) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	fs.String("host", "localhost", "host")
	fs.Int("port", 8080, "port")
	fs.Bool("debug", false, "debug")
	fs.StringSlice("tag", nil, "tags")

	args := []string{
		"--host=example.com",
		"--port=9090",
		"--debug",
		"--tag=alpha,beta,gamma",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := fs.Parse(args); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseDynamicFlags benchmarks parsing common dynamic flags.
func BenchmarkParseDynamicFlags(b *testing.B) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	http := fs.DynamicGroup("http")
	http.String("addr", "", "addr")
	http.Int("port", 80, "port")
	http.StringSlice("tags", nil, "tags")

	args := []string{
		"--http.alpha.addr=127.0.0.1",
		"--http.alpha.port=8080",
		"--http.alpha.tags=api,blue",
		"--http.beta.addr=10.0.0.2",
		"--http.beta.port=9090",
		"--http.beta.tags=jobs,green",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := fs.Parse(args); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseEnvOverrides benchmarks loading values from environment variables.
func BenchmarkParseEnvOverrides(b *testing.B) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	fs.SetGetEnvFn(func(key string) string {
		values := map[string]string{
			"APP_HOST": "env.example.com",
			"APP_PORT": "9443",
			"APP_TAG":  "blue,green",
		}
		return values[key]
	})
	fs.String("host", "localhost", "host")
	fs.Int("port", 8080, "port")
	fs.StringSlice("tag", nil, "tags")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := fs.Parse(nil); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUsageRendering benchmarks help rendering for larger flag sets.
func BenchmarkUsageRendering(b *testing.B) {
	fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
	fs.EnvPrefix("APP")
	fs.Description("Benchmark usage output for many flags and help annotations.")
	for i := 0; i < 40; i++ {
		fs.String(fmt.Sprintf("name-%d", i), "default", "test flag").Section("Static")
	}
	http := fs.DynamicGroup("http")
	http.Description("Dynamic HTTP targets.")
	for i := 0; i < 10; i++ {
		http.String(fmt.Sprintf("field-%d", i), "", "dynamic field")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fs.PrintUsage(ioDiscard{}, tinyflags.PrintFlags)
		fs.PrintDescription(ioDiscard{}, 0, 120)
		fs.PrintStaticDefaults(ioDiscard{}, 2, 30, 120)
		fs.PrintDynamicDefaults(ioDiscard{}, 2, 30, 120)
	}
}

type ioDiscard struct{}

// Write discards all bytes and reports success.
func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
