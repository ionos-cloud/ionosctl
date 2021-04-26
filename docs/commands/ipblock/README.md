---
description: IpBlock Operations
---

# IpBlock

## Usage

```text
ionosctl ipblock [command]
```

## Description

The sub-commands of `ionosctl ipblock` allow you to create/reserve, list, get, update, delete IpBlocks.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force            Force command to execute without user input
  -h, --help             help for ipblock
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl ipblock add-label](add-label.md) | Add a Label on a IpBlock |
| [ionosctl ipblock create](create.md) | Create/Reserve an IpBlock |
| [ionosctl ipblock delete](delete.md) | Delete an IpBlock |
| [ionosctl ipblock get](get.md) | Get an IpBlock |
| [ionosctl ipblock get-label](get-label.md) | Get a Label from a IpBlock |
| [ionosctl ipblock list](list.md) | List IpBlocks |
| [ionosctl ipblock list-labels](list-labels.md) | List Labels from a IpBlock |
| [ionosctl ipblock remove-label](remove-label.md) | Remove a Label from a IpBlock |
| [ionosctl ipblock update](update.md) | Update an IpBlock |

