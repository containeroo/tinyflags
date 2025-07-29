package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

func main() {
	tf := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
	tf.Version("1.0.1")

	tf.Bool("debug", false, "debug mode").Strict()

	z := tf.DynamicGroup("tcp").Title("TCP").Note("this is a group note")
	z.Bool("verbose", false, "verbose mode").Strict()
	z.StringSlice("addresses", []string{}, "API address")
	z.Int("port", 8080, "API port")

	a := tf.DynamicGroup("http")
	a.String("address", "", "API address")
	a.Bool("verbose", false, "verbose mode").Strict()
	a.String("log-level", "info", "log level").
		Choices("debug", "info", "warn", "error")
	a.Int("port", 8080, "API port")
	a.SortFlags()
	a.Title("HTTP")
	a.Description("this is a group description")

	tf.SortedGroups()

	// parse two dynamic flags
	args := []string{
		"--http.alpha.address=127.0.0.1",
		"--http.alpha.port=8080",
		"--http.beta.address=10.0.0.1",
		"--tcp.beta.addresses=10.0.0.1",
		"--tcp.beta.port=9090",
	}
	args = append(args, os.Args[1:]...)

	if err := tf.Parse(args); err != nil {
		if tinyflags.IsHelpRequested(err) || tinyflags.IsVersionRequested(err) {
			fmt.Println(err.Error()) // nolint:errcheck
			os.Exit(0)
		}
		panic(err)
	}

	type HTTPChecker struct {
		Name    string
		Address string
		Port    int
	}

	var httpCheckers []HTTPChecker

	i, _ := a.Lookup("verbose")
	if i == nil {
		panic("verbose not set")
	}
	fmt.Println(i.GetAny("alpha"))

	for _, group := range tf.DynamicGroups() {
		for _, id := range group.Instances() {
			name := group.Name()

			fmt.Println("id:", id)

			switch name {
			case "http":
				addr := tinyflags.GetOrDefaultDynamic[string](group, id, "address")
				port := tinyflags.GetOrDefaultDynamic[int](group, id, "port")
				checker := HTTPChecker{
					Name:    name,
					Address: addr,
					Port:    port,
				}
				httpCheckers = append(httpCheckers, checker)
			}
		}
	}

	fmt.Println(httpCheckers)
}
