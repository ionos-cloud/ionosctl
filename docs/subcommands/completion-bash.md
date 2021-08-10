---
description: Generate code to enable auto-completion with `TAB` key for BASH terminal
---

# CompletionBash

## Usage

```text
ionosctl completion bash [flags]
```

## Aliases

For `completion` command:

```text
[comp]
```

For `bash` command:

```text
[b]
```

## Description

Use this command to generate completion code for BASH terminal. IonosCTL supports completion for commands and flags.

To load completions for the current session, execute:

```text
source <(ionosctl completion bash)
```

To make these changes permanent, append the above line to your `.bashrc` file and use:

```text
source ~/.bashrc
```

You will need to start a new shell for this setup to take effect.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for bash
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          see step by step process when running a command
```

