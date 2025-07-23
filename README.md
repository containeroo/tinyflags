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
- **Custom Placeholders** and **rich help** formatting
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

  // builder returns the flag – call .Value() to get corresponding value
  host := fs.String("host", "localhost", "server host").
              Required().
              Value()
  port := fs.Int("port", 8080, "server port").
              Short("p").
              Value()
  debug := fs.Bool("debug", false, "enable debug").
               Short("d").
               Value()
  tags  := fs.StringSlice("tag", nil, "tags list").
               Value()

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

````go
package main

import (
  "fmt"
  "os"
  "github.com/containeroo/tinyflags"
)

func main() {
  fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
  dg := fs.DynamicGroup("http")

  // define per-instance flags; builder returned until you call Get()/MustGet()
  portFlag    := dg.Int("port", "backend port")
  timeoutFlag := dg.Duration("timeout", "request timeout")

  if err := fs.Parse(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  // iterate over all IDs seen
  for _, id := range dg.Instances() {
    port, _    := portFlag.Get(id)
    timeout, _ := timeoutFlag.Get(id)
    fmt.Printf("%s → port=%d, timeout=%s\n", id, port, timeout)
  }
}

Call it like:

```bash
./app --http.alpha.port=8080 --http.alpha.timeout=30s \
      --http.beta.port=9090 --http.beta.timeout=1m
````

Outputs:

```text
alpha → port=8080, timeout=30s
beta → port=9090, timeout=1m
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

| Go Type               | Methods                                |
| --------------------- | -------------------------------------- |
| `bool`                | `Bool`, `BoolVar`                      |
| `counter`\*           | `Counter`, `CounterVar`                |
| `int`                 | `Int`, `IntVar`                        |
| `string`              | `String`, `StringVar`                  |
| `[]string`            | `StringSlice`, `StringSliceVar`        |
| `time.Duration`       | `Duration`, `DurationVar`              |
| `net.IP` / `[]net.IP` | `IP`, `IPVar`, `IPSlice`, `IPSliceVar` |
| `*net.TCPAddr`        | `TCPAddr`, `TCPAddrVar`                |
| `url.URL`             | `URL`, `URLVar`                        |
| `*os.File`            | `File`, `FileVar`                      |

\* _Counter_ flags are special: they increment on each occurrence, `-vv` → `2`, `-vvv` → `3`, etc.

> All slice flags also support repeated usage and custom delimiters.

## FlagSet Methods

| Method                                                                | Description                           |
| --------------------------------------------------------------------- | ------------------------------------- |
| `NewFlagSet(name, mode)`                                              | create a new named FlagSet            |
| `Parse(args)`                                                         | parse flags and env vars              |
| `Version(s)`                                                          | set version string for `--version`    |
| `EnvPrefix(s)`                                                        | set global env-var prefix             |
| `Authors(s)`                                                          | add author info to help               |
| `Title(s)`                                                            | set help title                        |
| `Description(s)`                                                      | add a description paragraph           |
| `Note(s)`                                                             | append a note paragraph               |
| `DisableHelp()`, `DisableVersion()`                                   | disable built-in help/version flags   |
| `Sorted(bool)`                                                        | enable/disable sorted flag output     |
| `SetOutput(w)`, `Output()`                                            | set/get output writer                 |
| `IgnoreInvalidEnv(bool)`                                              | skip invalid env-var values           |
| `SetGetEnvFn(fn)`                                                     | customize env-var lookup              |
| `Globaldelimiter(s)`                                                  | set default slice-delimiter           |
| `RequirePositional(n)`                                                | require at least n positional args    |
| `Args()`, `Arg(i)`                                                    | retrieve positional args              |
| `GetGroup(name)`                                                      | create/get a mutual-exclusion group   |
| `PrintDefaults()`                                                     | print all flags with defaults         |
| `PrintUsage(w, mode)`                                                 | print usage in specified mode         |
| `PrintTitle(w)`, `PrintNotes(w, width)`, `PrintDescription(w, width)` | various help sections                 |
| `PrintAuthors(w)`                                                     | print author info                     |
| `DynamicGroup(name)`                                                  | create/get a dynamic-flag group       |
| `DefaultDelimiter()`                                                  | get default delimiter for slice flags |

> **Note**
> Every _static_ flag method returns the builder so you can chain calls (e.g. `.Required()`, `.Choices()`, `.Short("x")`).
> Only when you call `.Value()` do you actually get back the `*T`.
>
> _Dynamic_ flags work differently: you define them on a `Group`, but after parsing you must fetch values with `Get(id)`, `MustGet(id)`, or `Values()`.

## License

Apache 2.0 — see [LICENSE](LICENSE)
