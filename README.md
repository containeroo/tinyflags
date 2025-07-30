# tinyflags

[![Go Reference](https://pkg.go.dev/badge/github.com/containeroo/tinyflags.svg)](https://pkg.go.dev/github.com/containeroo/tinyflags)

A minimal, fast, and extensible CLI flag-parsing library for Go.
Zero dependencies. Full generics support. Rich usage output.

## Features

- **Short & long flags** (`-d`, `--debug`)
- **Boolean strict mode** (`--flag=true/false`, `--no-flag`)
- **Environment variable overrides** (`EnvPrefix`, per-flag opt-out)
- **Required, deprecated, and grouped flags**
- **Slice flags** (`[]T`) with custom delimiters
- **Allowed choices & validation**
- **Mutual-exclusion groups**
- **Require-together groups**
- **Custom placeholders & help sections**
- **Dynamic flags** (`--group.id.field=value`)
- **Typed values** (`*os.File`, `*net.TCPAddr`, `url.URL`, `time.Duration`, etc.)

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
	fs.EnvPrefix("MYAPP")
	fs.Version("v1.2.3")

	host := fs.String("host", "localhost", "Server hostname").Required().Value()
	port := fs.Int("port", 8080, "Port to bind").Short("p").Value()
	debug := fs.Bool("debug", false, "Enable debug logging").Short("d").Value()
	tags  := fs.StringSlice("tag", nil, "Optional tags").Value()

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

You can disable env binding per-flag:

```go
fs.Bool("internal", false, "internal use only").DisableEnv()
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
	http := fs.DynamicGroup("http")

	port    := http.Int("port", "Backend port")
	timeout := http.Duration("timeout", "Request timeout")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, id := range http.Instances() {
		p, _ := port.Get(id)
		t, _ := timeout.Get(id)
		fmt.Printf("%s: port=%d, timeout=%s\n", id, p, t)
	}
}
```

```bash
./app --http.a.port=8080 --http.a.timeout=30s \
      --http.b.port=9090 --http.b.timeout=1m
```

```text
a: port=8080, timeout=30s
b: port=9090, timeout=1m
```

## Help Output

```text
Usage: app [flags]

Flags:
  --host HOST         Server hostname (Default: localhost) (Env: MYAPP_HOST) (Required)
  -p, --port PORT     Port to bind (Default: 8080) (Env: MYAPP_PORT)
  -d, --debug         Enable debug logging (Env: MYAPP_DEBUG)
      --tag TAG...    Optional tags
  -v, --version       Show version
```

## Supported Types

| Type            | Methods                                  |
| --------------- | ---------------------------------------- |
| `bool`          | `Bool`, `BoolVar`                        |
| `int`           | `Int`, `IntVar`                          |
| `string`        | `String`, `StringVar`                    |
| `[]string`      | `StringSlice`, `StringSliceVar`          |
| `counter`       | `Counter`, `CounterVar` (auto-increment) |
| `time.Duration` | `Duration`, `DurationVar`                |
| `net.IP`        | `IP`, `IPVar`                            |
| `[]net.IP`      | `IPSlice`, `IPSliceVar`                  |
| `*net.TCPAddr`  | `TCPAddr`, `TCPAddrVar`                  |
| `url.URL`       | `URL`, `URLVar`                          |
| `*os.File`      | `File`, `FileVar`                        |

> Slice flags accept repeated use or custom-delimited strings.

## FlagSet API

### Core Methods

| Method                   | Description            |
| ------------------------ | ---------------------- |
| `NewFlagSet(name, mode)` | Create new flag set    |
| `Parse(args)`            | Parse CLI args + env   |
| `Usage func()`           | Custom usage handler   |
| `Version(s)`             | Set `--version` string |
| `DisableHelp()`          | Disable `--help`       |
| `DisableVersion()`       | Disable `--version`    |

### Help & Output

| Method                                             | Description               |
| -------------------------------------------------- | ------------------------- |
| `SetOutput(io.Writer)`                             | Set help output writer    |
| `Output()`                                         | Get current output writer |
| `Title(s)`                                         | Set usage title           |
| `Authors(s)`                                       | Set author section        |
| `Description(s)`                                   | Set help description      |
| `Note(s)`                                          | Set help footer           |
| `PrintUsage(w, mode)`                              | Print `Usage:` line       |
| `PrintTitle(w)`                                    | Print title header        |
| `PrintAuthors(w)`                                  | Print authors             |
| `PrintDescription(w, indent, width)`               | Print description block   |
| `PrintNotes(w, indent, width)`                     | Print footer block        |
| `PrintStaticDefaults(w, indent, startCol, width)`  | Static flag help          |
| `PrintDynamicDefaults(w, indent, startCol, width)` | Dynamic flag help         |

### Help Formatting

| Method                        | Description                     |
| ----------------------------- | ------------------------------- |
| `SetDescIndent(n)`            | Indent description lines        |
| `SetDescWidth(n)`             | Max description width           |
| `SetNoteIndent(n)`            | Indent for footer               |
| `SetNoteWidth(n)`             | Width for footer                |
| `SetUsageIndent(n)`           | Left padding for all flags      |
| `SetUsageColumn(n)`           | Column where description starts |
| `SetUsageWidth(n)`            | Max flag help line width        |
| `StaticAutoUsageColumn(pad)`  | Auto-calculate static col       |
| `DynamicAutoUsageColumn(pad)` | Auto-calculate dynamic col      |

### Environment

| Method                   | Description                    |
| ------------------------ | ------------------------------ |
| `EnvPrefix(s)`           | Prefix for all env vars        |
| `IgnoreInvalidEnv(true)` | Skip unknown envs              |
| `SetGetEnvFn(fn)`        | Custom `os.Getenv`             |
| `Globaldelimiter(s)`     | Slice delimiter (default: ",") |
| `DefaultDelimiter()`     | Get current delimiter          |

### Positional

| Method                 | Description             |
| ---------------------- | ----------------------- |
| `RequirePositional(n)` | Require at least N args |
| `Args()`               | Get all args            |
| `Arg(i)`               | Get i-th arg            |

### Mutual-Exclusive group

| Method                            | Description                          |
| --------------------------------- | ------------------------------------ |
| `MutualGroups()`                  | List all mutual-exlusive groups      |
| `AddMutualGroup(name, group)`     | Add mutual-exlusive group            |
| `GetMutualGroup(name)`            | Get mutual-exlusive group            |
| `AttachToMutualGroup(flag, name)` | Assign flag to mutual-exlusive group |

### Require-Together group

| Method                                 | Description                           |
| -------------------------------------- | ------------------------------------- |
| `RequireTogetherGroups()`              | List all require-together groups      |
| `AddRequireTogetherGroup(name, group)` | Add require-together group            |
| `GetRequireTogetherGroup(name)`        | Get require-together group            |
| `AttachToRequireTogetherGroup`         | Assign flag to require-together group |

### Dynamic Flags

| Method               | Description                 |
| -------------------- | --------------------------- |
| `DynamicGroup(name)` | Define a dynamic flag group |
| `DynamicGroups()`    | Get all registered groups   |

Also exposed:

```go
GetDynamic[T any](group, id, flag string) (T, error)
MustGetDynamic[T any](group, id, flag string) T
GetOrDefaultDynamic[T any](group, id, flag string) T
```

## License

Apache 2.0 — see [LICENSE](LICENSE)
