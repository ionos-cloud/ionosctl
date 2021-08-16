---
description: Generate code to enable auto-completion with `TAB` key for BASH terminal
---

# CompletionBash

## Usage

```text
ionosctl completion bash [flags]
```

## Description

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package. If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

```text
$ source <(ionosctl completion bash)
```

To load completions for every new session, execute once:

* Linux:

```text
$ ionosctl completion bash > /etc/bash_completion.d/ionosctl
```

* MacOS:

```text
$ ionosctl completion bash > /usr/local/etc/bash_completion.d/ionosctl
```

You will need to start a new shell for this setup to take effect.

## Options

```text
      --no-descriptions  disable completion descriptions
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for bash
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          see step by step process when running a command
```

