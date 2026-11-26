# Commands

`Command` lets you build subcommand trees with local flags and inherited globals. If you want parsing to fail unless a subcommand is chosen, call `RequireCommand()`.

Commands inherit the nearest explicitly configured `EnvPrefix`, including when it
is set after creating subcommands. `child.EnvPrefix("SERVER")` replaces the prefix
for that subtree; `child.EnvPrefix("")` stops inheritance and disables automatic
env lookup for that subtree. Explicit per-flag `.Env("KEY")` still works.
A command's own `Globals()` inherits its prefix, while ancestor globals keep their
owning command's prefix. You can override the prefix on `Globals()` separately.

Subcommands without an explicit prefix read environment variables using their ancestor's prefix.

```go
app := tinyflags.NewCommand("app", tinyflags.ExitOnError).RequireCommand()
app.Command("serve", "Run the server")

if err := app.Parse(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
```

If parsing fails because a required subcommand is missing, you can detect that and render contextual command help:

```go
if err := app.Parse(os.Args[1:]); err != nil {
    switch {
    case tinyflags.IsHelpRequested(err), tinyflags.IsVersionRequested(err):
        return err
    case tinyflags.IsCommandRequired(err):
        help, _ := tinyflags.HelpText(err)
        return fmt.Errorf("%s\n\n%s", err, help)
    default:
        return err
    }
}
```

## Local and inherited flags

```go
app := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
verbose := app.Globals().Bool("verbose", false, "Verbose logging").Value()
serve := app.Command("serve", "Run a server")
port := serve.Int("port", 8080, "Listen port").Value()
```

For `app --verbose serve --port=9090`, `verbose` and `port` are bound to their respective scopes. Use `SelectedCommand()` after parsing to identify the selected command. Nested command scopes can register their own persistent flags with `Globals()`.

## Positional argument hooks

Commands inherit the `FlagSet` positional hooks, including `SetPositionalValidateAll`
and `SetPositionalFinalizeAll`. Each participating flag set processes its own
positional slice. See [Parsing and positional arguments](parsing.md) for the
execution order and an IP-address deduplication example.

## Typed execution

Register a handler with `Run(handler, bindings...)` or a runner factory with `BuildCommand(builder)`. `ParseRunner` parses the command tree and returns the selected runner; the caller then executes it with a context.

```go
app := tinyflags.NewCommand("app", tinyflags.ContinueOnError).RequireCommand()
serve := app.Command("serve", "Run the server").NoPositional()
port := serve.Int("port", 8080, "Listen port").Value()

serve.Run(func(ctx context.Context, port int) error {
    fmt.Println(port)
    return ctx.Err()
}, port)

runner, err := app.ParseRunner(os.Args[1:])
if err != nil {
    return err
}
return runner.Run(ctx)
```

See the [typed commands example](https://github.com/containeroo/tinyflags/tree/main/examples/typedcommands) for a complete program, including cancellation and graceful HTTP shutdown.

## Command API

| Method                                                         | Description                                                                           |
| :------------------------------------------------------------- | :------------------------------------------------------------------------------------ |
| `NewCommand(name string, mode ErrorHandling)`                  | Create a root command tree.                                                           |
| `Command(name, summary string)`                                | Register a child command.                                                             |
| `Globals()`                                                    | Access persistent flags inherited by that subtree.                                    |
| `RequireCommand()`                                             | Return an error if this command is selected without a child command.                  |
| `MinPositional(n)` / `ExactPositional(n)` / `MaxPositional(n)` | Configure positional argument counts and return the command for fluent chaining.      |
| `NoPositional()`                                               | Reject all positional arguments and return the command for fluent chaining.           |
| `RequirePositional(n)`                                         | Deprecated fluent compatibility alias for `MinPositional(n)`.                         |
| `HelpText()`                                                   | Return rendered help for the selected command when available, otherwise the receiver. |
| `WriteHelp(w io.Writer)`                                       | Write rendered help for the selected command when available, otherwise the receiver.  |
| `Parse(args)`                                                  | Parse flags and select the active command.                                            |
| `SelectedCommand()`                                            | Return the selected leaf command from the last parse.                                 |
| `Run(handler, bindings...)` / `BuildCommand(builder)`          | Register execution for a command.                                                     |
| `ParseRunner(args)` / `ParseRunnable(args)`                    | Parse and build the selected runnable.                                                |
