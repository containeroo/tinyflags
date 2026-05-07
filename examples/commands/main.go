package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

// main demonstrates nested subcommands with persistent and local flags.
func main() {
	app := tinyflags.NewCommand("app", tinyflags.ExitOnError)
	verbose := app.Globals().Bool("verbose", false, "Enable verbose logging").Value()

	serve := app.Command("serve", "Run the HTTP server")
	port := serve.Int("port", 8080, "HTTP port").Value()

	admin := app.Command("admin", "Administrative commands")
	audit := admin.Globals().Bool("audit", false, "Enable audit mode").Value()

	users := admin.Command("users", "Manage user accounts")
	userName := users.String("name", "alice", "User name").Value()

	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"admin", "--audit", "users", "--name=bob", "--verbose"}
	}

	if err := app.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	selected := app.SelectedCommand()
	fmt.Printf("selected command: %s\n", selected.FullName())
	fmt.Printf("verbose: %t\n", *verbose)

	switch selected {
	case serve:
		fmt.Printf("port: %d\n", *port)
	case users:
		fmt.Printf("audit: %t\n", *audit)
		fmt.Printf("user: %s\n", *userName)
	default:
		fmt.Println("no command-specific output")
	}
}
