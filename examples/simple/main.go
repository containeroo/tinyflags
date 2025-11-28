package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

func main() {
	fs := tinyflags.NewFlagSet("hello", tinyflags.ExitOnError)

	name := fs.String("name", "world", "Who to greet").Value()
	verbose := fs.Bool("verbose", false, "Enable verbose logging").Value()

	if err := fs.Parse(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("parsed flags successfully")
	}
	fmt.Printf("Hello, %s!\n", *name)
}
