package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

// main demonstrates sub-subcommands with inherited persistent flags.
func main() {
	app := tinyflags.NewCommand("app", tinyflags.ExitOnError)
	verbose := app.Globals().Bool("verbose", false, "Enable verbose logging").Value()

	cluster := app.Command("cluster", "Cluster operations")
	region := cluster.Globals().String("region", "eu-central", "Target region").Value()

	node := cluster.Command("node", "Node operations")
	drain := node.Command("drain", "Drain one node")
	nodeName := drain.String("name", "worker-1", "Node name").Value()
	force := drain.Bool("force", false, "Force draining").Value()

	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"cluster", "--region=us-east", "node", "drain", "--name=worker-42", "--force", "--verbose"}
	}

	if err := app.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("selected command: %s\n", app.SelectedCommand().FullName())
	fmt.Printf("verbose: %t\n", *verbose)
	fmt.Printf("region: %s\n", *region)
	fmt.Printf("node: %s\n", *nodeName)
	fmt.Printf("force: %t\n", *force)
}
