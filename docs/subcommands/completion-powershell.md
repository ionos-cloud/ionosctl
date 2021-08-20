---
description: Generate code to enable auto-completion with `TAB` key for PowerShell terminal
---

# CompletionPowershell

## Usage

```text
ionosctl completion powershell [flags]
```

## Description

Generate the autocompletion script for powershell.

To load completions in your current shell session:

```text
PS C:\> ionosctl completion powershell | Out-String | Invoke-Expression
```

To load completions for every new session, add the output of the above command to your powershell profile.

## Options

```text
      --no-descriptions  disable completion descriptions
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for powershell
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          see step by step process when running a command
```

