---
description: LAN Operations
---

# Lan

## Usage

```text
ionosctl lan [command]
```

## Description

The sub-commands of `ionosctl lan` allow you to create, list, get, update, delete LANs.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [LanId,Name,Public])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
  -h, --help                   help for lan
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl lan create](create.md) | Create a LAN |
| [ionosctl lan delete](delete.md) | Delete a LAN |
| [ionosctl lan get](get.md) | Get a LAN |
| [ionosctl lan list](list.md) | List LANs |
| [ionosctl lan update](update.md) | Update a LAN |

