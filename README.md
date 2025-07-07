# tinyflags

[![Go Reference](https://pkg.go.dev/badge/github.com/containeroo/tinyflags.svg)](https://pkg.go.dev/github.com/containeroo/tinyflags)

`tinyflags` is a minimal, fast, and extensible CLI flag parsing library for Go, heavily inspired by [`pflag`](https://github.com/spf13/pflag) and [`kingpin`](https://github.com/alecthomas/kingpin). It supports advanced features like environment overrides, validation, required flags, shorthand syntax, and flexible usage formatting -- all while staying light and dependency-free.

## Features

- 🧩 Short and long flags (`-d`, `--debug`)
- 🧪 Boolean flags with optional `--no-flag` strict parsing (`--no-flag=true`)
- 🌲 Environment variable support (per-flag and global prefix)
- 🧯 Required flags and deprecation notices
- 🧵 Support for `[]T` slice flags with custom delimiters
- ⚠️ Choices and validation
- 🧑‍🤝‍🧑 Mutual exclusion groups
- 🔒 Disable environment lookups per flag
- 🎛 Custom metavars for help
- 🆘 Rich help output with wrapping and alignment
- ✅ Direct value assignment via pointer
- 🛠 Fully type-safe API with generics

## Install

```bash
go get github.com/containeroo/tinyflags
```

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

func main() {
	// Usually you would parse args from os.Args[1:]
	// but for this example we'll just hard-code them.
	args := []string{
		"--port=9000",
		"--host=example.com",
		"-vtrue",
	}

	fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)
	fs.Version("v1.0")    // optional, enables -v, --version
	fs.EnvPrefix("MYAPP") // optional, enables --env-key for all flags

	host := fs.String("host", "localhost", "host to use").
		Required().
		Value()

	port := fs.Int("port", 8080, "port to listen on").
		Env("MYAPP_CUSTOM_PORT"). // overrides default env key (otherwise "MYAPP_PORT")
		Required().
		Value()

	debug := fs.BoolP("debug", "d", false, "enable debug mode").
		Strict().
		Value()

	tags := fs.StringSlice("tag", []string{}, "list of tags").
		Value()

	loglevel := fs.String("log-level", "info", "log level to use").
		Choices("debug", "info", "warn", "error").
		Value()

	if err := fs.Parse(args); err != nil {
		if tinyflags.IsHelpRequested(err) || tinyflags.IsVersionRequested(err) {
			fmt.Fprint(os.Stdout, err.Error()+"\n") // nolint:errcheck
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err.Error()+"\n") //nolint:errcheck
		os.Exit(1)
	}

	fmt.Println("Host:", *host)
	fmt.Println("Port:", *port)
	fmt.Println("Debug:", *debug)
	fmt.Println("Tags:", *tags)
	fmt.Println("Loglevel:", *loglevel)
}
```

Help output looks like this:

```text
Usage: test.exe -d -i -v <true|false>
      --port PORT                             port to use (Default: 8080) (Env: MYAPP_CUSTOM_PORT) (Required)
      --host HOST                             host to use (Default: localhost) (Env: MYAPP_HOST) (Required)
      --host-ip HOST-IP                       host ip to use. Must be in range 10.0.10.0/32
                                              (Default: 10.0.10.8) (Env: MYAPP_HOST_IP)
      --log-level <debug|info|warn|error>     log level to use (Allowed: debug, info, warn, error)
                                              (Default: info) (Env: MYAPP_LOG_LEVEL)
  -d, --debug                                 debug mode (Env: MYAPP_DEBUG)
  -i, --insecure                              insecure mode (Env: MYAPP_INSECURE)
  -v, --verbose <true|false>                  verbose mode (Default: false) (Env: MYAPP_VERBOSE)
      --version                               show version

```

## Environment Variables

You can override flags via environment variables:

```bash
MYAPP_PORT=9090 ./myapp --debug
```

To control this behavior globally:

```go
fs.EnvPrefix("MYAPP")
```

To disable env lookup for a specific flag:

```go
fs.Bool("internal-flag", false, "something internal").
   DisableEnv()
```

## Help Output

```text
Usage: myapp [flags] // tinyflags.PrintUsage()

My Title (tinyFlags.Title("My Title")

Some Description (tinyFlags.Description("Some Description")

      --port PORT       port to listen on (Default: 8080) (Env: MYAPP_PORT) (Required)
  -d, --debug           enable debug mode (Env: MYAPP_DEBUG)
      --tag TAG...      list of tags (Env: MYAPP_TAG)

Note (tinyFlags.Note("Text")
```

Supports aligned columns, wrapped descriptions, and positional arguments.

## 📦 Flag Types

| Go Type           | Flag Methods (shorthand versions use `*P` suffix)                                  | Notes                                         |
| ----------------- | ---------------------------------------------------------------------------------- | --------------------------------------------- |
| `bool`            | `Bool`, `BoolP`, `BoolVar`, `BoolVarP`                                             |                                               |
| `int`             | `Int`, `IntP`, `IntVar`, `IntVarP`                                                 |                                               |
| `string`          | `String`, `StringP`, `StringVar`, `StringVarP`                                     |                                               |
| `float64`         | `Float64`, `Float64P`, `Float64Var`, `Float64VarP`                                 |                                               |
| `time.Duration`   | `Duration`, `DurationP`, `DurationVar`, `DurationVarP`                             |                                               |
| `[]string`        | `Strings`, `StringsP`, `StringsVar`, `StringsVarP`                                 | Also supports delimiter                       |
| `[]int`           | `Ints`, `IntsP`, `IntsVar`, `IntsVarP`                                             |                                               |
| `[]float64`       | `Float64s`, `Float64sP`, `Float64sVar`, `Float64sVarP`                             |                                               |
| `[]time.Duration` | `Durations`, `DurationsP`, `DurationsVar`, `DurationsVarP`                         |                                               |
| `net.IP`          | `IP`, `IPP`, `IPVar`, `IPVarP`                                                     |                                               |
| `[]net.IP`        | `IPSlice`, `IPSliceP`, `IPSliceVar`, `IPSliceVarP`                                 |                                               |
| `net.IPMask`      | `IPMask`, `IPMaskP`, `IPMaskVar`, `IPMaskVarP`                                     | IPv4 mask only                                |
| `[]net.IPMask`    | `IPMaskSlice`, `IPMaskSliceP`, `IPMaskSliceVar`, `IPMaskSliceVarP`                 |                                               |
| `*net.TCPAddr`    | `ListenAddr`, `ListenAddrP`, `ListenAddrVar`, `ListenAddrVarP`                     | Parsed using `net.ResolveTCPAddr("tcp", ...)` |
| `[]*net.TCPAddr`  | `ListenAddrSlice`, `ListenAddrSliceP`, `ListenAddrSliceVar`, `ListenAddrSliceVarP` |                                               |
| `*os.File`        | `File`, `FileP`, `FileVar`, `FileVarP`                                             | Opened for reading (`os.Open`)                |
| `[]*os.File`      | `FileSlice`, `FileSliceP`, `FileSliceVar`, `FileSliceVarP`                         | All files opened                              |
| `url.URL`         | `URL`, `URLP`, `URLVar`, `URLVarP`                                                 | Parsed using `url.Parse`                      |
| `[]url.URL`       | `URLSlice`, `URLSliceP`, `URLSliceVar`, `URLSliceVarP`                             |                                               |
| `string` (custom) | `SchemaHostPort`, `SchemaHostPortP`, `SchemaHostPortVar`, `SchemaHostPortVarP`     | Requires format `scheme://host:port`          |

> 💡 All slice types support configurable delimiters (e.g., `--ips=1.1.1.1,8.8.8.8`) and repeated usage (e.g., `--ips=1.1.1.1 --ips=8.8.8.8`).

All flags return a flag with methods:

### Builder Methods

| Method                      | Description                                                                     |
| --------------------------- | ------------------------------------------------------------------------------- |
| `.Env("KEY")`               | Set a custom environment variable                                               |
| `.Required()`               | Mark the flag as required                                                       |
| `.Deprecated("...")`        | Mark the flag as deprecated                                                     |
| `.Group("name")`            | Place flag in a mutual exclusion group                                          |
| `.Choices(...)`             | Allow only specific values (auto-formatted in help)                             |
| `.Validator(func)`          | Add a custom validation function                                                |
| `.Metavar("FOO")`           | Customize the metavar name in help                                              |
| `.Strict()` (`bool`)        | For bool flags, require explicit `--flag=true/false`, also supports `--no-flag` |
| `.Delimiter(",")` (`slice`) | For slice flags, split using delimiter                                          |
| `.DisableEnv()`             | Prevent the flag from being set via environment variable                        |
| `.Value()`                  | Get a pointer to the parsed flag value                                          |

## Mutual Exclusion Groups

```go
fs.String("file", "", "use a file").Group("input")
fs.String("url", "", "fetch from url").Group("input")
```

Only one of `--file` or `--url` can be set.

## Positional Arguments

```go
fs.RequirePositional(1) // Require at least 1 positional
```

Use `fs.Args()` to access them.

## Advanced Output Customization

```go
fs.DescriptionIndent(50)
fs.DescriptionMaxLen(100)
fs.UsagePrintMode(tinyflags.PrintShort)
fs.PrintUsage()
fs.PrintTitle("My awsome App")
fs.PrintDescription("MyApp is a simple tool...")
fs.PrintDefaults() // show all flags with their current values
```

## Behavior

- Unrecognized flags return an error.
- Positional arguments begin after `--` or after the first non-flag arg.
- Bools default to `false`, or `true` if present (unless `Strict()` is used).
- Slice flags can be repeated like `--flag=a --flag=b --flag=c` or `--flag=a,b,c`
- Short flags can be combined like `-abc` or `-a -b -c`.
- Short flags can be like `-p 8080` or `-p=8080`.

## Nice to know

Use `fs.Args()` to access all positional arguments.
Use `IsHelpRequested(err)` and `IsVersionRequested(err)` to check for help/version flags.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.
