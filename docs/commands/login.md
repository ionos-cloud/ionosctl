---
description: Authentication command for SDK
---

# Login

## Usage

```text
ionosctl login [flags]
```

## Description

Use this command to authenticate. User data will be saved in `$XDG_CONFIG_HOME/ionosctl-config.json` file. 

You can use another configuration file for authentication with `--config` global option.

Note: The command can also be used without `--user` and `--password` flags (see Examples).

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
  -h, --help              help for login
      --ignore-stdin      Force command to execute without user input
  -o, --output string     Desired output format [text|json] (default "text")
      --password string   Password to authenticate
  -q, --quiet             Quiet output
      --user string       Username to authenticate
  -v, --verbose           Enable verbose output
```

## Examples

```text
ionosctl login --user USERNAME --password PASSWORD
✔ Status: Authentication successful!

ionosctl login 
Enter your username:
USERNAME
Enter your password:

✔ Status: Authentication successful!
```

