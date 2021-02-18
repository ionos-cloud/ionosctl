---
description: Generate code to enable auto-completion with `TAB` key for ZSH terminal
---

# Zsh

## Usage

```text
ionosctl completion zsh [flags]
```

## Description

Use this command to generate completion code for ZSH terminal. IonosCTL supports completion for commands and flags.
Add the following line to your .profile or .bashrc.

`source  <(ionosctl completion zsh)`

Note:
- ZSH completions require zsh 5.2 or newer.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
  -h, --help             help for zsh
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Enable verbose output
```

## See also

* [ionosctl completion](./)

