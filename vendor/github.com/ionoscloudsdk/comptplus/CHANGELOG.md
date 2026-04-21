# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.2] - 2026-04-21

### Fixed
* `findSuggestions`, `getCurrentFlagAndValueContext`, and `getFlagValueSuggestions` now use `shellquote.Split` instead of `strings.Fields`, matching the parser change from v1.1.0. Previously, suggestions would break when typing quoted arguments (e.g. `--name "John Oliver" --env `).
* stale doc comments: `InArgsParser` referenced `strings.Fields` as default (now `shellquote.Split`), `GoPromptOptions` referenced `c-bata/go-prompt` (now `elk-language/go-prompt`).

### Added
* comprehensive test suite covering: shell-quote parsing in suggestions, command execution with `SetArgs`, hook lifecycle (before/after, abort, resolved command), flag reset (string, slice, persist, custom behaviour), error handling, dynamic suggestions with `ValidArgsFunction`, hidden commands/flags, bool flag exclusion, flag value descriptions, inherited flags, suggestion filters, async cache concurrency, cache invalidation, deep command trees, command aliases, and context propagation.

## [1.1.1] - 2026-04-16

### Added
* `AsyncFlagValueSuggestions` option for non-blocking flag value completions. When enabled, flag value suggestions are fetched in a background goroutine instead of blocking the prompt.

## [1.1.0] - 2026-04-16

### Changed
* **BREAKING**: `DynamicSuggestionsFunc` signature now receives the resolved `*cobra.Command` as its first parameter, enabling use of `cmd.ValidArgsFunction()` for generic auto-completion logic.
* command execution now uses `cobra.Command.SetArgs()` instead of manipulating the global `os.Args`, avoiding side effects and fixing usage as a sub-command of a larger CLI.
* default argument parser now uses [go-shellquote](https://github.com/kballard/go-shellquote) instead of `strings.Fields`, so quoted arguments (e.g. `--name "John Oliver"`) are parsed correctly. Falls back to `strings.Fields` on parse errors. Custom parsers via `InArgsParser` are still supported.
* `handleUserError` no longer calls `os.Exit(1)` when no `OnErrorFunc` is set. Errors are printed and control returns to the prompt loop.

### Updated
* `elk-language/go-prompt` v1.1.5 → v1.3.1
* `spf13/cobra` v1.8.0 → v1.10.2
* `spf13/pflag` v1.0.5 → v1.0.10
* all indirect dependencies bumped to latest

### Fixed
* tests updated for go-prompt API change (`InsertText` → `InsertTextMoveCursor`)

## [1.0.3] - 2023-11-24

### Changed
* changed `HookBefore` and `HookAfter` to have access to the command object that is about to be / has been executed.


## [1.0.2] - 2023-11-24

### Fixed
* fixed flag reset behaviour for slice/array flags by using casting to `SliceValue` and using `SliceValue.Reset()` by @avirtopeanu-ionos in https://github.com/ionoscloudsdk/comptplus/pull/3
  * In v1.0.1 and before, including in the original cobra-prompt repository, the defaults would be appended to the values of the previous execution

## [1.0.1] - 2023-11-23

### Added
*  added the option to set custom flag reset behaviours by @avirtopeanu-ionos in https://github.com/ionoscloudsdk/comptplus/pull/2

## [1.0.0] - 2023-11-22

### Added
* added completions for flag values. by @avirtopeanu-ionos in https://github.com/ionoscloudsdk/comptplus/pull/1
    * default cache duration for responses set to 500ms - to prevent laggy user interaction
    * support flag descriptions by splitting on `\t`
* added `HookBefore` and `HookAfter` for additional actions before and after command execution.

### Changed
* `PersistFlagValues` behavior:
  * instead of adding a flag, setting PersistFlagValues to true will directly influence persistance throughout the entire shell session.
  * instead of resetting flags to their default value every time a new character is typed, flag defaults are set after a command execution.

## [0.5.0] - 2023-01-28

### Added

- `RunContext` - option to pass context into nested command execututions. ([#9](https://github.com/stromland/cobra-prompt/pull/9) by [@klowdo](https://github.com/klowdo))

## [0.4.0] - 2022-10-04

### Added

- `SuggestionFilter` to `CobraPrompt`. Function to decide which suggestions that should be presentet to the user. Overrides the current filter from go-prompt. ([#8](https://github.com/stromland/cobra-prompt/pull/8) by [@klowdo](https://github.com/klowdo))

## [0.3.0] - 2022-04-25

### Added

- `InArgsParser` to `CobraPrompt`. This makes it possible to decide how arguments should be structured before passing them to Cobra. ([#7](https://github.com/stromland/cobra-prompt/pull/7) by [@klowdo](https://github.com/klowdo))
