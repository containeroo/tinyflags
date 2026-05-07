package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

// main demonstrates grouped flag constraints.
func main() {
	fs := tinyflags.NewFlagSet("groups", tinyflags.ExitOnError)
	fs.Description("Choose exactly one authentication method.")

	email := fs.String("email", "", "Login email").
		AllOrNone("password-auth")
	password := fs.String("password", "", "Login password").
		AllOrNone("password-auth")
	bearer := fs.String("bearer-token", "", "Bearer token").
		OneOfGroup("auth")

	fs.GetOneOfGroup("auth").
		Required().
		Title("Authentication Method")
	fs.AttachGroupToOneOf("auth", "password-auth")

	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{
			"--email=alice@example.com",
			"--password=super-secret",
		}
	}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch {
	case email.Changed():
		fmt.Printf("auth method: password\nemail: %s\npassword length: %d\n", *email.Value(), len(*password.Value()))
	case bearer.Changed():
		fmt.Printf("auth method: bearer\nbearer token prefix: %.6s...\n", *bearer.Value())
	default:
		fmt.Println("auth method: none")
	}
}
