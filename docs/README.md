---
description: >-
  IonosCTL CLI is currently under development. We are working on adding new
  commands and use-cases in order to support all the operations available in the
  Cloud API.
---

# Introduction

## Overview

IonosCTL is a tool to help you manage your IONOS Cloud resources directly from your terminal.

## Getting started

An IONOS account is required for access to the Cloud API; credentials from your registration are used to authenticate against the IONOS Cloud API.

### Installing `ionosctl`

#### Downloading a Release from GitHub

Check the [Release Page](https://github.com/ionos-cloud/ionosctl/releases) and find the corresponding archive for your operating system and architecture. You can download the archive from your browser or you can follow the next steps if you are using a linux operating system:

```text
# Check if /usr/local/bin is part of your PATH
echo $PATH

# Download and extract the binary (<version> is the full semantic version): 
curl -sL https://github.com/ionos-cloud/ionosctl/releases/download/v<version>/ionosctl-<version>-linux-amd64.tar.gz | tar -xzv

# Move the binary somewhere in your $PATH:
sudo mv ionosctl /usr/local/bin

# Use the ionosctl CLI
ionosctl help
```

For Windows users, you can download the latest release available on [Release Page](https://github.com/ionos-cloud/ionosctl/releases), unzip it and follow this \[official guide\]\([https://msdn.microsoft.com/en-us/library/office/ee537574\(v=office.14\).aspx](https://msdn.microsoft.com/en-us/library/office/ee537574%28v=office.14%29.aspx)\) that explains how to add tools to your `PATH`. The path that you need to add is the path to the folder where you unzipped the ionosctl release.

#### Building a local version\(on a Linux machine\)

If you have a Go environment \(e.g. Go 1.17\) configured, you can build and install the development version of `ionosctl` with:

```text
git clone https://github.com/ionos-cloud/ionosctl.git
```

After cloning the repository, you can build `ionosctl` locally with:

```text
make build
```

To install `ionosctl` locally, you can use:

```text
make install
```

Note that the development version is a work-in-progress of a future stable release and can include bugs. Officially released versions will generally be more stable. Check the latest releases in the [Release Page](https://github.com/ionos-cloud/ionosctl/releases).

Dependencies: `ionosctl` uses [Go Modules](https://github.com/golang/go/wiki/Modules) with vendoring.

On Windows, you just need to run the command

```text
go install
```

from the folder where you cloned the ionosctl git.

### Authenticating with Ionos Cloud

Before using `ionosctl` to perform any operations, you will need to set your credentials for IONOS Cloud account. The authenticating mechanism is first checking the environment variables and if these are not set, it is checking if a configuration file exists and if the user has the right permissions for it.

You can provide your credentials:

* Using environment variables

You can set the environment variables for HTTP basic authentication:

```text
export IONOS_USERNAME="username"
export IONOS_PASSWORD="password"
```

Or you can use token authentication:

```text
export IONOS_TOKEN="token"
```

Also, you can overwrite the api endpoint: `api.ionos.com` via the `--api-url` global flag or via the following environment variable:

```text
export IONOS_API_URL="api-url"
```

* Using `login` command

```text
ionosctl login --user username --password password -v
```

The command can also be used without setting the `--user` and `--password` flags:

```text
ionosctl login
Enter your username:
username
Enter your password:
```

You can also authenticate via `--token` flag exclusively:

```text
ionosctl login --token IONOS_TOKEN
```

After providing credentials, you will be notified if you logged in successfully or not:

```text
Status: Authentication successful!
```

```text
Error: 401 Unauthorized
```

Setting `--api-url` or `IONOS_API_URL` will overwrite the default value of `https://api.ionos.com` for subsequent requests.

After a successful authentication, you will no longer need to provide credentials unless you want to change them. By default, they will be stored in

* macOS: `${HOME}/Library/Application Support/ionosctl/config.json`
* Linux: `${XDG_CONFIG_HOME}/ionosctl/config.json`
* Windows: `%APPDATA%\ionosctl\config.json`

  and retrieved every time you will perform an action on your account.

### Environment Variables

Environment Variable | Description
--- | --- 
`IONOS_USERNAME` | Specify the username used to login, to authenticate against the IONOS Cloud API | 
`IONOS_PASSWORD` | Specify the password used to login, to authenticate against the IONOS Cloud API | 
`IONOS_TOKEN` | Specify the token used to login, if a token is being used instead of username and password |
`IONOS_API_URL` | Specify the API URL. It will overwrite the API endpoint default value `api.ionos.com`. Note: the host URL does not contain the `/cloudapi/v5` path, so it should _not_ be included in the `IONOS_API_URL` environment variable | 

### Enabling Shell Auto-Completion

`ionosctl` provides completions for various shells, for both commands and flags. If you partially type a command or a flag and then press `TAB`, the rest of the command will be automatically filled in.

To enable auto-completion, you need to use `ionosctl completion [shell]`, depending on the shell you are using.

`ionosctl` uses the latest release of Cobra framework, which supports by default completion with descriptions for commands and flags. To disable it, `--no-descriptions` flag is available.

#### Enabling Bash Shell Auto-Completion

To load completions for the current session, execute:

```text
source <(ionosctl completion bash)
```

To make these changes permanent, append the above line to your `.bashrc` file and use:

```text
source ~/.bashrc
```

By default, `TAB` key in Bash is bound to `complete` readline command. If you want to use `menu-complete` append the following line to `.bashrc` file:

```text
bind 'TAB':menu-complete
```

You will need to start a new shell for this setup to take effect.

#### Enabling Fish Shell Auto-Completion

To load completions into the current shell execute:

```text
ionosctl completion fish | source
```

In order to make the completions permanent execute once:

```text
ionosctl completion fish > ~/.config/fish/completions/ionosctl.fish
```

#### Enabling Zsh Shell Auto-Completion

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

#### Enabling PowerShell Auto-Completion

PowerShell supports three different completion modes:

* TabCompleteNext \(default Windows style - on each key press the next option is displayed\)
* Complete \(works like Bash\)
* MenuComplete \(works like Zsh\)

You set the mode with `Set-PSReadLineKeyHandler -Key Tab -Function <mode>`

Descriptions will only be supported for Complete and MenuComplete.

Follow the next steps to enable it:

To load completions for the current session, execute:

```text
PS> ionosctl completion powershell | Out-String | Invoke-Expression
```

To load completions for every new session, run:

```text
PS> ionosctl completion powershell > ionosctl.ps1
```

and source this file from your PowerShell profile or you can append the above line to your PowerShell profile file.

Regarding the PowerShell profile, you can follow the next steps:

* You need to find the PowerShell Profile path using the command `$PROFILE` and verify it is created with `Test-Path $PROFILE`.
* If the result of the previous command is false, the profile doesn’t exist you need to create one, so you can use the command `New-Item -Type File -Force $PROFILE`.
* Now, you created the profile and you can oopen file with a text editor and add the following line: `. $PATH\ionosctl.ps1`, where $PATH is absolute path to ionosctl.ps1 \(for example . D:\ionoscloud\ionosctl.ps1\)

In case you want more details, the profile creating steps are detailed in this link: [https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about\_profiles?view=powershell-7.1](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_profiles?view=powershell-7.1)

* If you get the following error:"path\Microsoft.PowerShell\_profile.ps1" cannot be loaded because running scripts is disabled on this system, you can run the command `Set-ExecutionPolicy RemoteSigned` and restart the terminal.

  After you finish your work with ionosctl, you can run `Set-ExecutionPolicy Restricted` to disable running scripts.

You will need to start a new PowerShell for this setup to take effect.

Note: PowerShell completions require version 5.0 or above, which comes with Windows 10 and can be downloaded separately for Windows 7 or 8.1.

### Output formatting

* Use the `--output` option

You can control the output format with the `--output` or `-o` option. IonosCTL supports JSON format for all commands output by specifying `--output=json`. By default, the output is set to human-readable format.

* Use the `--quiet` option

To redirect all the output to `dev/null`, except for error messages, you can use `--quiet` or `-q` option.

* Use the `--force` option

For deletion/removal commands, you will need to provide a confirmation to perform the action. To force the command to execute without a confirmation, you can use `--force` or `-f` option.

* Use the `--all` option

For deletion/removal commands, you can use the `--all` flag to delete all the resources. The command iterates through all the resources and deletes them. If an error happens, it will be displayed after the entire iteration is done.

* Use the `--cols` option

To obtain only a specific field/column, or a collection of columns on output, you can use the `--cols` option with the list of desired fields.

For example, if you want to print only the Datacenter ID and the Location for your existing Virtual Data Centers, you can use the following command:

```text
ionosctl datacenter list --cols "DatacenterId,Location"
DatacenterId     Location
DATACENTER_ID1   us/ewr
DATACENTER_ID2   us/las
DATACENTER_ID3   us/las
```

Note: When using `TAB` in autocompletion, on `--cols` option on a specific resource, the available columns for that resource will be displayed.

* Use the `--no-headers` option

To skip printing the column headers in output format `text`.

* Use the `--verbose` option

You will see step-by-step process when running a command.

This flag can be used with any command(in general create, read, update, delete, but it's available also for the other specific command) of any resource.

* Use the `--filters` option

You can use the filters option for the list commands, in order to filter the results based on properties or on metadata information. In order to set one or multiple filters, you must use the following format: `--filters KEY1=VALUE1,KEY2=VALUE2`. You can also use the `--max-results` or `--order-by` options.

### Help Information

You can see all available options for each command, use:

```text
ionosctl help [command]

ionosctl help [command] [command]

ionosctl [command] --help

ionosctl [command] -h
```

### Testing

```text
make test
```

### Examples

For each runnable command, use `ionosctl [command] --help`, `ionosctl [command] -h` or `ionosctl help [command]` or see the [full reference documentation](https://github.com/ionos-cloud/ionosctl/tree/decc11ee0cb1b00d1201dd3765219adec53cedc8/docs/subcommands/README.md) to see examples.

### Uninstalling `ionosctl`

#### Local version

To uninstall a local version built with the steps from [Installing Ionosctl](README.md#building-a-local-version), use:

```text
make clean
```

## Feature Reference

The IONOS Cloud CLI aims to offer access to all resources in the IONOS Cloud API and also offers some additional features that make the integration easier:

* authentication for API calls
* handling of asynchronous requests 

## FAQ

* How can I open a bug/feature request?

Bugs & feature requests can be open on the repository issues: [https://github.com/ionos-cloud/ionosctl/issues/new/choose](https://github.com/ionos-cloud/ionosctl/issues/new/choose)

* Can I contribute to the IONOS Cloud CLI?

Sure! Our repository is public, feel free to fork it and file a PR for one of the issues opened in the issues list. We will review it and work together to get it released.

