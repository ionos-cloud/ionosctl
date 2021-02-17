---
description: Generate code to enable auto-completion with TAB key for BASH terminal
---

# Bash

## Usage

```text
ionosctl completion bash [flags]
```

## Description

Use this command to generate completion code for BASH terminal. IonosCTL supports completion for commands and flags.

Follow the next steps to enable it:

Linux:

- Generate completion code:
`ionosctl completion bash > $PWD/ionosctl_bash_completion.sh`
- Copy generated code to `/etc/bash_completion.d/`:
`sudo cp $PWD/ionosctl_bash_completion.sh /etc/bash_completion.d/`
- Restart terminal to use auto-completion with TAB key.
- Clean-up:
`rm $PWD/ionosctl_bash_completion.sh`

Mac OS:

`ionosctl generate completion bash > /usr/local/etc/bash_completion.d/ionosctl_bash_completion.sh`

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

