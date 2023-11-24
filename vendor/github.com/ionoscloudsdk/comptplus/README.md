# Cobra-Prompt Plus (Comptplus)

![Comptplus Banner](https://github.com/avirtopeanu-ionos/cobra-prompt/assets/100703584/9a4b23f1-5f7e-4e76-89f3-010a799158f5)

> Comptplus is a fork of [Cobra Prompt](https://github.com/stromland/cobra-prompt) with added features, intended to enhance the original implementation by making it more versatile.

## Projects using this fork
- [Ionos Cloud CLI](https://github.com/ionos-cloud/ionosctl/)

## Features unique to this fork
### Flag Value Completions and Persistence
- **Flag Completions**: Added completions for flag values, facilitating easier and more intuitive command usage.
    - Default cache duration for responses is set to 500ms, reducing latency in user interaction.
    - Support for flag descriptions by splitting on `\t`.

### Customizable Flag Reset Behaviors
- **Custom Flag Reset Behaviors**: Ability to set custom behaviors for resetting flag values.
    - The default flag reset behaviour has also been changed to reset flags to their default values after each command execution, improving consistency and predictability of the shell environment.
    - A bug in the original repo caused slice/array flags to be reset incorrectly, by appending the default values to the previous execution's values. This has been fixed in this fork.

### Pre and Post Execution Hooks
- **Execution Hooks**: Added `HookBefore` and `HookAfter` for performing actions before and after command execution.

### Enhanced Flag Persistence
- **Persist Flag Values**: Refined behavior of `PersistFlagValues` to offer more seamless flag value persistence across shell sessions.

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
