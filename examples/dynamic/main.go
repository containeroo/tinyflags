package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/containeroo/tinyflags"
)

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

	// Global unknown flag handler to ignore extra args
	fs.OnUnknownFlag(func(name string) error { return nil })

	if err := fs.Parse(nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, id := range svc.Instances() {
		fmt.Printf("[%s] addr=%s ports=%v\n", id, addr.MustGet(id), ports.MustGet(id))
	}
}
