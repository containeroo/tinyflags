# Parsing and positional arguments

Tinyflags applies input in this order:

1. Parse command-line arguments.
2. Handle built-in `--help` / `--version` exits.
3. Load unset static and dynamic flags from environment variables.
4. Apply default finalizers for unset values.
5. Run required/group/dependency checks and positional validation and finalization.
6. Run registered cross-flag validators.

Additional behavior:

- Explicit CLI arguments win over environment variables.
- Reusing a `FlagSet` across multiple `Parse(...)` calls is supported; parser state is reset before each parse.
- Automatic static env lookup requires `EnvPrefix(...)`; explicit static `.Env("KEY")` works without a prefix.
- Dynamic env lookup requires `EnvPrefix(...)` and uses `PREFIX_GROUP_ID_FIELD` keys such as `MYAPP_HTTP_API_PORT`.

## Positional argument counts

Count constraints are checked during command-line parsing, before environment loading and positional hooks. Help and version requests bypass the count checks.

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

## Positional validation and finalization

Both `FlagSet` and `Command` expose four optional hooks:

| Method                     | Callback signature                 |
| -------------------------- | ---------------------------------- |
| `SetPositionalValidate`    | `func(string) error`               |
| `SetPositionalFinalize`    | `func(string) string`              |
| `SetPositionalValidateAll` | `func([]string) error`             |
| `SetPositionalFinalizeAll` | `func([]string) ([]string, error)` |

After checking the input count, tinyflags validates each argument, finalizes each
argument, validates the entire resulting slice, and finally replaces it with the
result of the collection finalizer. Cross-flag validators run afterwards and can
read the finalized values through `Args()` and `Arg(i)`.

The collection hooks run once per participating flag set, even when there are no
positional arguments. They receive a copy of the current slice (which may be nil).
The validator cannot change the stored arguments by editing its slice. A finalizer
may return fewer or more values, including nil. Its result is stored only on success;
any callback error aborts parsing and is handled by the configured error mode.
Count constraints are not rechecked against the finalizer's output. Passing nil
removes a hook; setting a hook again replaces the previous one. Help and version
requests skip the hooks. Command scopes retain their own positional arguments.

For example, normalize and deduplicate IP addresses while preserving input order
using the standard library's `net/netip` package:

```go
command := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
command.MinPositional(1)
command.SetPositionalFinalizeAll(func(args []string) ([]string, error) {
    addresses := make([]string, 0, len(args))
    seen := make(map[netip.Addr]struct{}, len(args))

    for _, arg := range args {
        address, err := netip.ParseAddr(arg)
        if err != nil {
            return nil, err
        }

        if _, exists := seen[address]; exists {
            continue
        }

        seen[address] = struct{}{}
        addresses = append(addresses, address.String())
    }

    return addresses, nil
})

if err := command.Parse([]string{"192.0.2.1", "2001:0db8::1", "2001:db8::1"}); err != nil {
    return err
}
// command.Args() == []string{"192.0.2.1", "2001:db8::1"}
```

Use `SetPositionalValidateAll` when you only need to check relationships between
arguments, such as requiring key-value pairs, without transforming the values.

## Handling parse errors

`ContinueOnError` returns parse failures, `ExitOnError` can exit with status 2, and `PanicOnError` can panic on parse failures. Help and version requests are returned as sentinel errors rather than ordinary failures. Use `ContinueOnError` when you want to handle errors in your application. `IsHelpRequested`, `IsVersionRequested`, `IsCommandRequired`, and `HelpText(err)` let you distinguish built-in requests and extract contextual help.
