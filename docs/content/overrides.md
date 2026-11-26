# Value provenance

`Overrides()` contains only explicitly provided CLI or environment values.

- `Overrides()` reports values provided by CLI or env, not untouched defaults. Each entry contains both the reporting value and its exact `ValueOrigin`.
- `Overrides().Values()` projects only reporting values, `Overrides().Origins()` projects only origins, and `Origin(name)` also reports `Default` for untouched values.
- CLI origins contain only the flag spelling, never its value, so sensitive `--token=...` input is reported as `--token`.

For example:

```go
fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
fs.EnvPrefix("APP")
fs.String("host", "localhost", "Server host")
fs.Int("port", 8080, "Server port").Short("p")
fs.String("dsn", "", "Database DSN").Env("DATABASE_DSN").
    OverriddenValueMaskFn(tinyflags.MaskPostgresURL)
fs.String("mode", "development", "Runtime mode")

// APP_HOST=db.internal DATABASE_DSN=postgres://... app -p 9090
if err := fs.Parse(os.Args[1:]); err != nil {
    return err
}

overrides := fs.Overrides()
fmt.Println(overrides.Values())
fmt.Println(overrides["host"].Origin) // Environment · APP_HOST
fmt.Println(overrides["port"].Origin) // Flag · -p
fmt.Println(overrides["dsn"].Origin)  // Environment · DATABASE_DSN
fmt.Println(fs.Origin("mode"))        // Default
```

## Masking sensitive values

Register a mask using `.OverriddenValueMaskFn(...)` to protect the reporting value. Masking affects `Overrides()` and its `Values()` projection, not the parsed value. Never log raw arguments or unmasked secrets.

```go
fs.String("token", "", "API token").
    OverriddenValueMaskFn(func(any) any { return "[redacted]" })
```

Origins are safe to display as flag or environment key names, but values are not automatically redacted. See [the API reference](reference.md) for masking helpers.
