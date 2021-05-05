---
description: Generate code to enable auto-completion with `TAB` key for ZSH terminal
---

# Zsh

## Usage

```text
ionosctl completion zsh [flags]
```

## Description

Use this command to generate completion code for ZSH terminal. IonosCTL supports completion for commands and flags.

If shell completions are not already enabled for your environment, you need to enable them. Add the following line to your `~/.zshrc` file:

```text
autoload -Uz compinit; compinit
```

To load completions for each session execute the following commands:

```text
mkdir -p ~/.config/ionosctl/completion/zsh
ionosctl completion zsh > ~/.config/ionosctl/completion/zsh/_ionosctl
```

Finally add the following line to your `~/.zshrc`file, _before_ you call the `compinit` function:

```text
fpath+=(~/.config/ionosctl/completion/zsh)
```

In the end your `~/.zshrc` file should contain the following two lines in the order given here:

```text
fpath+=(~/.config/ionosctl/completion/zsh)
#  ... anything else that needs to be done before compinit
autoload -Uz compinit; compinit
# ...
```

You will need to start a new shell for this setup to take effect. Note: ZSH completions require zsh 5.2 or newer.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for zsh
      --force            Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

