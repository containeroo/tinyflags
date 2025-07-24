package main

import (
	"fmt"

	"github.com/containeroo/tinyflags"
)

func main() {
	fmt.Println("start")

	fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
	fs.Version("1.0.0")

	http := fs.DynamicGroup("http")
	http.String("address", "", "API address")
	http.Int("port", 8080, "API port")

	tcp := fs.DynamicGroup("tcp")
	tcp.String("address", "", "API address")
	tcp.Int("port", 8080, "API port")

	// parse two dynamic flags
	if err := fs.Parse([]string{
		"--http.alpha.address=127.0.0.1",
		"--http.alpha.port=8080",
		"--http.beta.address=10.0.0.1",
		"--tcp.beta.address=10.0.0.1",
		"--tcp.beta.port=9090",
	}); err != nil {
		panic(err)
	}

	type HTTPChecker struct {
		Name    string
		Address string
		Port    int
	}

	var httpCheckers []HTTPChecker

	for _, group := range fs.DynamicGroups() {
		for _, id := range group.Instances() {
			name := group.Name()
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
