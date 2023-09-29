---
description: "Convenience command for deletion of config file credentials. To also remove your account's active tokens, use `ionosctl token delete --all`"
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

Convenience command for deletion of config file credentials. To also remove your account's active tokens, use `ionosctl token delete --all`

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl logout
```

