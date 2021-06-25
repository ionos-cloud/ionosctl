---
description: Authentication command for SDK
---

# Login

## Usage

```text
ionosctl login [flags]
```

## Aliases

For `login` command:

```text
[log auth]
```

## Description

Use this command to authenticate. By default, user data will be saved in:

* macOS: `${HOME}/Library/Application Support/ionosctl/config.json`
* Linux: `${XDG_CONFIG_HOME}/ionosctl/config.json`
* Windows: `%APPDATA%\ionosctl\config.json`.

You can use another configuration file for authentication with `--config` global option.

Note: The command can also be used without `--user` and `--password` flags. For more details, see Examples.

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -h, --help              help for login
  -o, --output string     Desired output format [text|json] (default "text")
  -p, --password string   Password to authenticate
  -q, --quiet             Quiet output
      --token string      Token to authenticate
      --user string       Username to authenticate
```

## Examples

```text
ionosctl login --user USERNAME --password PASSWORD
Status: Authentication successful!

ionosctl login 
Enter your username:
USERNAME
Enter your password:

Status: Authentication successful!
```

