# Validation and finalization

## Cross-flag validation

Flag sets can register multiple cross-flag validators with
`Validate(func() error)`. They run in registration order after command-line
and environment values, defaults, and built-in relationships have been
processed. Parsing stops at the first validation error. Help and version
requests do not run validators.

```go
mode := fs.String("mode", "development", "Runtime mode").Value()
token := fs.String("token", "", "Production token").Value()

fs.Validate(func() error {
    if *mode == "production" && *token == "" {
        return fmt.Errorf("production mode requires --token")
    }
    return nil
})
```

## Per-flag validation and finalization

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

   Tinyflags also includes reusable validators for common constraints:

   ```go
   timeout := fs.Duration("timeout", 5*time.Second, "Request timeout").
       Validate(tinyflags.Positive[time.Duration]()).
       Value()

   retries := fs.Int("retries", 3, "Retry count").
       Validate(tinyflags.NonNegative[int]()).
       Value()

   port := fs.Int("port", 8080, "Listen port").
       Validate(tinyflags.Between(1, 65535)).
       Value()

   name := fs.String("name", "app", "Application name").
       Validate(tinyflags.NotBlank()).
       Value()

   sessionSecret := fs.String("oidc-session-secret", "", "OIDC session secret").
       Validate(tinyflags.Optional(tinyflags.MinLength(32))).
       Value()
   ```

   `AtLeast` and `AtMost` provide one-sided inclusive bounds. `Between` is inclusive
   and panics when its minimum is greater than its maximum. Numeric validators also
   support named numeric types such as `time.Duration`. `MinLength`, `MaxLength`, and
   `LengthBetween` use Go's `len` for string length. `Optional` skips the wrapped
   validator only when the value equals the type's zero value, such as `""`, `0`, or
   `false`.

2. **Finalize**
   - Only after validation succeeds, Tinyflags passes the parsed value through your finalizer:

   ```go
   fs := tinyflags.NewFlagSet("app", tinyflags.ExitOnError)

   name := fs.String("name", "", "User name").
     Finalize(func(s string) string {
             return strings.TrimSpace(s)
             }).
     Value()

   url := fs.URL("url", nil, "URL to use").
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


    if err := fs.Parse([]string{"--name=   Alice Smith  ", "--url", "https://containeroo.ch"}); err != nil {
        return err
    }
    fmt.Println("Hello,", *name)
    fmt.Println("Visit:", *url)
   ```

3. **FinalizeDefaultValue**
   - Runs only when a flag is _unset_ (after env parsing), and does not mark it as changed.
   - Uses the flag’s `Finalize(...)` function; no separate function is provided.

   ```go
   fs.String("name", " app ", "Name").
       Finalize(strings.TrimSpace).
       FinalizeDefaultValue()
   ```

   - Useful for trimming whitespace, applying normalization, setting derived defaults, etc.

## Validator helpers

- `Positive[T]()` / `NonNegative[T]()` — reusable numeric validators for positive and non-negative values.
- `AtLeast(min)` / `AtMost(max)` / `Between(min, max)` — reusable inclusive validators for ordered values.
- `NotBlank()` — rejects empty or whitespace-only strings.
- `MinLength(min)` / `MaxLength(max)` / `LengthBetween(min, max)` — reusable string-length validators.
- `Optional(validator)` — skips a wrapped validator when the value is the zero value for its type.
