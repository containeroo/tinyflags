package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

// main demonstrates basic tinyflags usage.
func main() {
	fs := tinyflags.NewFlagSet("hello", tinyflags.ExitOnError)

	name := fs.String("name", "world", "Who to greet").Value()

	silent := fs.Bool("silent", false, "Silent mode").
		Strict().
		Value()

	debug := fs.Bool("debug", false, "Enable debug logging").Value()

	verbose := fs.Counter("verbose", 0, "Enable verbose logging").
		Short("v").
		Value()

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Hello, %s!\n", *name)
	fmt.Printf("verbosity: %d\n", *verbose)
	fmt.Printf("debug: %t\n", *debug)
	fmt.Printf("silent: %t\n", *silent)
}
