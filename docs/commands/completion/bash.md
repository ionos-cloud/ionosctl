---
description: Generate code to enable auto-completion with `TAB` key for BASH terminal
---

# Bash

## Usage

```text
ionosctl completion bash [flags]
```

## Description

Use this command to generate completion code for BASH terminal. IonosCTL supports completion for commands and flags.

To load completions for the current session, execute: 

    source <(ionosctl completion bash)

To make these changes permanent, append the above line to your `.bashrc` file and use:

    source ~/.bashrc

You will need to start a new shell for this setup to take effect.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
  -h, --help             help for bash
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Enable verbose output
```

