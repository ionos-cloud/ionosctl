# Cobra-Prompt Plus (Comptplus)

![Comptplus Banner](https://github.com/avirtopeanu-ionos/cobra-prompt/assets/100703584/9a4b23f1-5f7e-4e76-89f3-010a799158f5)

> Comptplus is a fork of [Cobra Prompt](https://github.com/stromland/cobra-prompt) with added features, intended to enhance the original implementation by making it more versatile.

## Projects using this fork
- [Ionos Cloud CLI](https://github.com/ionos-cloud/ionosctl/)

## Features unique to this fork
### Flag Value Completions and Persistence
- **Flag Completions**: Added completions for flag values, for easier and more intuitive command usage.
    - Default cache duration for refreshing flag value completions is set to 500ms, to avoid lag in user interaction.
    - Support for flag descriptions by splitting on `\t`.

### Customizable Flag Reset Behaviors
- **Custom Flag Reset Behaviors**: Ability to set custom behaviors for resetting flag values.
    - The default flag reset behaviour has also been changed to reset flags to their default values after each command execution (rather than on each typed character).
    - A bug in the original repo caused slice/array flags to be reset incorrectly, by appending the default values to the previous execution's values. This has been fixed in this fork.

### Pre and Post Execution Hooks
- **Execution Hooks**: Ability to use `HookBefore` and `HookAfter` for performing actions before and after command execution.

### elk-language/go-prompt
- **Bump go-prompt**: Switched to an actively maintained fork of go-prompt with more features and bug fixes: https://github.com/elk-language/go-prompt. Sadly this involves some breaking changes.
    - This fixes a bug which would disallow usage of CTRL + C in after usage of go-prompt.

### Proper Argument Passing via SetArgs
- **SetArgs instead of os.Args**: Command execution now uses `cobra.Command.SetArgs()` instead of manipulating the global `os.Args`. This avoids side effects with global state and works correctly when comptplus is used as a sub-command of a larger CLI (e.g. apps that call `root.SetArgs` as part of their bootstrap).

### Shell-Aware Argument Parsing
- **Quoted argument support**: Default argument parser now uses [go-shellquote](https://github.com/kballard/go-shellquote) instead of `strings.Fields`. This means arguments with spaces work correctly when quoted, e.g. `get -n "John Oliver" food apple` is parsed as `["get", "-n", "John Oliver", "food", "apple"]` rather than splitting on every space. Falls back to `strings.Fields` on parse errors (e.g. unclosed quotes). Custom parsers via `InArgsParser` are still supported.

### DynamicSuggestionsFunc with Command Context
- **Command-aware dynamic suggestions**: `DynamicSuggestionsFunc` now receives the resolved `*cobra.Command` as its first parameter, allowing consumers to use `cmd.ValidArgsFunction()` for generic auto-completion logic.

### Graceful Error Handling
- **No more os.Exit on error**: `handleUserError` no longer calls `os.Exit(1)` when no `OnErrorFunc` is set. Instead it prints the error and returns control to the prompt loop, so a single bad command doesn't kill the entire shell session.


## Original README below

-----

# Cobra-Prompt

Cobra-prompt makes every Cobra command and flag available for go-prompt.
- https://github.com/spf13/cobra
- https://github.com/c-bata/go-prompt


## Features

- Traverse cobra command tree. Every command and flag will be available.
- Persist flag values.
- Add custom functions for dynamic suggestions.

## Getting started

Get the module:

```
go get github.com/stromland/cobra-prompt
```

## Explore the example

```
cd _example
go build -o cobra-prompt
./cobra-prompt
```
