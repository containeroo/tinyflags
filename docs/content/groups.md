# Flag groups

When certain flags must be used **together**, or must be **exclusive**, tinyflags makes that easy.

## Require-together groups

You can define a group where **either all or none** of the flags must be set:

```go
email := fs.String("email", "", "User email").
    AllOrNone("authpair").
    Value()

password := fs.String("password", "", "User password").
    AllOrNone("authpair").
    Value()
```

## Mutual-exclusion groups

You can define a group where **only one** of the flags (or groups!) may be set:

```go
bearer := fs.String("bearer-token", "", "Bearer token").
    OneOfGroup("authmethod").
    Value()

fs.GetOneOfGroup("authmethod").
    Title("Authentication Method").
    Required().
    AddGroup(fs.GetAllOrNoneGroup("authpair")) // <-- email+password
```

By default, a one-of group permits zero or one selection. Calling `Required()` on the group makes one selection mandatory. An all-or-none group can itself be required with `fs.GetAllOrNoneGroup("authpair").Required()`, but do not do that when it is intended to be an alternative to a bearer token.

This example enforces:

- Either `--email` and `--password` **together**, or
- `--bearer-token`, but **not both**

### Valid combinations

```text
# Valid:
--email=user --password=secret
--bearer-token=abc123

# Invalid:
--email=user                ❌ password missing
--email=user --bearer-token=abc123 ❌ both auth methods
```

### Illustrative help output

```text
Usage: app [flags]

Flags:
  --email EMAIL         User email (Group: Authentication Method) (AllOrNone: authpair)
  --password PASSWORD   User password (Group: Authentication Method) (AllOrNone: authpair)
  --bearer-token TOKEN  Bearer token (Group: Authentication Method)

Authentication Method
(Exactly one required)
  --bearer-token TOKEN
  [--email, --password] (Required together)

authpair
  --email EMAIL
  --password PASSWORD
```

## Multiple one-of groups

You can attach a flag to more than one one-of group by chaining `OneOfGroup(...)`. By default, all visible memberships are rendered in help. If you want help to mention only a subset, add `HelpOneOfGroups(...)`.

```go
searchHistoryOnly := fs.Bool("history-only", false, "Search shell history commands only.").
    OneOfGroup("search-mode").
    Value()

searchPinsOnly := fs.Bool("pins-only", false, "Search pinned commands only.").
    OneOfGroup("search-mode").
    OneOfGroup("pins-scope").
    HelpOneOfGroups("search-mode").
    Value()

searchExcludePins := fs.Bool("exclude-pins", false, "Exclude pinned commands from history results.").
    OneOfGroup("pins-scope").
    Value()
```
