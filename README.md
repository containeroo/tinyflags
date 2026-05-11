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

- **Validate and finalize on the fly**
  Tinyflags lets you **validate** and **finalize** values as they are parsed, so input shaping and business rules can live directly on the flag definition.

- **Group-aware constraints**
  Tinyflags includes built-in **one-of** and **all-or-none** groups, plus per-flag dependency rules like `Requires(...)`.

- **Built-in help and version handling**
  `--help` and `--version` are supported out of the box, including library-friendly error sentinels for callers that want to intercept them.

- **Dynamic flags**
  Tinyflags supports `--group.id.field=value` flags where `id` is only known at parse time, so one parser can handle both regular and instance-scoped settings.

In short, tinyflags aims to reduce flag boilerplate while keeping the parser predictable, typed, and easy to extend.

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

## Commands

`Command` lets you build subcommand trees with local flags and inherited globals. If you want parsing to fail unless a subcommand is chosen, call `RequireCommand()`.

```go
app := tinyflags.NewCommand("app", tinyflags.ExitOnError).RequireCommand()
app.Command("serve", "Run the server")

if err := app.Parse(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
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

## Parse Model

Tinyflags applies input in this order:

1. Parse command-line arguments.
2. Handle built-in `--help` / `--version` exits.
3. Load unset static flags from environment variables.
4. Apply default finalizers for unset values.
5. Run required/group/dependency/positional validation.

Additional behavior:

- Explicit CLI arguments win over environment variables.
- `OverriddenValues()` reports values provided by CLI or env, not untouched defaults.
- Reusing a `FlagSet` across multiple `Parse(...)` calls is supported; parser state is reset before each parse.
- Dynamic flags are CLI-only and are not populated from environment variables.

### Handling toggles with multiple flags

You can model toggles with paired flags (e.g., `--debug` and `--no-debug`) and pick the first one the user set:

```go
debug := fs.Bool("debug", false, "Enable debug").Short("d").OneOfGroup("debug")
noDebug := fs.Bool("no-debug", false, "Disable debug").Short("n").OneOfGroup("debug")

// Use the order you prefer; first changed wins
enabled, set := tinyflags.FirstChanged(false, debug, noDebug)
fmt.Printf("debug: %t (set: %v)\n", enabled, set)
```

### Helpers and exported utilities

- `FirstChanged[T](defaultValue, flags...)` — returns the value of the first changed flag (by order) plus whether any flag was set.
- `IsHelpRequested(err)` / `IsVersionRequested(err)` — detect help/version parse exits.
- `RequestHelp(msg)` / `RequestVersion(msg)` — trigger help/version errors manually.
- `Flag[T]` — minimal interface implemented by flag handles (`Changed() bool`, `Value() *T`).

## FlagSet API

### Common Flag-Builder Methods

| Method                      | Applies to  | Description                                                                             |
| :-------------------------- | :---------- | :-------------------------------------------------------------------------------------- |
| `Short(s string)`           | static only | One-letter alias (`-p`). Must be exactly one rune (panics otherwise).                   |
| `Required()`                | all flags   | Mark the flag as required; parser errors if unset.                                      |
| `HideRequired()`            | all flags   | Hide the "(Required)" suffix from help.                                                 |
| `Hidden()`                  | all flags   | Omit this flag from generated help output.                                              |
| `Deprecated(reason string)` | all flags   | Mark flag deprecated; includes `DEPRECATED` note in help.                               |
| `OneOfGroup(group string)`  | all flags   | Assign to a named mutual-exclusion group. Parsing errors if more than one in group set. |
| `HelpOneOfGroups(names...)` | all flags   | Override which one-of groups for this flag are shown in help output.                    |
| `AllOrNone(group string)`   | all flags   | Assign to a named require-together group. All or none in group must be set.             |
| `Env(key string)`           | all flags   | Override the environment-variable name (panics if `DisableEnv` already called).         |
| `HideEnv()`                 | all flags   | Hide the environment-variable name from help output.                                     |
| `DisableEnv()`              | all flags   | Disable environment lookup for this flag (panics if `Env(...)` already called).         |
| `Placeholder(text string)`  | all flags   | Customize the `<VALUE>` placeholder in help.                                            |
| `Allowed(vals ...string)`   | all flags   | Restrict help to show only these allowed values.                                        |
| `HideAllowed()`             | all flags   | Hide the allowed values from help.                                                      |
| `Requires(names ...string)` | all flags   | Mark flag as required by the given flag.                                                |
| `HideRequires()`            | all flags   | Hide the “(Requires)” suffix from help.                                                 |
| `OverriddenValueMaskFn(fn)` | all flags   | Provide a mask function used by `OverriddenValues()`.                                   |
| `Value() *T`                | static only | Return the pointer to the parsed value (after `Parse`).                                 |

### Static-Flag Extras

| Method                                      | Description                                                                                  | Example                                                                                                                           |
| :------------------------------------------ | :------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------- |
| `Choices(v1, v2, ...)`                      | Only allow the provided literal values; automatically adds them to help output.              | `fs.String("env","dev","...").Choices("dev","staging","prod")`                                                                    |
| `Validate(fn func(v T) error)`              | Run custom check on parsed value; if `fn` returns non-nil, `Parse` returns an error.         | `go<br>fs.Int("count",0,"...").Validate(func(n int) error {<br>  if n<0 {return fmt.Errorf("must ≥0")}<br>  return nil<br>})<br>` |
| `Finalize(fn func(v T) T)`                  | Transform the parsed value before storing; e.g. trimming, normalization, applying defaults.  | `go<br>fs.String("name","","...").Finalize(func(s string) string {<br>  return strings.TrimSpace(s)<br>})<br>`                    |
| `FinalizeDefaultValue()`                    | Run the existing finalizer on default values when the flag is unset.                          | `go<br>fs.String("name","","...").Finalize(strings.TrimSpace).FinalizeDefaultValue()<br>`                                        |
| `FinalizeWithID(fn func(id string, v T) T)` | _(dynamic only)_ Finalize with access to the instance ID.                                    | `http.String("addr","","").FinalizeWithID(func(id, v string) string { return id+":"+v })`                                         |
| `Delimiter(sep string)`                     | _(slice flags only)_ Use a custom separator instead of the default comma when parsing lists. | `fs.StringSlice("tags",nil,"...").Delimiter(";")`                                                                                 |
| `AllowEmpty()`                              | _(slice flags only)_ Allow empty items (e.g. `"a,,b"`).                                      |                                                                                                                                   |
| `HideDefault()`                             | Hide the default value from help output.                                                     |                                                                                                                                   |
| `Section(name string)`                      | Group flags under a section header in help output.                                           |                                                                                                                                   |

### Dynamic-Flag Extras

### Multiple one-of groups

You can attach a flag to more than one one-of group by chaining `OneOfGroup(...)`. By default, all visible memberships are rendered in help. If you want help to mention only a subset, add `HelpOneOfGroups(...)`.

```go
searchHistoryOnly := fs.Bool("history-only", false, "Search shell history commands only.").
    OneOfGroup("search-mode").
    Value()

searchPinsOnly := fs.Bool("pins-only", false, "Search pinned commands only.").
    OneOfGroup("search-mode").
    OneOfGroup("pins-scope").
    HelpOneOfGroups("search-mode").
    Value()

searchExcludePins := fs.Bool("exclude-pins", false, "Exclude pinned commands from history results.").
    OneOfGroup("pins-scope").
    Value()
```

| Method                               | Description                                                                       | Example                                                      |
| :----------------------------------- | :-------------------------------------------------------------------------------- | :----------------------------------------------------------- |
| `Has(id string) bool`                | Return whether the given instance-ID was provided on the command line or via env. | `if port.Has("a") { fmt.Println("port a is set") }`          |
| `Get(id string) (value T, ok bool)`  | Retrieve the parsed value for that ID; `ok==false` if unset (returns default).    | `p, ok := port.Get("a"); if ok { fmt.Println("a →", p) }`    |
| `MustGet(id string) T`               | Like `Get`, but panics if the instance wasn't provided.                           | `timeout.MustGet("b")`                                       |
| `Values() map[string]T`              | Get all parsed values keyed by instance ID.                                       | `for id, v := range timeout.Values() { fmt.Println(id, v) }` |
| `ValuesAny() map[string]interface{}` | Same as `Values()`, but values are `interface{}`.                                 |                                                              |
| `AllowOverride()`                    | Allow re-assignment of a dynamic flag per-ID. Only for Scalar flags.              |                                                              |

> All **common** methods (Required, Hidden, etc.) and **static extras** (Choices, Validate, Finalize, FinalizeDefaultValue, Delimiter) also apply to dynamic flags.

### FlagSet Core & Help Configuration

| Method                                             | Description                                                          |
| :------------------------------------------------- | :------------------------------------------------------------------- |
| `NewFlagSet(name string, mode ErrorHandling)`      | Create a new flag set (e.g. `ExitOnError`, `ContinueOnError`).       |
| `EnvPrefix(prefix string)`                         | Prefix all environment-variable lookups (e.g. `MYAPP_`).             |
| `SetEnvKeyFunc`                                    | Set a function to derive env keys from prefix+flag name.             |
| `EnvKeyForFlag`                                    | Derive the env key for a flag.                                       |
| `NewReplacerEnvKeyFunc`                            | Build an `EnvKeyFunc` that applies the given replacer.               |
| `Version(version string)`                          | Enable the `--version` flag, printing this string.                   |
| `Help()`                                           | Access grouped helpers for title/authors/description/note/help text. |
| `Layout()`                                         | Access grouped helpers for usage/indent/width/note layout.           |
| `BeforeParse(fn func([]string) ([]string, error))` | Mutate arguments before parsing (e.g., expand @files).               |
| `OnUnknownFlag(fn func(name string) error)`        | Handle or ignore unknown flags instead of failing.                   |
| `VersionText(text string)`                         | Override the `--version` text. Default: `"Show version"`.            |
| `HelpText(text string)`                            | Override the `--help` text. Default: `"Show help"`.                  |
| `DisableHelp()` / `DisableVersion()`               | Remove `--help` or `--version`.                                      |
| `Usage func()`                                     | Optional custom usage function on `FlagSet` that replaces the default renderer. |
| `Title(text string)`                               | Override the "Usage:" title heading.                                 |
| `Authors(text string)`                             | Add an `Authors:` section to help output.                            |
| `Description(text string)`                         | Add a free-form description block under the title.                   |
| `Note(text string)`                                | Add a footer note under the flags listing.                           |
| `SetOneOfGroupVerbose(bool)`                       | Toggle detailed OneOfGroup errors with conflicting flags.            |
| `SetOutput(w io.Writer)` / `Output()`              | Redirect or retrieve where help/version is written.                  |
| `PrintUsage(w, mode)`                              | Print the `Usage:` line.                                             |
| `PrintTitle(w)`                                    | Print title and description.                                         |
| `PrintAuthors(w)`                                  | Print authors section.                                               |
| `PrintDescription(w,indent,width)`                 | Print the description block.                                         |
| `PrintNotes(w,indent,width)`                       | Print footer notes.                                                  |
| `PrintStaticDefaults(w,indent,startCol,width)`     | Print static flags help.                                             |
| `PrintDynamicDefaults(w,indent,startCol,width)`    | Print dynamic flags help.                                            |
| `RequirePositional(n int)`                         | Enforce at least `n` positional arguments.                           |
| `Args() []string` / `Arg(i int) (string, bool)`    | Access leftover positional args safely.                              |
| `OverriddenValues() map[string]any`                | Return flags explicitly set via args/env (dynamic keys: `group.id.flag`). |
| `MaskFirstLast(value any) any`                     | Helper mask that keeps first/last character (strings, `[]string`).   |
| `MaskPostgresURL(value any) any`                   | Helper mask for `postgres://user:pass@host/db` credentials.          |
| `AddOneOfGroup(name string, group *core.OneOfGroupGroup)` | Register a pre-built mutual-exclusion group.                 |
| `AddAllOrNoneGroup(name string, group *core.AllOrNoneGroup)` | Register a pre-built require-together group.           |

### Command API

| Method | Description |
| :-- | :-- |
| `NewCommand(name string, mode ErrorHandling)` | Create a root command tree. |
| `Command(name, summary string)` | Register a child command. |
| `Globals()` | Access persistent flags inherited by that subtree. |
| `RequireCommand()` | Return an error if this command is selected without a child command. |
| `Parse(args)` | Parse flags and select the active command. |
| `SelectedCommand()` | Return the selected leaf command from the last parse. |
| `Run(handler, bindings...)` / `BuildCommand(builder)` | Register execution for a command. |
| `ParseRunner(args)` / `ParseRunnable(args)` | Parse and build the selected runnable. |

### Naming notes

- `GlobalDelimiter(...)` is the preferred spelling.
- `Globaldelimiter(...)` remains available as a compatibility alias.
- `AllOrNoneGroups()` is the preferred plural accessor.
- `AllOrNoneGroup()` remains available as a compatibility alias.

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

3. **FinalizeDefaultValue**
   - Runs only when a flag is *unset* (after env parsing), and does not mark it as changed.
   - Uses the flag’s `Finalize(...)` function; no separate function is provided.

   ```text
   Hello, Alice Smith
   Visti: https://containeroo.ch/
   ```

   - Useful for trimming whitespace, applying normalization, setting derived defaults, etc.

## Environment Variables

Flags can be set from environment variables in addition to CLI arguments.
By default, environment keys are derived from:

- the global prefix set via `fs.EnvPrefix("MYAPP")`
- `-`, `_`, `.`, `_`, `/` are replaced with `_` in the flag name
- the whole key upper-cased

For example:

```bash
MYAPP_HOST=example.com MYAPP_PORT=9090 ./app --debug
```

A flag `--db.user` will look up `MYAPP_DB_USER`.

You can disable env binding per-flag:

```go
fs.Bool("internal", false, "internal use only").DisableEnv()
```

### Custom key mapping

You can override how keys are derived with `SetEnvKeyFunc`:

```go
fs.EnvPrefix("MYAPP")
fs.SetEnvKeyFunc(engine.NewReplacerEnvKeyFunc(
    strings.NewReplacer("-", "_", ".", "_", "/", "_"),
    true, // upper-case
))
```

This would map:

- `--log.level` → `MYAPP_LOG_LEVEL`
- `--db-user` → `MYAPP_DB_USER`
- `--api/v1` → `MYAPP_API_V1`

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
