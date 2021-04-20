---
description: Resource Operations
---

# Resource

## Usage

```text
ionosctl resource [command]
```

## Description

The sub-command of `ionosctl resource` allows you to list, get Resources.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [ResourceId,Name,SecAuthProtection,Type])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for resource
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl resource get](get.md) | Get all Resources of a Type or a specific Resource Type |
| [ionosctl resource list](list.md) | List Resources |

