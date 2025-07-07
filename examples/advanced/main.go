package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/containeroo/tinyflags"
)

type Config struct {
	Port     int
	Host     string
	Addr     string
	HostIP   net.IP
	Verbose  bool
	Insecure bool
	Debug    bool
	IpAddr   string
	LogLevel string
	Paths    []string
}

func parseArgs(args []string) (*Config, error) {
	tf := tinyflags.NewFlagSet("test.exe", tinyflags.ContinueOnError)
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
	}
	showHelp := tf.Bool("help", false, "show help"). // Register own without shorthand
								Value()

	port := tf.Int("port", 8080, "port to use").
		Env("MYAPP_CUSTOM_PORT").
		Required().
		Value()

	host := tf.StringP("host", "h", "localhost", "host to use").
		Required().
		Value()

	defaultAddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	addr := tf.ListenAddrP("addr", "a", defaultAddr, "listen address").Value()

	hostip := tf.IP("host-ip", net.ParseIP("10.0.10.8"), "host ip to use. Must be in range 10.0.10.0/24").
		Validator(func(ip net.IP) error {
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

	debug := tf.BoolP("debug", "d", false, "debug mode").
		Value()

	insecure := tf.BoolP("insecure", "i", false, "insecure mode").
		Value()

	verbose := tf.BoolP("verbose", "v", false, "verbose mode").
		Strict().
		Value()

	if err := tf.Parse(args); err != nil {
		return nil, err
	}

	if showHelp {
		var buf bytes.Buffer
		tf.SetOutput(&buf)
		tf.Usage()
		return &Config{}, tinyflags.RequestHelp(buf.String())
	}

	// Positional arguments are captured after all flags
	paths := tf.Args()

	return &Config{
		Port:     port,
		Host:     host,
		HostIP:   hostip,
		Addr:     addr.String(),
		Verbose:  verbose,
		Insecure: insecure,
		Debug:    debug,
		LogLevel: loglevel,
		Paths:    paths,
	}, nil
}

func main() {
	args := []string{
		"--port=9000",
		"--host=example.com",
		"--host-ip=10.0.10.12",
		"-vtrue",
		"-di",
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
	fmt.Println("positional:", strings.Join(cfg.Paths, ", "))
}
