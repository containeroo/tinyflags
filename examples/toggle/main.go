package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

// main demonstrates paired toggle flags.
func main() {
	fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)

	debugFlag := fs.Bool("debug", false, "Enable debug logs").Short("d").OneOfGroup("debug")
	noDebugFlag := fs.Bool("no-debug", false, "Disable debug logs").Short("n").OneOfGroup("debug")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	debug, _ := tinyflags.FirstChanged(false, debugFlag, noDebugFlag)
	origin := fs.Origin("debug")
	if noDebugFlag.Changed() {
		origin = fs.Origin("no-debug")
	}
	fmt.Printf("debug enabled: %t (source: %s)\n", debug, origin)
}
