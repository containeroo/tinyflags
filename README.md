# tinyflags

[![Go Reference](https://pkg.go.dev/badge/github.com/containeroo/tinyflags.svg)](https://pkg.go.dev/github.com/containeroo/tinyflags)

A minimal, fast, and extensible CLI flag-parsing library for Go.
Zero dependencies, full generics support, advanced features.

## Features

- **Short & long flags** (`-d`, `--debug`)
- **Boolean strict mode** (`--flag=true/false`, `--no-flag`)
- **Environment overrides** (per-flag or `EnvPrefix`)
- **Required & deprecated flags**
- **Slice flags** (`[]T`) with custom delimiters
- **Choices & validation**
- **Mutual-exclusion groups**
- **Custom metavars** and **rich help** formatting
- **Dynamic flags** (`--group.id.field=value`)
- **TCP-addr**, **URL**, **IP**, **File**, **Duration**, and more built-in types

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
  fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
  fs.EnvPrefix("MYAPP")   // use MYAPP_… env vars
  fs.Version("v1.0.0")    // enable -v/--version

  host := fs.String("host", "localhost", "server host").Required().Value()
  port := fs.IntP("port", "p", 8080, "server port").Value()
  debug := fs.BoolP("debug", "d", false, "enable debug").Value()
  tags  := fs.StringSlice("tag", nil, "tags list").Value()

  if err := fs.Parse(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  fmt.Println("Host:", *host)
  fmt.Println("Port:", *port)
  fmt.Println("Debug:", *debug)
  fmt.Println("Tags:", *tags)
}
```

## Environment Variables

```bash
MYAPP_HOST=example.com MYAPP_PORT=9090 ./app --debug
```

Disable per-flag with:

```go
fs.Bool("internal", false, "").DisableEnv()
```

## Dynamic Flags

```go
package main

import (
  "fmt"
  "os"
  "github.com/containeroo/tinyflags"
)

func main() {
  fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
  dg := fs.DynamicGroup("http")

  // define per-instance flags
  portFlag     := dg.Int("port", "backend port")
  timeoutFlag  := dg.Duration("timeout", "request timeout")

  if err := fs.Parse(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  // iterate over all instance IDs seen
  for _, id := range dg.Instances() {
    port, _    := portFlag.Get(id)
    timeout, _ := timeoutFlag.Get(id)
    fmt.Printf("%s → port=%d, timeout=%s\n", id, port, timeout)
  }
}
```

Call it like:

```bash
./app --http.alpha.port=8080 --http.alpha.timeout=30s \
      --http.beta.port=9090 --http.beta.timeout=1m
```

## Help Output

```text
Usage: app [flags]

  --host HOST         server host (Default: localhost) (Env: MYAPP_HOST) (Required)
  -p, --port PORT     server port (Default: 8080) (Env: MYAPP_PORT)
  -d, --debug         enable debug (Env: MYAPP_DEBUG)
      --tag TAG...    tags list
  -v, --version       show version
```

## Supported Types

| Go Type               | Methods                                 |
| --------------------- | --------------------------------------- |
| `bool`                | `Bool`, `BoolP`, `BoolVarP`             |
| `int`                 | `Int`, `IntP`, `IntVarP`                |
| `string`              | `String`, `StringP`, `StringVarP`       |
| `[]string`            | `StringSlice`, `StringSliceP`, …        |
| `time.Duration`       | `Duration`, `DurationP`, `DurationVarP` |
| `net.IP` / `[]net.IP` | `IP`, `IPSlice`, …                      |
| `*net.TCPAddr`        | `TCPAddr`, `TCPAddrP`, `TCPAddrVarP`    |
| `url.URL`             | `URL`, `URLP`, `URLVarP`                |
| `*os.File`            | `File`, `FileP`, `FileVarP`             |

> All slice flags also support repeated usage and custom delimiters.

## FlagSet Methods

| Method                                                                | Description                            |
| --------------------------------------------------------------------- | -------------------------------------- |
| `NewFlagSet(name, mode)`                                              | create a new named FlagSet             |
| `Parse(args)`                                                         | parse flags and env vars               |
| `Version(s)`                                                          | set version string for `--version`     |
| `EnvPrefix(s)`                                                        | set global env-var prefix              |
| `Authors(s)`                                                          | add author info to help                |
| `Title(s)`                                                            | set help title                         |
| `Description(s)`                                                      | add a description paragraph            |
| `Note(s)`                                                             | append a note paragraph                |
| `DisableHelp()`, `DisableVersion()`                                   | disable built-in help/version flags    |
| `Sorted(bool)`                                                        | enable/disable sorted flag output      |
| `SetOutput(w)`, `Output()`                                            | set/get output writer                  |
| `IgnoreInvalidEnv(bool)`                                              | skip invalid env-var values            |
| `SetGetEnvFn(fn)`                                                     | customize env-var lookup               |
| `Globaldelimiter(s)`                                                  | set default slice-delimiter            |
| `RequirePositional(n)`                                                | require at least n positional args     |
| `Args()`, `Arg(i)`                                                    | retrieve positional args               |
| `GetGroup(name)`                                                      | get or create a mutual-exclusion group |
| `PrintDefaults()`                                                     | print all flags with defaults          |
| `PrintUsage(w, mode)`                                                 | print usage in specified mode          |
| `PrintTitle(w)`, `PrintNotes(w, width)`, `PrintDescription(w, width)` | various help sections                  |
| `PrintAuthors(w)`                                                     | print author info                      |
| `DynamicGroup(name)`                                                  | create/get a dynamic-flag group        |
| `DefaultDelimiter()`                                                  | get default delimiter for slice flags  |
| `RegisterDynamic(group,field,val)`                                    | register a custom dynamic flag         |
| `RegisterFlag(name,bf)`                                               | register a custom BaseFlag             |
| `Groups()`                                                            | list all mutual-exclusion groups       |

## License

Apache 2.0 — see [LICENSE](LICENSE)
