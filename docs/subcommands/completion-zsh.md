---
description: Generate code to enable auto-completion with `TAB` key for ZSH terminal
---

# CompletionZsh

## Usage

```text
ionosctl completion zsh [flags]
```

## Description

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need to enable it. You can execute the following once:

```text
$ echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions for every new session, execute once:

* Linux:

  ```text
  $ ionosctl completion zsh > "${fpath[1]}/_ionosctl"
  ```

* MacOS:

```text
$ ionosctl completion zsh > /usr/local/share/zsh/site-functions/_ionosctl
```

You will need to start a new shell for this setup to take effect.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for zsh
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          see step by step process when running a command
```

