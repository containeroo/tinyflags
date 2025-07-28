package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

func main() {
	fmt.Println("start")

	fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
	fs.Version("1.0.1")

	fs.Bool("debug", false, "debug mode").Strict()

	http := fs.DynamicGroup("http")
	http.String("address", "", "API address")
	http.Bool("verbose", false, "verbose mode").Strict()
	http.String("log-level", "info", "log level").
		Choices("debug", "info", "warn", "error")
	http.Int("port", 8080, "API port")
	http.SortFlags()

	tcp := fs.DynamicGroup("tcp")
	tcp.StringSlice("addresses", []string{}, "API address")
	tcp.Int("port", 8080, "API port")
	tcp.Bool("verbose", false, "verbose mode")

	// parse two dynamic flags
	args := []string{
		"--http.alpha.address=127.0.0.1",
		"--http.alpha.port=8080",
		"--http.beta.address=10.0.0.1",
		"--tcp.beta.addresses=10.0.0.1",
		"--tcp.beta.port=9090",
	}
	args = append(args, os.Args[1:]...)

	if err := fs.Parse(args); err != nil {
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

	i, _ := http.Lookup("verbose")
	if i == nil {
		panic("verbose not set")
	}
	fmt.Println(i.GetAny("alpha"))

	for _, group := range fs.DynamicGroups() {
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
