# Parsing and positional arguments

Tinyflags applies input in this order:

1. Parse command-line arguments.
2. Handle built-in `--help` / `--version` exits.
3. Load unset static and dynamic flags from environment variables.
4. Apply default finalizers for unset values.
5. Run required/group/dependency/positional validation.
6. Run registered cross-flag validators.

Additional behavior:

- Explicit CLI arguments win over environment variables.
- Reusing a `FlagSet` across multiple `Parse(...)` calls is supported; parser state is reset before each parse.
- Automatic static env lookup requires `EnvPrefix(...)`; explicit static `.Env("KEY")` works without a prefix.
- Dynamic env lookup requires `EnvPrefix(...)` and uses `PREFIX_GROUP_ID_FIELD` keys such as `MYAPP_HTTP_API_PORT`.

## Positional argument counts

Use `MinPositional`, `ExactPositional`, `MaxPositional`, and `NoPositional` to constrain leftover positional arguments:

```go
fs.MinPositional(1)   // one or more
fs.MaxPositional(3)   // at most three

// Or require one exact count:
fs.ExactPositional(2)

// Or reject positional arguments entirely:
fs.NoPositional()
```

The same methods are available on `Command` and return the command for fluent chaining. `MinPositional` and `MaxPositional` can be combined to define a range. `RequirePositional(n)` remains available for backward compatibility and is deprecated in favor of `MinPositional(n)`.

## Accessing positional values

After parsing, `fs.Args()` returns the remaining arguments and `fs.Arg(i)` returns a value and a presence boolean. Use `SetPositionalValidate` and `SetPositionalFinalize` to validate and normalize each positional argument.

## Handling parse errors

`ContinueOnError` returns parse failures, `ExitOnError` can exit with status 2, and `PanicOnError` can panic on parse failures. Help and version requests are returned as sentinel errors rather than ordinary failures. Use `ContinueOnError` when you want to handle errors in your application. `IsHelpRequested`, `IsVersionRequested`, `IsCommandRequired`, and `HelpText(err)` let you distinguish built-in requests and extract contextual help.
