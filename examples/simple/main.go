package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

func main() {
	// Usually you would parse args from os.Args[1:]
	// but for this example we'll just hard-code them.
	args := []string{
		"--port=9000",
		"--host=example.com",
		"-dtrue",
	}

	args = append(args, os.Args[1:]...) // append remaining args

	fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
	fs.Version("v1.0")    // optional, enables -v, --version
	fs.EnvPrefix("MYAPP") // optional, enables --env-key for all flags

	host := fs.String("host", "localhost", "host to use").
		Required().
		Value()

	port := fs.Int("port", 8080, "port to listen on").
		Env("MYAPP_CUSTOM_PORT"). // overrides default env key (otherwise "MYAPP_PORT")
		Required().
		Value()

	debug := fs.BoolP("debug", "d", false, "enable debug mode").
		Strict(). // strict bools require explicit value (--debug=true | --debug=false)
		Value()

	tags := fs.StringSlice("tag", []string{}, "list of tags").
		Value()

	loglevel := fs.String("log-level", "info", "log level to use").
		Choices("debug", "info", "warn", "error").
		Value()

	if err := fs.Parse(args); err != nil {
		if tinyflags.IsHelpRequested(err) || tinyflags.IsVersionRequested(err) {
			fmt.Fprint(os.Stdout, err.Error()+"\n") // nolint:errcheck
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err.Error()+"\n") //nolint:errcheck
		os.Exit(1)
	}

	fmt.Println("Host:", *host)
	fmt.Println("Port:", *port)
	fmt.Println("Debug:", *debug)
	fmt.Println("Tags:", *tags)
	fmt.Println("Loglevel:", *loglevel)
}
