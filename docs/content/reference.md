# Help and API reference

## Go API

For the complete exported API and signatures, see [pkg.go.dev](https://pkg.go.dev/github.com/containeroo/tinyflags).

## Configuring help

```go
fs.Title("Example application")
fs.Description("Manage a service.")
fs.SectionOrder("General", "Network")
fs.Int("port", 8080, "Listen port").Section("Network")
fs.Help().BuiltinSection("General")
```

Use `HideEnvs()` to omit environment keys from all help, `HideEnv()` on individual flags, and `Layout()` for indentation and column settings. These settings affect help presentation, not parsing.

## Help output

```text
Usage: app [flags]

Flags:
  --host HOST         Server hostname (Default: localhost) (Env: MYAPP_HOST) (Required)
  -p, --port PORT     Port to bind (Default: 8080) (Env: MYAPP_PORT)
  -d, --debug         Enable debug logging (Env: MYAPP_DEBUG)
      --tag TAG...    Optional tags
  -v, --version       Show version
```

## FlagSet core and help configuration

| Method                                                       | Description                                                                                                      |
| :----------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------- |
| `NewFlagSet(name string, mode ErrorHandling)`                | Create a new flag set (e.g. `ExitOnError`, `ContinueOnError`).                                                   |
| `EnvPrefix(prefix string)`                                   | Prefix all environment-variable lookups (e.g. `MYAPP_`).                                                         |
| `SetEnvKeyFunc`                                              | Set a function to derive env keys from prefix+flag name.                                                         |
| `EnvKeyForFlag`                                              | Derive the env key for a flag.                                                                                   |
| `NewReplacerEnvKeyFunc`                                      | Build an `EnvKeyFunc` that applies the given replacer.                                                           |
| `Version(version string)`                                    | Enable the `--version` flag, printing this string.                                                               |
| `Help()`                                                     | Access grouped helpers for title/authors/description/note/help text.                                             |
| `Layout()`                                                   | Access grouped helpers for usage/indent/width/note layout.                                                       |
| `BeforeParse(fn func([]string) ([]string, error))`           | Mutate arguments before parsing (e.g., expand @files).                                                           |
| `OnUnknownFlag(fn func(name string) error)`                  | Handle or ignore unknown flags instead of failing.                                                               |
| `VersionText(text string)`                                   | Override the `--version` text. Default: `"Show version"`.                                                        |
| `HelpText(text string)`                                      | Override the `--help` text. Default: `"Show help"`.                                                              |
| `DisableHelp()` / `DisableVersion()`                         | Remove `--help` or `--version`.                                                                                  |
| `Usage func()`                                               | Optional custom usage function on `FlagSet` that replaces the default renderer.                                  |
| `Title(text string)`                                         | Override the "Usage:" title heading.                                                                             |
| `Authors(text string)`                                       | Add an `Authors:` section to help output.                                                                        |
| `Description(text string)`                                   | Add a free-form description block under the title.                                                               |
| `Note(text string)`                                          | Add a footer note under the flags listing.                                                                       |
| `SetOneOfGroupVerbose(bool)`                                 | Toggle detailed OneOfGroup errors with conflicting flags.                                                        |
| `SetOutput(w io.Writer)` / `Output()`                        | Redirect or retrieve where help/version is written.                                                              |
| `PrintUsage(w, mode)`                                        | Print the `Usage:` line.                                                                                         |
| `PrintTitle(w)`                                              | Print title and description.                                                                                     |
| `PrintAuthors(w)`                                            | Print authors section.                                                                                           |
| `PrintDescription(w,indent,width)`                           | Print the description block.                                                                                     |
| `PrintNotes(w,indent,width)`                                 | Print footer notes.                                                                                              |
| `PrintStaticDefaults(w,indent,startCol,width)`               | Print static flags help.                                                                                         |
| `PrintDynamicDefaults(w,indent,startCol,width)`              | Print dynamic flags help.                                                                                        |
| `MinPositional(n int)`                                       | Enforce at least `n` positional arguments.                                                                       |
| `ExactPositional(n int)`                                     | Enforce exactly `n` positional arguments.                                                                        |
| `MaxPositional(n int)`                                       | Enforce at most `n` positional arguments.                                                                        |
| `NoPositional()`                                             | Reject all positional arguments.                                                                                 |
| `RequirePositional(n int)`                                   | Deprecated compatibility alias for `MinPositional(n)`.                                                           |
| `Args() []string` / `Arg(i int) (string, bool)`              | Access leftover positional args safely.                                                                          |
| `Overrides() Overrides`                                      | Return explicit args/env overrides with both reporting values and origins (dynamic keys: `group.id.flag`).       |
| `Overrides.Values() map[string]any`                          | Project only the masked reporting values from an `Overrides` map.                                                |
| `Overrides.Origins() map[string]ValueOrigin`                 | Project only the exact origins from an `Overrides` map.                                                          |
| `Origin(name) ValueOrigin`                                   | Return the exact winning flag spelling or environment key for the effective value.                               |
| `MaskFirstLast(value any) any`                               | Helper mask that keeps first/last character (strings, `[]string`).                                               |
| `MaskFirstLastLimit(limit int) func(any) any`                | Returns a helper mask that keeps the first/last character and limits masked output length (strings, `[]string`). |
| `MaskPostgresURL(value any) any`                             | Helper mask for `postgres://user:pass@host/db` credentials.                                                      |
| `AddOneOfGroup(name string, group *core.OneOfGroupGroup)`    | Register a pre-built mutual-exclusion group.                                                                     |
| `AddAllOrNoneGroup(name string, group *core.AllOrNoneGroup)` | Register a pre-built require-together group.                                                                     |

## Helpers and exported utilities

- `FirstChanged[T](defaultValue, flags...)` — returns the value of the first changed flag (by order) plus whether any flag was set.
- `IsHelpRequested(err)` / `IsVersionRequested(err)` — detect help/version parse exits.
- `IsCommandRequired(err)` — detect missing required subcommand errors, even when wrapped with usage help.
- `HelpText(err)` — extract rendered help text from usage-bearing parse errors.
- `RequestHelp(msg)` / `RequestVersion(msg)` — trigger help/version errors manually.
- `Flag[T]` — minimal interface implemented by flag handles (`Changed() bool`, `Value() *T`).
- `Positive[T]()` / `NonNegative[T]()` — reusable numeric validators for positive and non-negative values.
- `AtLeast(min)` / `AtMost(max)` / `Between(min, max)` — reusable inclusive validators for ordered values.
- `NotBlank()` — rejects empty or whitespace-only strings.
- `MinLength(min)` / `MaxLength(max)` / `LengthBetween(min, max)` — reusable string-length validators.
- `Optional(validator)` — skips a wrapped validator when the value is the zero value for its type.
- `Override` / `Overrides` — one reporting value together with its `ValueOrigin`, and the map of all explicit overrides.
- `ValueOrigin` — exact provenance (`Source` plus the winning flag or environment key); its `String()` method returns labels such as `Flag · -p` and `Environment · APP_HOST`.
- `ValueSourceDefault` / `ValueSourceFlag` / `ValueSourceEnvironment` — source kinds stored in `ValueOrigin.Source`.
