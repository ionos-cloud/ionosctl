# Cobra-Prompt Plus (Comptplus)

Comptplus is a fork of [Cobra Prompt](https://github.com/stromland/cobra-prompt) with added features:
  - Automatic Flag Value Completion: Extracts possible flag values and recommends them to the user.
    - Works with descriptions registered via `RegisterFlagCompletionFunc` 
    - This is only possible with Cobra 1.8.0 (specifically PRs [spf13/cobra#1943](https://github.com/spf13/cobra/pull/1943) and [spf13/cobra#2063](https://github.com/spf13/cobra/pull/2063))
    - Flag descriptions are also added correctly (in Cobra `\t` is used to split the flag values vs. flag descriptions).
  - `HookBefore` and `HookAfter` custom hooks which can be used to specify custom behaviour before/after each command
  -  Refactor to use some stateless funcs

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
