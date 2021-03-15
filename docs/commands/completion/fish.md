---
description: Generate code to enable auto-completion with `TAB` key for Fish terminal
---

# Fish

## Usage

```text
ionosctl completion fish [flags]
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
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for fish
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

