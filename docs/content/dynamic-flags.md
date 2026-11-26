# Dynamic flags

Dynamic flags model configuration whose instance names are not known at registration time. Register fields once and supply `--group.id.field=value` at runtime. Each ID gets its own typed values.

```go
package main

import (
    "fmt"
    "os"
    "github.com/containeroo/tinyflags"
)

func main() {
    fs := tinyflags.NewFlagSet("app", tinyflags.ContinueOnError)
    http := fs.DynamicGroup("http")

    port    := http.Int("port", 0, "Backend port")
    timeout := http.Duration("timeout", 0, "Request timeout")

    if err := fs.Parse(os.Args[1:]); err != nil {
        if tinyflags.IsHelpRequested(err) {
            fmt.Println(err)
            return
        }
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

### Disable slice splitting

Call `.NoSplit()` on any static or dynamic slice to pass each flag input to its
parser as one string. For example, `fs.StringSlice("value", nil, "Values").NoSplit()`
parses `--value=a,b --value=c,d` as `[]string{"a,b", "c,d"}`.
Whitespace handling and `AllowEmpty()` still apply. Calling `.Delimiter(sep)`
after `.NoSplit()` enables splitting again.

Custom dynamic slice parsers can use `.NoSplit()` to handle their own separators:

```go
tinyflags.DynamicSlice(httpGroup, "expected-status-codes", []int{200},
    "Expected HTTP status codes", httputils.ParseStatusCodes, strconv.Itoa).
    NoSplit()
```

### Custom dynamic slices

Use `tinyflags.DynamicSlice` for a parser with signature `func(string) ([]T, error)`
that expands one input into multiple elements, such as HTTP status ranges:

```go
httpGroup := fs.DynamicGroup("http")
tinyflags.DynamicSlice(
    httpGroup,
    "expected-status-codes",
    []int{200},
    "Expected HTTP status codes",
    httputils.ParseStatusCodes,
    strconv.Itoa,
)

// After parsing --http.web.expected-status-codes=200-299,304:
codes := tinyflags.GetOrDefaultDynamic[[]int](httpGroup, "web", "expected-status-codes")
```

Here, `httputils.ParseStatusCodes` is your parser and `strconv.Itoa` formats each
integer. Each comma-separated chunk is trimmed and passed to the parser; its
results are appended in order. Repeated flags append more elements. The returned
handle supports `Delimiter`, `PreserveSpace`, `Validate`, `Choices`, and finalizers,
with validation and finalization applied to each expanded element. Unset IDs use
the default; a successful empty expansion stores an empty slice.

## Dynamic flag accessors

| Method                               | Description                                                                       | Example                                                      |
| :----------------------------------- | :-------------------------------------------------------------------------------- | :----------------------------------------------------------- |
| `Has(id string) bool`                | Return whether the given instance-ID was provided on the command line or via env. | `if port.Has("a") { fmt.Println("port a is set") }`          |
| `Get(id string) (value T, ok bool)`  | Retrieve the parsed value for that ID; `ok==false` if unset (returns default).    | `p, ok := port.Get("a"); if ok { fmt.Println("a →", p) }`    |
| `MustGet(id string) T`               | Like `Get`, but panics if the instance wasn't provided.                           | `timeout.MustGet("b")`                                       |
| `Values() map[string]T`              | Get all parsed values keyed by instance ID.                                       | `for id, v := range timeout.Values() { fmt.Println(id, v) }` |
| `ValuesAny() map[string]interface{}` | Same as `Values()`, but values are `interface{}`.                                 |                                                              |
| `AllowOverride()`                    | Allow re-assignment of a dynamic flag per-ID. Only for Scalar flags.              |                                                              |

## Required fields and instance-aware finalizers

```go
http := fs.DynamicGroup("http")
port := http.Int("port", 8080, "Listen port")
port.Validate(tinyflags.Between(1, 65535))
port.Required()

name := http.String("name", "", "Display name")
name.FinalizeWithID(func(id, value string) string {
    return id + ":" + value
})
```

Required fields are checked for each discovered instance, not for nonexistent instances. Use `Instances()` to enumerate IDs. Dynamic flags use `Get(id)` rather than a static `Value()` pointer. Dynamic scalar and slice builders expose type-specific methods; keep the original typed handle when calling metadata methods such as `Required()` because those methods return the embedded metadata builder.

See [Environment variables](environment.md) for discovery through environment variables. Use `.AllowOverride()` on a dynamic scalar if repeated assignments should be accepted; by default, assigning the same field and ID twice is an error.
