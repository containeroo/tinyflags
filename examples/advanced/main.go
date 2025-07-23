package main

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/containeroo/tinyflags"
)

type Config struct {
	Port           int
	Host           string
	ListenAddr     string
	SchemaHostPort string
	HostIP         net.IP
	Verbose        int
	Insecure       bool
	LogLevel       string
	Paths          []string
}

func parseArgs(args []string) (*Config, error) {
	tf := tinyflags.NewFlagSet("test.exe", tinyflags.ContinueOnError)
	tf.Authors("me@containeroo.ch")
	tf.EnvPrefix("MYAPP")    // optional, enables --env-key for all flags
	tf.Version("v1.2.3")     // optional, enables -v, --version
	tf.DisableHelp()         // optional, disables automatic help flag registration
	tf.DescriptionIndent(45) // optional, sets indentation for descriptions

	// Since we disabled help, we need to define our own usage function
	tf.Usage = func() {
		out := tf.Output()
		tf.PrintUsage(out, tinyflags.PrintBoth)
		tf.PrintTitle(out)
		tf.PrintDescription(out, 80)
		tf.PrintDefaults()
		tf.PrintNotes(out, 80)
		tf.PrintAuthors(out)
	}
	showHelp := tf.Bool("help", false, "show help"). // Register own without shorthand
								Value()

	port := tf.Int("port", 8080, "port to use").
		Env("MYAPP_CUSTOM_PORT").
		Required().
		Value()

	host := tf.String("host", "localhost", "host to use").
		Short("h").
		Required().
		Value()

	listenAddr := tf.TCPAddr("listen-addr", &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	}, "listen address to use").
		Short("l").
		Value()

	schemaHostPort := tf.String("schema-host-port", "scheme://host:port", "schema://host:port").
		Validate(func(s string) error {
			u, err := url.Parse(s)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return fmt.Errorf("invalid scheme://host:port format")
			}
			return nil
		}).Value()

	hostip := tf.IP("host-ip", net.ParseIP("10.0.10.8"), "host ip to use. Must be in range 10.0.10.0/24").
		Validate(func(ip net.IP) error {
			_, ipNet, _ := net.ParseCIDR("10.0.10.0/24")
			if !ipNet.Contains(ip) {
				return fmt.Errorf("must be in range %s", ipNet.String())
			}
			return nil
		}).
		Value()

	loglevel := tf.String("log-level", "info", "log level to use").
		Choices("debug", "info", "warn", "error").
		Value()

	insecure := tf.Bool("insecure", false, "insecure mode").
		Short("i").
		Value()

	verbose := tf.Counter("verbose", "verbose mode").
		Short("v").
		Value()

	if err := tf.Parse(args); err != nil {
		return nil, err
	}

	if *showHelp {
		var buf bytes.Buffer
		tf.SetOutput(&buf)
		tf.Usage()
		return &Config{}, tinyflags.RequestHelp(buf.String())
	}

	// Positional arguments are captured after all flags
	paths := tf.Args()

	return &Config{
		Port:           *port,
		Host:           *host,
		ListenAddr:     (*listenAddr).String(),
		SchemaHostPort: *schemaHostPort,
		HostIP:         *hostip,
		Verbose:        *verbose,
		Insecure:       *insecure,
		LogLevel:       *loglevel,
		Paths:          paths,
	}, nil
}

func main() {
	args := []string{
		"--port=9000",
		"--host=example.com",
		"--host-ip=10.0.10.12",
		"-vv",
		"--log-level=debug",
		"/first/path", "/second/path",
	}
	args = append(args, os.Args[1:]...)

	cfg, err := parseArgs(args)
	if err != nil {
		if tinyflags.IsHelpRequested(err) || tinyflags.IsVersionRequested(err) {
			fmt.Fprint(os.Stdout, err.Error()) // nolint:errcheck
			os.Exit(0)
		}
		fmt.Fprint(os.Stderr, err.Error()+"\n") // nolint:errcheck
		os.Exit(1)
	}

	fmt.Println("port:", cfg.Port)
	fmt.Println("host:", cfg.Host)
	fmt.Println("verbose:", cfg.Verbose)
	fmt.Println("insecure:", cfg.Insecure)
	fmt.Println("log-level:", cfg.LogLevel)
	fmt.Println("schema host port", cfg.SchemaHostPort)
	fmt.Println("positional:", strings.Join(cfg.Paths, ", "))
}
