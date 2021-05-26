---
description: Generate code to enable auto-completion with `TAB` key for Fish terminal
---

# CompletionFish

## Usage

```text
ionosctl completion fish [flags]
```

## Aliases

For `completion` command:
```text
[comp]
```

For `fish` command:
```text
[f]
```

## Description

Use this command to generate completion code for Fish terminal. IonosCTL supports completion for commands and flags.

To load completions into the current shell execute:

```text
ionosctl completion fish | source
```

In order to make the completions permanent execute once:

```text
ionosctl completion fish > ~/.config/fish/completions/ionosctl.fish
```

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for fish
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

