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
- **Allowed choices, validation and finalizers**
- **One Of groups**
- **All or None groups**
- **Custom placeholders & help sections**
- **Dynamic flags** (`--group.id.field=value`)
- **Typed values** (`*os.File`, `*net.TCPAddr`, `url.URL`, `time.Duration`, etc.)

**Why yet another flag library?**

- **Validate & Finalize on the fly**
  I got tired of the two-step tango--parse first, then wade through a swamp of `if`-statements just to check and tweak values. Tinyflags lets you **validate** and **finalize** your flags as they're parsed, so you can slap on your business logic (and data massaging) in one go.

- **Group therapy for flags**
  Ever tried juggling "onf of" or "all-or-nothing" flags with plain `flag`? It's like herding cats. Tinyflags brings built-in **onf-of** and **all-or-none** groups so your flags behave like well-trained puppies.

- **Self-service help & version**
  Want to bail out with `--help` or `--version` at just the right moment, without writing extra `if`-blocks? Tinyflags handles the exit routine for you, so you can spend less time plumbing and more time coding.

- **Dynamic flags--finally!**
  I looked high and low for a Go library that lets you declare `--group.id.field=value` flags, where `id` is dynamic. No luck. So I built one, folded it into tinyflags, and voila: one library to rule both "regular" **and** "shape-shifting" dynamic flags.

In short, tinyflags slices away boilerplate, stitches in the goodies I actually needed, and keeps my codebase lean--no extra flag-parsing baggage required. 🚀

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

## Supported Types

| Type            | Methods                                  |
| :-------------- | :--------------------------------------- |
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

### Common Flag-Builder Methods

| Method                      | Applies to  | Description                                                                             |
| :-------------------------- | :---------- | :-------------------------------------------------------------------------------------- |
| `Short(s string)`           | static only | One-letter alias (`-p`). Must be exactly one rune (panics otherwise).                   |
| `Required()`                | all flags   | Mark the flag as required; parser errors if unset.                                      |
| `HideRequired()`            | all flags   | Hide the “(Required)” suffix from help.                                                 |
| `Hidden()`                  | all flags   | Omit this flag from generated help output.                                              |
| `Deprecated(reason string)` | all flags   | Mark flag deprecated; includes `DEPRECATED` note in help.                               |
| `OneOfGroup(group string)`  | all flags   | Assign to a named mutual-exclusion group. Parsing errors if more than one in group set. |
| `AllOrNone(group string)`   | all flags   | Assign to a named require-together group. All or none in group must be set.             |
| `Env(key string)`           | all flags   | Override the environment-variable name (panics if `DisableEnv` already called).         |
| `HideEnv()`                 | all flags   | Hide the environment-variable name from help.                                           |
| `DisableEnv()`              | all flags   | Disable environment lookup for this flag (panics if `Env(...)` already called).         |
| `Placeholder(text string)`  | all flags   | Customize the `<VALUE>` placeholder in help.                                            |
| `Allowed(vals ...string)`   | all flags   | Restrict help to show only these allowed values.                                        |
| `HideAllowed()`             | all flags   | Hide the allowed values from help.                                                      |
| `Value() *T`                | static only | Return the pointer to the parsed value (after `Parse`).                                 |

### Static-Flag Extras

| Method                         | Description                                                                                  | Example                                                                                                                           |
| :----------------------------- | :------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------- |
| `Choices(v1, v2, ...)`         | Only allow the provided literal values; automatically adds them to help output.              | `fs.String("env","dev","...").Choices("dev","staging","prod")`                                                                    |
| `Validate(fn func(v T) error)` | Run custom check on parsed value; if `fn` returns non-nil, `Parse` returns an error.         | `go<br>fs.Int("count",0,"...").Validate(func(n int) error {<br>  if n<0 {return fmt.Errorf("must ≥0")}<br>  return nil<br>})<br>` |
| `Finalize(fn func(v T) T)`     | Transform the parsed value before storing; e.g. trimming, normalization, applying defaults.  | `go<br>fs.String("name","","...").Finalize(func(s string) string {<br>  return strings.TrimSpace(s)<br>})<br>`                    |
| `Delimiter(sep string)`        | _(slice flags only)_ Use a custom separator instead of the default comma when parsing lists. | `fs.StringSlice("tags",nil,"...").Delimiter(";")`                                                                                 |

### Dynamic-Flag Extras

| Method                               | Description                                                                       | Example                                                      |
| :----------------------------------- | :-------------------------------------------------------------------------------- | :----------------------------------------------------------- |
| `Has(id string) bool`                | Return whether the given instance-ID was provided on the command line or via env. | `if port.Has("a") { fmt.Println("port a is set") }`          |
| `Get(id string) (value T, ok bool)`  | Retrieve the parsed value for that ID; `ok==false` if unset (returns default).    | `p, ok := port.Get("a"); if ok { fmt.Println("a →", p) }`    |
| `MustGet(id string) T`               | Like `Get`, but panics if the instance wasn't provided.                           | `timeout.MustGet("b")`                                       |
| `Values() map[string]T`              | Get all parsed values keyed by instance ID.                                       | `for id, v := range timeout.Values() { fmt.Println(id, v) }` |
| `ValuesAny() map[string]interface{}` | Same as `Values()`, but values are `interface{}`.                                 |                                                              |

> All **common** methods (Required, Hidden, etc.) and **static extras** (Choices, Validate, Finalize, Delimiter) also apply to dynamic flags.

### FlagSet Core & Help Configuration

| Method                                            | Description                                                    |
| :------------------------------------------------ | :------------------------------------------------------------- |
| `NewFlagSet(name string, mode ErrorHandling)`     | Create a new flag set (e.g. `ExitOnError`, `ContinueOnError`). |
| `EnvPrefix(prefix string)`                        | Prefix all environment-variable lookups (e.g. `MYAPP_`).       |
| `Version(version string)`                         | Enable the `--version` flag, printing this string.             |
| `DisableHelp()` / `DisableVersion()`              | Remove `--help` or `--version`.                                |
| `Usage(fn func())`                                | Install a custom usage function in place of the default.       |
| `Title(text string)`                              | Override the "Usage:" title heading.                           |
| `Authors(names ...string)`                        | Add an "Authors:" section to help.                             |
| `Description(text string)`                        | Add a free-form description block under the title.             |
| `Note(text string)`                               | Add a footer note under the flags listing.                     |
| `SetOutput(w io.Writer)` / `Output()`             | Redirect or retrieve where help/version is written.            |
| **Help Printers:**                                |                                                                |
|   `PrintUsage(w, mode)`                           | Print the `Usage:` line.                                       |
|   `PrintTitle(w)`                                 | Print title and description.                                   |
|   `PrintAuthors(w)`                               | Print authors section.                                         |
|   `PrintDescription(w,indent,width)`              | Print the description block.                                   |
|   `PrintNotes(w,indent,width)`                    | Print footer notes.                                            |
|   `PrintStaticDefaults(w,indent,startCol,width)`  | Print static flags help.                                       |
|   `PrintDynamicDefaults(w,indent,startCol,width)` | Print dynamic flags help.                                      |
| `RequirePositional(n int)`                        | Enforce at least `n` positional arguments.                     |
| `Args() []string` / `Arg(i int) string`           | Access leftover positional args.                               |
| `AddOneOfGroup(name string, flags []string)`      | Manually define a mutual-exclusion group.                      |
| `AddAllOrNoneGroup(name string, flags []string)`  | Manually define a require-together group.                      |

### How `Validate` and `Finalize` Work

1. **Validate**

   - After parsing a flag's raw input into `T`, Tinyflags calls your validator:

   ```go
   fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)

   count := fs.Int("count", 0, "Number of items").
       Validate(func(n int) error {
           if n < 0 {
               return fmt.Errorf("must be non-negative")
           }
           return nil
       }).
       Value()

   fs.Parse(os.Args[1:])
   fmt.Println("Count:", *count)
   ```

   - If the user does `--count=-5`, they see:

     ```text
     invalid value for --count: must be non-negative
     ```

   - On error, parsing aborts and your message is shown to the user.

2. **Finalize**

   - Only after validation succeeds, Tinyflags passes the parsed value through your finalizer:

   ```go
   fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)

   name := fs.String("name", "", "User name").
     Finalize(func(s string) string {
             return strings.TrimSpace(strings.ToTitle(s))
             }).
     Value()

   url := fs.String("url", "", "URL to use").
     // Ensure the URL ends with a slash
     Finalize(func(u *url.URL) *url.URL {
        // Clone to avoid mutating the original (optional, if needed)
        u2 := *u
        if len(u2.Path) > 0 && u2.Path[len(u2.Path)-1] != '/' {
            u2.Path += "/"
        }
        return &u2
     }).
     Value()


    fs.Parse([]string{"--name=   alice smith  ", "--url", "https://containeroo.ch"})
    fmt.Println("Hello,", *name)
    fmt.Println("Visist:", *url)
   ```

   _Output:_

   ```text
   Hello, Alice Smith
   Visti: https://containeroo.ch/
   ```

   - Useful for trimming whitespace, applying normalization, setting derived defaults, etc.

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

## Grouped Flags: Mutual-Exclusion & Require-Together

When certain flags must be used **together**, or must be **exclusive**, tinyflags makes that easy.

### 🔁 Require-Together Group

You can define a group where **either all or none** of the flags must be set:

```go
email := fs.String("email", "", "User email").
    AllOrNone("authpair").
    Value()

password := fs.String("password", "", "User password").
    AllOrNone("authpair").
    Value()
```

### 🔀 Mutual-Exclusion Group

You can define a group where **only one** of the flags (or groups!) may be set:

```go
bearer := fs.String("bearer-token", "", "Bearer token").
    OneOfGroup("authmethod").
    Value()

fs.GetOneOfGroup("authmethod").
    Title("Authentication Method").
    AddGroup(fs.GetAllOrNoneGroup("authpair")) // <-- email+password
```

This enforces:

- Either `--email` and `--password` **together**, or
- `--bearer-token`, but **not both**

### ✅ Valid combinations

```text
# Valid:
--email=user --password=secret
--bearer-token=abc123

# Invalid:
--email=user                ❌ password missing
--email=user --bearer-token=abc123 ❌ both auth methods
```

### Help Output Example

```text
Usage: app [flags]

Flags:
  --email EMAIL         User email (Group: Authentication Method) (AllOrNone: authpair)
  --password PASSWORD   User password (Group: Authentication Method) (AllOrNone: authpair)
  --bearer-token TOKEN  Bearer token (Group: Authentication Method)

Authentication Method
(Exactly one required)
  --bearer-token TOKEN
  [--email, --password] (Required together)

authpair
(Required)
  --email EMAIL
  --password PASSWORD
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

## License

Apache 2.0 -- see [LICENSE](LICENSE)
