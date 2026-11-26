# tinyflags

[![Go Reference](https://pkg.go.dev/badge/github.com/containeroo/tinyflags.svg)](https://pkg.go.dev/github.com/containeroo/tinyflags)

A minimal, fast, and extensible CLI flag-parsing library for Go. Zero third-party runtime dependencies, generic flag builders, environment variable overrides, subcommands, validation, and rich help output.

## Getting started

Requires Go 1.25 or later.

```sh
go get github.com/containeroo/tinyflags
```

```go
package main

import (
    "fmt"
    "os"

    "github.com/containeroo/tinyflags"
)

func main() {
    fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
    fs.EnvPrefix("APP")
    fs.NoPositional()

    port := fs.Int("port", 8080, "Listen port").Short("p").
        Validate(tinyflags.Between(1, 65535)).Value()

    if err := fs.Parse(os.Args[1:]); err != nil {
        if tinyflags.IsHelpRequested(err) {
            fmt.Println(err)
            return
        }
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    fmt.Println(*port)
}
```

```sh
go run . --port=9090
APP_PORT=9090 go run .
go run . --help
```

## Documentation

Full guides, examples and reference: **[containeroo.github.io/tinyflags](https://containeroo.github.io/tinyflags/)**.

For exported Go types and signatures, see [pkg.go.dev](https://pkg.go.dev/github.com/containeroo/tinyflags).

## License

Licensed under the [Apache License, Version 2.0](LICENSE).
