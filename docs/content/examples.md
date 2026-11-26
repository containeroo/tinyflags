# Examples

Clone the repository to run the examples against the checked-out source:

```sh
git clone https://github.com/containeroo/tinyflags.git
cd tinyflags
```

## Simple

Basic flags, strict booleans and verbosity counters. [Source](https://github.com/containeroo/tinyflags/tree/main/examples/simple).

```sh
go run ./examples/simple --help
```

## Advanced

Sections, slices, defaults and finalizers. [Source](https://github.com/containeroo/tinyflags/tree/main/examples/advanced).

```sh
go run ./examples/advanced --help
```

## Commands

Nested commands and inherited globals. [Source](https://github.com/containeroo/tinyflags/tree/main/examples/commands).

```sh
go run ./examples/commands --help
```

## Subsubcommands

Deeper subcommand trees. [Source](https://github.com/containeroo/tinyflags/tree/main/examples/subsubcommands).

```sh
go run ./examples/subsubcommands --help
```

## Typedcommands

Typed runners and graceful cancellation. [Source](https://github.com/containeroo/tinyflags/tree/main/examples/typedcommands).

```sh
go run ./examples/typedcommands --help
```

## Dynamic

Dynamic flag groups and per-instance finalizers. [Source](https://github.com/containeroo/tinyflags/tree/main/examples/dynamic).

```sh
go run ./examples/dynamic --help
```

## Groups

Mutually exclusive and require-together flags. [Source](https://github.com/containeroo/tinyflags/tree/main/examples/groups).

```sh
go run ./examples/groups --help
```

## Toggle

Paired boolean flags and explicit changes. [Source](https://github.com/containeroo/tinyflags/tree/main/examples/toggle).

```sh
go run ./examples/toggle --help
```
