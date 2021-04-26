---
description: IPBlock Operations
---

# IPBlock

## Usage

```text
ionosctl ipblock [command]
```

## Description

The sub-commands of `ionosctl ipblock` allow you to create/reserve, list, get, update, delete IPBlocks.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for ipblock
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl ipblock add-label](add-label.md) | Add a Label on a IPBlock |
| [ionosctl ipblock create](create.md) | Create/Reserve an IPBlock |
| [ionosctl ipblock delete](delete.md) | Delete an IPBlock |
| [ionosctl ipblock get](get.md) | Get an IPBlock |
| [ionosctl ipblock get-label](get-label.md) | Get a Label from a IPBlock |
| [ionosctl ipblock list](list.md) | List IPBlocks |
| [ionosctl ipblock list-labels](list-labels.md) | List Labels from a IPBlock |
| [ionosctl ipblock remove-label](remove-label.md) | Remove a Label from a IPBlock |
| [ionosctl ipblock update](update.md) | Update an IPBlock |

