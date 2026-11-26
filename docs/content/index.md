# tinyflags

**A minimal, extensible CLI flag-parsing library for Go.**

Tinyflags provides typed flags, fluent builders, environment-variable overrides, subcommands and configurable help output without third-party runtime dependencies.

```sh
go get github.com/containeroo/tinyflags
```

## Features

- **Short & long flags** (`-d`, `--debug`)
- **Boolean strict mode** (`--flag=true/false`)
- **Environment variable overrides** (`EnvPrefix`, per-flag opt-out)
- **Value provenance** (exact `-x` / `--long` flag spelling or environment-variable key)
- **Required, deprecated, and grouped flags**
- **Slice flags** (`[]T`) with custom delimiters
- **Allowed choices, validation and finalizers**
- **One Of groups**
- **All or None groups**
- **Custom placeholders & help sections**
- **Dynamic flags** (`--group.id.field=value`)
- **Nested commands and typed runners**
- **Typed values** (`*os.File`, `*net.TCPAddr`, `*url.URL`, `time.Duration`, etc.)

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

## Documentation

- [Getting started](getting-started.md): install tinyflags and parse your first CLI.
- [Flags and types](flags.md): scalars, slices, enums and fluent flag builders.
- [Commands](commands.md): nested subcommands, global flags and typed runners.
- [Parsing and positional arguments](parsing.md): precedence, positional constraints and parse errors.
- [Environment variables](environment.md): automatic and explicit bindings.
- [Dynamic flags](dynamic-flags.md): instance-scoped configuration and custom slices.
- [Validation and finalization](validation.md): reusable validators and value transformations.
- [Flag groups](groups.md): mutually exclusive and require-together options.
- [Value provenance](overrides.md): inspect explicit overrides and mask secrets.
- [Help and API reference](reference.md): help layout and API tables.
- [Examples](examples.md): runnable examples in the repository.

[Source on GitHub](https://github.com/containeroo/tinyflags) · [Go API reference](https://pkg.go.dev/github.com/containeroo/tinyflags) · [Releases](https://github.com/containeroo/tinyflags/releases)
