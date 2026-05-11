package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/containeroo/tinyflags"
)

// main demonstrates dynamic tinyflags groups.
func main() {
	fs := tinyflags.NewFlagSet("dynamic", tinyflags.ExitOnError)

	// Dynamic group for services
	svc := fs.DynamicGroup("svc")
	addr := svc.String("addr", "", "Service address")
	addr.Required()
	addr.Validate(func(s string) error {
		if strings.Contains(s, "localhost") {
			return fmt.Errorf("localhost is not allowed")
		}
		return nil
	})
	ports := svc.IntSlice("port", nil, "Service ports").
		AllowEmpty().
		FinalizeWithID(func(id string, v int) int {
			// Example: offset port by length of ID
			return v + len(id)
		})
	headers := fs.StringSlice("headers", []string{}, "HTTP headers").Value()

	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{
			"--svc.api.addr=10.0.0.10",
			"--svc.api.port=8080,8443",
			"--svc.worker.addr=10.0.0.20",
			"--svc.worker.port=9000",
		}
	}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ids := svc.Instances()
	sort.Strings(ids)
	for _, id := range ids {
		fmt.Printf("[%s] addr=%s ports=%v headers=%v\n", id, addr.MustGet(id), ports.MustGet(id), *headers)
	}
}
