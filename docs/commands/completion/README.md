---
description: Generate code to enable auto-completion with `TAB` key
---

# Completion

## Usage

```text
ionosctl completion [command]
```

## Description

Use this command to generate completion code for specific shell for `ionosctl` commands and flags.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for completion
      --force            Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Enable verbose output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl completion bash](bash.md) | Generate code to enable auto-completion with `TAB` key for BASH terminal |
| [ionosctl completion fish](fish.md) | Generate code to enable auto-completion with `TAB` key for Fish terminal |
| [ionosctl completion powershell](powershell.md) | Generate code to enable auto-completion with `TAB` key for PowerShell terminal |
| [ionosctl completion zsh](zsh.md) | Generate code to enable auto-completion with `TAB` key for ZSH terminal |

