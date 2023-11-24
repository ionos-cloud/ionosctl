# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
