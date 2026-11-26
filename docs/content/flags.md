# Flags and types

## Supported types

| Type                 | Scalar constructors                      | Slice constructors                  |
| -------------------- | ---------------------------------------- | ----------------------------------- |
| `bool`               | `Bool`, `BoolVar`                        | —                                   |
| `int`                | `Int`, `IntVar`                          | `IntSlice`, `IntSliceVar`           |
| `int32`              | `Int32`, `Int32Var`                      | `Int32Slice`, `Int32SliceVar`       |
| `int64`              | `Int64`, `Int64Var`                      | `Int64Slice`, `Int64SliceVar`       |
| `float32`            | `Float32`, `Float32Var`                  | `Float32Slice`, `Float32SliceVar`   |
| `float64`            | `Float64`, `Float64Var`                  | `Float64Slice`, `Float64SliceVar`   |
| `string`             | `String`, `StringVar`, `Enum`, `EnumVar` | `StringSlice`, `StringSliceVar`     |
| Counter (`int`)      | `Counter`, `CounterVar`                  | —                                   |
| `time.Duration`      | `Duration`, `DurationVar`                | `DurationSlice`, `DurationSliceVar` |
| `time.Time`          | `Time`, `TimeVar`                        | `TimeSlice`, `TimeSliceVar`         |
| Byte size (`uint64`) | `Bytes`, `BytesVar`                      | `BytesSlice`, `BytesSliceVar`       |
| `net.IP`             | `IP`, `IPVar`                            | `IPSlice`, `IPSliceVar`             |
| `net.IPMask`         | `IPMask`, `IPv4MaskVar`                  | `IPMaskSlice`, `IPMaskSliceVar`     |
| `*net.TCPAddr`       | `TCPAddr`, `TCPAddrVar`                  | `TCPAddrSlice`, `TCPAddrSliceVar`   |
| `*url.URL`           | `URL`, `URLVar`                          | `URLSlice`, `URLSliceVar`           |
| `*os.File`           | `File`, `FileVar`                        | `FileSlice`, `FileSliceVar`         |

Slice flags accept repeated use or custom-delimited strings. Byte-size flags accept raw bytes, SI units such as `32MB`, and IEC units such as `32MiB`.

## Common flag-builder methods

| Method                      | Applies to    | Description                                                                                                |
| :-------------------------- | :------------ | :--------------------------------------------------------------------------------------------------------- |
| `Short(s string)`           | static only   | One-letter alias (`-p`). Must be exactly one rune (panics otherwise).                                      |
| `Required()`                | all flags     | Mark the flag as required; parser errors if unset.                                                         |
| `NotEmpty()`                | all flags     | Reject an explicitly set empty/zero value; combine with `Required()` to require both presence and content. |
| `NotBlank()`                | string values | Reject explicitly set empty or whitespace-only strings.                                                    |
| `HideRequired()`            | all flags     | Hide the "(Required)" suffix from help.                                                                    |
| `Hidden()`                  | all flags     | Omit this flag from generated help output.                                                                 |
| `Deprecated(reason string)` | all flags     | Mark flag deprecated; includes `DEPRECATED` note in help.                                                  |
| `OneOfGroup(group string)`  | all flags     | Assign to a named mutual-exclusion group. Parsing errors if more than one in group set.                    |
| `HelpOneOfGroups(names...)` | all flags     | Override which one-of groups for this flag are shown in help output.                                       |
| `AllOrNone(group string)`   | all flags     | Assign to a named require-together group. All or none in group must be set.                                |
| `Env(key string)`           | static only   | Override the environment-variable name (panics if `DisableEnv` already called).                            |
| `HideEnv()`                 | all flags     | Hide the environment-variable name from help output.                                                       |
| `DisableEnv()`              | all flags     | Disable environment lookup for this flag (panics if `Env(...)` already called).                            |
| `Placeholder(text string)`  | all flags     | Customize the `<VALUE>` placeholder in help.                                                               |
| `Allowed(vals ...string)`   | all flags     | Advertise allowed values in help; use `Choices` to enforce them.                                           |
| `HideAllowed()`             | all flags     | Hide the allowed values from help.                                                                         |
| `Requires(names ...string)` | all flags     | Mark flag as required by the given flag.                                                                   |
| `HideRequires()`            | all flags     | Hide the “(Requires)” suffix from help.                                                                    |
| `OverriddenValueMaskFn(fn)` | all flags     | Provide a mask function used by `Overrides()` and its `Values()` projection.                               |
| `Value() *T`                | static only   | Return the pointer to the parsed value (after `Parse`).                                                    |

`Required()` and `NotEmpty()` are intentionally different: `Required()` checks whether a flag was supplied, while `NotEmpty()` checks the resulting value. For example, `--token=` satisfies `Required()` but fails `NotEmpty()`. Chaining both enforces a present, non-empty value.

## Type-specific builder methods

| Method                                                | Description                                                                                                                | Example                                                                                                         |
| :---------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------- |
| `Choices(v1, v2, ...)`                                | Only allow the provided literal values; automatically adds them to help output.                                            | `fs.String("env","dev","...").Choices("dev","staging","prod")`                                                  |
| `Enum(name, def, usage, values...)`                   | Convenience constructor for a string flag with fixed allowed values.                                                       | `fs.Enum("env","dev","...","dev","staging","prod")`                                                             |
| `tinyflags.Enum(fs, name, def, usage, ...)`           | Generic helper for typed string enum constants.                                                                            | `tinyflags.Enum(fs,"env",dev,"...",dev,staging,prod)`                                                           |
| `tinyflags.EnumMap(fs, name, def, usage, choices...)` | Generic helper for typed enums with display names, useful for `iota` enums.                                                | `tinyflags.EnumMap(fs,"level",info,"...",tinyflags.Choice("debug",debug),tinyflags.Choice("info",info))`        |
| `Validate(fn func(v T) error)`                        | Run custom check on parsed value; if `fn` returns non-nil, `Parse` returns an error.                                       | `fs.Int("count",0,"...").Validate(func(n int) error {   if n<0 {return fmt.Errorf("must ≥0")}   return nil }) ` |
| `Finalize(fn func(v T) T)`                            | Transform the parsed value before storing; e.g. trimming, normalization, applying defaults.                                | `fs.String("name","","...").Finalize(func(s string) string {   return strings.TrimSpace(s) }) `                 |
| `FinalizeDefaultValue()`                              | Run the existing finalizer on default values when the flag is unset.                                                       | `fs.String("name","","...").Finalize(strings.TrimSpace).FinalizeDefaultValue() `                                |
| `FinalizeWithID(fn func(id string, v T) T)`           | _(dynamic only)_ Finalize with access to the instance ID.                                                                  | `http.String("addr","","").FinalizeWithID(func(id, v string) string { return id+":"+v })`                       |
| `Delimiter(sep string)`                               | _(slice flags only)_ Use a custom separator instead of the default comma when parsing lists.                               | `fs.StringSlice("tags",nil,"...").Delimiter(";")`                                                               |
| `TrimSpace()`                                         | _(slice flags only)_ Trim leading/trailing whitespace from each parsed item. String slices preserve whitespace by default. | `fs.StringSlice("tags",nil,"...").TrimSpace()`                                                                  |
| `PreserveSpace()`                                     | _(slice flags only)_ Preserve leading/trailing whitespace in each parsed item. Typed slices trim whitespace by default.    | `fs.IntSlice("ports",nil,"...").PreserveSpace()`                                                                |
| `AllowEmpty()`                                        | _(slice flags only)_ Allow empty items (e.g. `"a,,b"`).                                                                    |                                                                                                                 |
| `HideDefault()`                                       | Hide the default value from help output.                                                                                   |                                                                                                                 |
| `Section(name string)`                                | Group flags under a section header in help output.                                                                         |                                                                                                                 |

## Binding to existing variables

Use the `Var` constructor to bind a flag to an existing variable or struct field.

```go
var cfg struct {
    Port int
    Host string
}
fs.IntVar(&cfg.Port, "port", 8080, "Listen port")
fs.StringVar(&cfg.Host, "host", "localhost", "Listen host")
```

## Slice flags

```go
tags := fs.StringSlice("tag", nil, "Tags").TrimSpace().Value()
```

`--tag=api,web --tag=worker` produces `[]string{"api", "web", "worker"}`. Use `Delimiter(";")` for a different separator, or `NoSplit()` to keep each occurrence intact.

## Boolean flags and counters

Use `.Strict()` on a boolean flag to require an explicit value, for example `--silent=true` or `--silent=false`. Ordinary boolean flags accept `--debug` (implicit `true`), `--debug=true`, and `--debug=false`. Explicit values are parsed and validated, and override environment values even when they equal the default. The same syntax applies to dynamic boolean flags. A counter increments each time its flag is supplied.

```go
silent := fs.Bool("silent", false, "Silent mode").Strict().Value()
verbose := fs.Counter("verbose", 0, "Verbosity").Short("v").Value()
```

## Handling toggles with multiple flags

You can model toggles with paired flags (e.g., `--debug` and `--no-debug`) and pick the first one the user set:

```go
debug := fs.Bool("debug", false, "Enable debug").Short("d").OneOfGroup("debug")
noDebug := fs.Bool("no-debug", false, "Disable debug").Short("n").OneOfGroup("debug")

// Use the order you prefer; first changed wins
enabled, set := tinyflags.FirstChanged(false, debug, noDebug)
fmt.Printf("debug: %t (set: %v)\n", enabled, set)
```
