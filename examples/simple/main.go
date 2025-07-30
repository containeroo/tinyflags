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
		"--user", "admin",
	}

	args = append(args, os.Args[1:]...) // append remaining args

	tf := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
	tf.Version("v1.0")    // optional, enables -v, --version
	tf.EnvPrefix("MYAPP") // optional, enables --env-key for all flags

	host := tf.String("host", "localhost", "host to use").
		Required().
		Value()

	port := tf.Int("port", 8080, "port to listen on").
		Env("MYAPP_CUSTOM_PORT"). // overrides default env key (otherwise "MYAPP_PORT")
		Required().
		Value()

	debug := tf.Bool("debug", false, "enable debug mode").
		Strict(). // strict bools require explicit value (--debug=true | --debug=false)
		Short("d").
		Value()

	log := tf.Bool("log", false, "enable logging").MutualExlusive("logging").Short("l").Value()
	noLog := tf.Bool("no-log", false, "disable logging").MutualExlusive("logging").Short("L").Value()
	tf.GetMutualGroup("logging").Hidden()

	tags := tf.StringSlice("tag", []string{}, "list of tags").
		Value()

	loglevel := tf.String("log-level", "info", "log level to use").
		Choices("debug", "info", "warn", "error").
		Value()

	user := tf.String("user", "admin", "user to use").RequireTogether("credentials").Value()
	pw := tf.String("password", "", "password to use").RequireTogether("credentials").Value()

	if err := tf.Parse(args); err != nil {
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
	fmt.Println("Log:", *log)
	fmt.Println("No Log:", *noLog)
	fmt.Println("User:", *user)
	fmt.Println("Password:", *pw)
}
