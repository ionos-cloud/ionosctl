---
description: "Convenience command for removing config file credentials"
---

# ConfigLogout

## Usage

```text
ionosctl config logout [flags]
```

## Aliases

For `config` command:

```text
[cfg]
```

## Description

This command is a 'Quality of Life' command which will parse your config file for fields that contain sensitive data.
If any such fields are found, their values will be replaced with an empty string.

AUTHENTICATION ORDER
ionosctl uses a layered approach for authentication, prioritizing sources in this order:
  1. Global flags
  2. Environment variables
  3. Config file entries
Within each layer, a token takes precedence over a username and password combination. For instance, if a token and a username/password pair are both defined in environment variables, ionosctl will prioritize the token. However, higher layers can override the use of a token from a lower layer. For example, username and password environment variables will supersede a token found in the config file.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -t, --timeout int      Timeout in seconds for polling the request (default 60)
  -v, --verbose          Print step-by-step process when running command
  -w, --wait             Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl logout
```

