# Environment variables

Flags can be set from environment variables in addition to CLI arguments.
Automatic environment lookup is enabled by setting a global prefix:

- static flags use `PREFIX_FLAG_NAME`
- dynamic flags use `PREFIX_GROUP_ID_FIELD`
- `-`, `.`, and `/` are replaced with `_`
- the derived key is upper-cased

For example:

```bash
MYAPP_HOST=example.com MYAPP_PORT=9090 ./app --debug
```

A static flag `--db.user` will look up `MYAPP_DB_USER`.
A dynamic flag `--http.api.port` will look up `MYAPP_HTTP_API_PORT` and create the `api` instance if needed.

Static flags can also bind to an exact key without a global prefix:

```go
fs.String("token", "", "API token").Env("TOKEN")
```

Dynamic flags do not use bare keys; dynamic env lookup requires `EnvPrefix(...)` so instance discovery stays bounded.
CLI arguments always win over env values.

You can disable env binding per-flag:

```go
fs.Bool("internal", false, "internal use only").DisableEnv()
```

### Custom key mapping

You can override how static keys are derived with `SetEnvKeyFunc`:

```go
fs.EnvPrefix("MYAPP")
fs.SetEnvKeyFunc(fs.NewReplacerEnvKeyFunc(
    strings.NewReplacer("-", "_", ".", "_", "/", "_"),
    true, // upper-case
))
```

This would map static flags as follows:

- `--log.level` → `MYAPP_LOG_LEVEL`
- `--db-user` → `MYAPP_DB_USER`
- `--api/v1` → `MYAPP_API_V1`

Dynamic env keys keep the canonical `PREFIX_GROUP_ID_FIELD` mapping.
