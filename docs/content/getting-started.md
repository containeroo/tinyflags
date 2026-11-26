# Getting started

Go 1.25 or later is required by the module.

## Install

```bash
go get github.com/containeroo/tinyflags
```

## Quickstart

```go
package main

import (
    "fmt"
    "os"

    "github.com/containeroo/tinyflags"
)

func main() {
    fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
    fs.EnvPrefix("MYAPP")
    fs.Version("v1.2.3")

    host := fs.String("host", "localhost", "Server hostname").Required().Value()
    port := fs.Int("port", 8080, "Port to bind").Short("p").Value()
    debug := fs.Bool("debug", false, "Enable debug logging").Short("d").Value()
    tags := fs.StringSlice("tag", nil, "Optional tags").Value()

    if err := fs.Parse(os.Args[1:]); err != nil {
        if tinyflags.IsHelpRequested(err) || tinyflags.IsVersionRequested(err) {
            fmt.Println(err)
            return
        }
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    fmt.Println("Host:", *host)
    fmt.Println("Port:", *port)
    fmt.Println("Debug:", *debug)
    fmt.Println("Tags:", *tags)
}
```

If starting from an empty directory, run `go mod init example.com/app` before `go get`. Save the program as `main.go` in the module and run:

```sh
go run . --host=example.com -p 9090 -d --tag=api,web
MYAPP_HOST=example.com go run .
go run . --help
```

See [Flags and types](flags.md) for more options and [Commands](commands.md) for subcommands.
