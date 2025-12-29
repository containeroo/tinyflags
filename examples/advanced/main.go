package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/containeroo/tinyflags"
)

func main() {
	fs := tinyflags.NewFlagSet("advanced", tinyflags.ExitOnError)
	fs.Authors("Tinyflags Team")
	fs.Title("Advanced Example")
	fs.Description("Demonstrates slices, sections, and advanced flag options.")
	fs.Note("This example demonstrates the following features:\n\n" +
		"- Slices\n" +
		"- Sections\n" +
		"- Finalize\n" +
		"- FinalizeDefaultValue\n" +
		"- Choices\n" +
		"- StrictDelimiter\n" +
		"- AllowEmpty\n" +
		"- Default values\n" +
		"- Environment variables\n" +
		"- Short options\n" +
		"- Debug logs\n")

	env := fs.String("env", "dev", "Deployment environment").
		Choices("dev", "staging", "prod").
		Section("General")
	envVal := env.Value()

	timeout := fs.Duration("timeout", 5*time.Second, "request timeout").
		Section("Networking")
	timeoutVal := timeout.Value()

	tags := fs.StringSlice("tag", nil, "Tags to attach").
		Delimiter(",").
		StrictDelimiter().
		AllowEmpty().
		Section("Metadata")
	tagsVal := tags.Value()

	configDir := fs.String("config-dir", " /etc/app ", "Config directory").
		Finalize(strings.TrimSpace).
		FinalizeDefaultValue().
		Section("General")
	configDirVal := configDir.Value()

	webhook := fs.String("webhook", "", "Webhook URL").
		Finalize(func(u string) string {
			if u == "" {
				return u
			}
			return strings.TrimSuffix(u, "/") + "/"
		}).
		Section("Networking")

	debug := fs.Bool("debug", false, "Enable debug logs").Short("d")
	noDebug := fs.Bool("no-debug", false, "Disable debug logs").Short("n")

	if err := fs.Parse(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("env:", *envVal)
	fmt.Println("timeout:", *timeoutVal)
	fmt.Println("tags:", *tagsVal)
	fmt.Println("config-dir:", *configDirVal)
	fmt.Println("webhook:", *webhook)
	debugVal, set := tinyflags.FirstChanged(false, debug, noDebug)
	source := "default"
	if set {
		switch {
		case debug.Changed():
			source = "--debug/--d"
		case noDebug.Changed():
			source = "--no-debug/--n"
		}
	}
	fmt.Printf("debug: %t (source: %s)\n", debugVal, source)
}
