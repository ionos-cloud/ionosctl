---
description: Generate code to enable auto-completion with `TAB` key for Fish terminal
---

# CompletionFish

## Usage

```text
ionosctl completion fish [flags]
```

## Description

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

```text
$ ionosctl completion fish | source
```

To load completions for every new session, execute once:

```text
$ ionosctl completion fish > ~/.config/fish/completions/ionosctl.fish
```

You will need to start a new shell for this setup to take effect.

## Options

```text
      --no-descriptions  disable completion descriptions
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for fish
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

