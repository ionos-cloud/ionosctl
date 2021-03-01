# Introduction

## Overview

IonosCTL is a tool to help you manage your Ionos Cloud resources directly from your terminal. IonosCTL uses [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) libraries in order to manage commands and options.
Cobra is both a library for creating powerful modern command-line interface (CLI) applications as well as a program to generate applications and command files and it is used in many Go projects together with Viper library. 

## Getting started

Before you begin you will need to have signed-up for a [Ionos Cloud](https://www.ionos.com/enterprise-cloud/signup) account. The credentials you establish during sign-up will be used to authenticate against the [Ionos Cloud API](https://dcd.ionos.com/latest/).

### Installing `ionosctl`

#### Downloading a Release from Github

Check the [Releases page](https://www.ionos.com/ionosctl/releases) and find the corresponding archive for your operating system and architecture. You can download the archive from your browser or you can follow the next steps:

```
# Check if /usr/local/bin is part of your PATH
echo $PATH

# Download and extract the binary (<version> is the full semantic version): 
curl -sL https://github.com/ionos-cloud/ionosctl/releases/download/v<version>/ionosctl-<version>-linux-amd64.tar.gz | tar -xzv

# Move the binary somewhere in your $PATH:
sudo mv ~/ionosctl /usr/local/bin

# Use the ionosctl CLI
ionosctl help
```

For Windows users, you can download the latest release available on [Releases page](https://www.ionos.com/ionosctl/releases) and follow [this official guide](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) that explains how to add tools to your `PATH`. 

#### Building a local version

If you have a Go environment configured, you can build and install the development version of `ionosctl` with: 
```
git clone https://github.com/ionos-cloud/ionosctl.git
```

After cloning the repository, you can build `ionosctl` locally with:
```
make build
```
To install `ionosctl` locally, you can use: 
```
make install 
```

Note that the development version is a work-in-progress of a future stable release and can include bugs. Officially released versions will generally be more stable.

Dependencies: `ionosctl` uses [Go Modules](https://github.com/golang/go/wiki/Modules) with vendoring.

### Authenticating with Ionos Cloud

Before using `ionosctl` to perform any operations, you will need to set your credentials for Ionos Cloud account: 

```
ionosctl login --user username --password password 
```
The command can also be used without setting the `--user` and `--password` flags:
```
ionosctl login
Enter your username:
username
Enter your password:

```

After providing credentials, you will be notified if you logged in successfully or not:

```
Status: Authentication successful!
```

```
Error: 401 Unauthorized
```

After a successful authentication, you will no longer need to provide credentials unless you want to change them. 
By default, they will be stored in 
- macOS: `${HOME}/Library/Application Support/ionosctl/config.json`
- Linux: `${XDG_CONFIG_HOME}/ionosctl/config.json`
- Windows: `%APPDATA%\ionosctl\config.json`
and retrieved every time you will perform an action on your account.

If you want to use a different configuration file, use `--config` option.

### Enabling Shell Auto-Completion

`ionosctl` provides completions for various shells, for both commands and flags. If you partially type a command or a flag and then press `TAB`, the rest of the command will be automatically filled in. 

To enable auto-completion, you need to use `ionosctl completion [shell]`, depending on the shell you are using.

#### Enabling Bash Shell Auto-Completion

To load completions for the current session, execute: 
```
source <(ionosctl completion bash)
```

To make these changes permanent, append the above line to your `.bashrc` file and use:
```
source ~/.bashrc
```

By default, `TAB` key in Bash is bound to `complete` readline command. 
If you want to use `menu-complete` append the following line to `.bashrc` file:
```
bind 'TAB':menu-complete
```

You will need to start a new shell for this setup to take effect.

#### Enabling Fish Shell Auto-Completion

To load completions into the current shell execute:
```
ionosctl completion fish | source
```

In order to make the completions permanent execute once:
```
ionosctl completion fish > ~/.config/fish/completions/ionosctl.fish
```

#### Enabling Zsh Shell Auto-Completion

If shell completions are not already enabled for your environment, you need to enable them. 
Add the following line to your `~/.zshrc` file:
```
autoload -Uz compinit; compinit
```

To load completions for each session execute the following commands:
```
mkdir -p ~/.config/ionosctl/completion/zsh
ionosctl completion zsh > ~/.config/ionosctl/completion/zsh/_ionosctl
```

Finally add the following line to your `~/.zshrc`file, *before* you
call the `compinit` function:
```
fpath+=(~/.config/ionosctl/completion/zsh)
```

In the end your `~/.zshrc` file should contain the following two lines in the order given here:
```
fpath+=(~/.config/ionosctl/completion/zsh)
#  ... anything else that needs to be done before compinit
autoload -Uz compinit; compinit
# ...
```

You will need to start a new shell for this setup to take effect.
Note: ZSH completions require zsh 5.2 or newer.

#### Enabling PowerShell Auto-Completion

PowerShell supports three different completion modes:

- TabCompleteNext (default Windows style - on each key press the next option is displayed)
- Complete (works like Bash)
- MenuComplete (works like Zsh)

You set the mode with `Set-PSReadLineKeyHandler -Key Tab -Function <mode>`

Descriptions will only be supported for Complete and MenuComplete.

Follow the next steps to enable it:

To load completions for the current session, execute: 
```
PS> ionosctl completion powershell | Out-String | Invoke-Expression
```

To load completions for every new session, run:
```
PS> ionosctl completion powershell > ionosctl.ps1
```

and source this file from your PowerShell profile or you can append the above line to your PowerShell profile file. 

You will need to start a new PowerShell for this setup to take effect.

Note: PowerShell completions require version 5.0 or above, which comes with Windows 10 and can be downloaded separately for Windows 7 or 8.1. 

### Output configuration

You can control the output format with the `--output` option. `ionosctl` supports JSON format for all commands output by specifying `--output=json`.

To redirect all the output to `dev/null`, except for error messages, you can use `--quiet` option. 

For `list` and `get` commands, you can also specify which information should be printed using `--cols` option.

For `delete`,`stop`,`detach` commands, you will need to provide a confirmation to perform the action. To force the command to execute without a confirmation, you can use `--ignore-stdin` flag.

### Testing 

```text
make test
```

### Examples

For more information about each available command, including examples, use `ionosctl [command] --help` or `ionosctl help [command]` or see the [full reference documentation](./commands/README.md). 

### Uninstalling `ionosctl` 

#### Local version

To uninstall a local version built with the steps from [Installing Ionosctl](#building-a-local-version), use:
```text
make clean
```

## Feature Reference 

The IONOS Cloud CLI aims to offer access to all resources in the IONOS Cloud API and also offers some additional features that make the integration easier: 
- authentication for API calls
- handling of asynchronous requests 

## FAQ
- How can I open a bug/feature request?

Bugs & feature requests can be open on the repository issues: https://github.com/ionos-cloud/ionosctl/issues/new/choose
