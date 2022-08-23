[![CI](https://github.com/ionos-cloud/ionosctl/workflows/CI/badge.svg)](https://github.com/ionos-cloud/ionosctl/actions)
[![Gitter](https://img.shields.io/gitter/room/ionos-cloud/sdk-general)](https://gitter.im/ionos-cloud/sdk-general)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=cli-ionosctl&metric=alert_status)](https://sonarcloud.io/dashboard?id=cli-ionosctl)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=cli-ionosctl&metric=bugs)](https://sonarcloud.io/dashboard?id=cli-ionosctl)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=cli-ionosctl&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=cli-ionosctl)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=cli-ionosctl&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=cli-ionosctl)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=cli-ionosctl&metric=security_rating)](https://sonarcloud.io/dashboard?id=cli-ionosctl)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=cli-ionosctl&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=cli-ionosctl)
[![Release](https://img.shields.io/github/v/release/ionos-cloud/ionosctl.svg)](https://github.com/ionos-cloud/ionosctl/releases/latest)
[![Release Date](https://img.shields.io/github/release-date/ionos-cloud/ionosctl.svg)](https://github.com/ionos-cloud/ionosctl/releases/latest)
[![Go](https://img.shields.io/github/go-mod/go-version/ionos-cloud/ionosctl.svg)](https://github.com/ionos-cloud/ionosctl)

![Alt text](.github/IONOS.CLOUD.BLU.svg?raw=true "Title")

# IONOSCTL CLI

## Overview

IonosCTL is a tool to help you manage your Ionos Cloud resources directly from your terminal. 

IonosCTL uses [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) libraries in order to manage commands and options.
Cobra is both a library for creating powerful modern command-line interface (CLI) applications as well as a program to generate applications and command files and it is used in many Go projects together with Viper library. 

[![overview](.github/images/overview.gif)](.github/images/overview.gif)

## Getting started

Before you begin you will need to have signed-up for a [Ionos Cloud](https://www.ionos.com/enterprise-cloud/signup) account. The credentials you establish during sign-up will be used to authenticate against the [Ionos Cloud API](https://dcd.ionos.com/latest/).

### Installing `ionosctl`

#### Installing on Linux

You can install ionosctl using snap package manager:

```
snap install ionosctl
```

### Installing on macOS

You can install `ionosctl` using the [Homebrew](https://brew.sh) package manager:

```bash
brew tap ionos-cloud/homebrew-ionos-cloud
brew install ionosctl
```

### Installing on Windows

You can install `ionosctl` using the [Scoop](https://scoop.sh/) package manager:

```bash
scoop bucket add ionos-cloud https://github.com/ionos-cloud/scoop-bucket.git
scoop install ionos-cloud/ionosctl
```

#### Downloading a Release from Github

Check the [Release Page](https://github.com/ionos-cloud/ionosctl/releases) and find the corresponding archive for your operating system and architecture. You can download the archive from your browser or you can follow the next steps:

```
# Check if /usr/local/bin is part of your PATH
echo $PATH

# Download and extract the binary (<version> is the full semantic version): 
curl -sL https://github.com/ionos-cloud/ionosctl/releases/download/v<version>/ionosctl-<version>-linux-amd64.tar.gz | tar -xzv

# Move the binary somewhere in your $PATH:
sudo mv ionosctl /usr/local/bin

# Use the ionosctl CLI
ionosctl help
```

For Windows users, you can download the latest release available on [Release Page](https://github.com/ionos-cloud/ionosctl/releases), unzip it and follow this [official guide](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) that explains how to add tools to your `PATH`.
The path that you need to add is the path to the folder where you unzipped the ionosctl release.

#### Building a local version(on a Linux machine)

If you have a Go environment (e.g. Go 1.17) configured, you can build and install the development version of `ionosctl` with:

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

Note that the development version is a work-in-progress of a future stable release and can include bugs. Officially released versions will generally be more stable. Check the latest releases in the [Release Page](https://github.com/ionos-cloud/ionosctl/releases).

Dependencies: `ionosctl` uses [Go Modules](https://github.com/golang/go/wiki/Modules) with vendoring.

On Windows, you just need to run the command
```text
go install
```
from the folder where you cloned the ionosctl git.

### Authenticating with Ionos Cloud

Before using `ionosctl` to perform any operations, you will need to set your credentials for IONOS Cloud account. The authentication mechanism is first checking the environment variables and if these are not set, it is checking if a configuration file exists and if the user has the right permissions for it.

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

| Environment Variable | Description                                                                                                                                                                                                                    |
|----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `IONOS_USERNAME`     | Specify the username used to login, to authenticate against the IONOS Cloud API                                                                                                                                                | 
| `IONOS_PASSWORD`     | Specify the password used to login, to authenticate against the IONOS Cloud API                                                                                                                                                | 
| `IONOS_TOKEN`        | Specify the token used to login, if a token is being used instead of username and password                                                                                                                                     |
| `IONOS_API_URL`      | Specify the API URL. It will overwrite the API endpoint default value `api.ionos.com`. Note: the host URL does not contain the `/cloudapi/v5` path, so it should _not_ be included in the `IONOS_API_URL` environment variable | 
| `IONOS_PINNED_CERT`  | Specify the SHA-256 public fingerprint here, enables certificate pinning                                                                                                                                                       |

### Certificate pinning:

You can enable certificate pinning if you want to bypass the normal certificate checking procedure,
by doing the following:

Set env variable IONOS_PINNED_CERT=<insert_sha256_public_fingerprint_here>

You can get the sha256 fingerprint most easily from the browser by inspecting the certificate.

### Enabling Shell Auto-Completion

`ionosctl` provides completions for various shells, for both commands and flags. If you partially type a command or a flag and then press `TAB`, the rest of the command will be automatically filled in. 

To enable auto-completion, you need to use `ionosctl completion [shell]`, depending on the shell you are using.

`ionosctl` uses the latest release of Cobra framework, which supports by default completion with descriptions for commands and flags. To disable it, `--no-descriptions` flag is available.

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

Regarding the PowerShell profile, you can follow the next steps:

* You need to find the PowerShell Profile path using the command ```$PROFILE``` and verify it is created with  ```Test-Path $PROFILE```.

* If the result of the previous command is false, the profile doesn’t exist you need to create one, so you can use the command ```New-Item -Type File -Force $PROFILE```.

* Now, you created the profile and you can oopen file with a text editor and add the following line: ```. $PATH\ionosctl.ps1```, where $PATH is absolute path to ionosctl.ps1 (for example . D:\ionoscloud\ionosctl.ps1) 

In case you want more details, the profile creating steps are detailed in this link: https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_profiles?view=powershell-7.1

* If you get the following error:"path\Microsoft.PowerShell_profile.ps1" cannot be loaded because running scripts is disabled on this system, you can run the command ```Set-ExecutionPolicy RemoteSigned``` and restart the terminal.
After you finish your work with ionosctl, you can run ```Set-ExecutionPolicy Restricted``` to disable running scripts.

You will need to start a new PowerShell for this setup to take effect.

Note: PowerShell completions require version 5.0 or above, which comes with Windows 10 and can be downloaded separately for Windows 7 or 8.1. 

### Output formatting

* Use the `--output` option

You can control the output format with the `--output` or `-o` option. IonosCTL supports JSON format for all commands output by specifying `--output=json`. By default, the output is set to `text`.

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

You can use the filters option for the majority of list commands, in order to filter the results based on properties or on metadata information. In order to set one or multiple filters, you must use the following format: `--filters KEY1=VALUE1,KEY2=VALUE2`. You can also use the `--max-results` or `--order-by` options.

### Help Information

You can see all available options for each command, use:

```text
ionosctl help [command]

ionosctl help [command] [command]

ionosctl [command] --help

ionosctl [command] -h
```

### Testing

#### What Are We Testing?

The purpose of our unit tests is to ensure that properties set via flags are handled as expected before sending API Requests. The tests are integrated into [GitHub Actions](https://github.com/ionos-cloud/ionosctl/actions) that run at every PR, commit and release.

We understand the importance of testing, and we put our best efforts to add integration tests as well.

#### How to Run Tests Locally

In order to run the tests locally, you can simply run:

```text
make test
```

### Examples

For each runnable command, use `ionosctl [command] --help`, `ionosctl [command] -h`  or `ionosctl help [command]` or see the [full reference documentation](docs/subcommands) to see examples.

### Uninstalling `ionosctl` 

#### Local version

To uninstall a local version built with the steps from [Installing Ionosctl](#building-a-local-versionon-a-linux-machine), use:
```text
make clean
``` 

## Contributing

Bugs & feature requests can be open on the repository issues: https://github.com/ionos-cloud/ionosctl/issues/new/choose

- Can I contribute to IonosCTL?

Sure! Our repository is public, feel free to fork it and file a PR for one of the issues opened in the issues list. We will review it and work together to get it released.
