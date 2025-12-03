package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

func main() {
	fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)

	debugFlag := fs.Bool("debug", false, "Enable debug logs").Short("d").OneOfGroup("debug")
	noDebugFlag := fs.Bool("no-debug", false, "Disable debug logs").Short("n").OneOfGroup("debug")

	if err := fs.Parse(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	debug, set := tinyflags.FirstChanged(false, debugFlag, noDebugFlag)
	source := "default"
	if set {
		if debugFlag.Changed() {
			source = "--debug/-d"
		} else if noDebugFlag.Changed() {
			source = "--no-debug/-n"
		}
	}
	fmt.Printf("debug enabled: %t (source: %s)\n", debug, source)
}
