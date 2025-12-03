package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

func main() {
	fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)

	debugFlag := fs.Bool("debug", false, "Enable debug logs").Short("d")

	if err := fs.Parse(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch {
	case debugFlag.Changed() && *debugFlag.Value():
		fmt.Println("debug explicitly enabled via --debug")
	case debugFlag.Changed() && !*debugFlag.Value():
		fmt.Println("debug explicitly disabled via --no-debug")
	default:
		fmt.Println("debug left at default (false)")
	}
}
