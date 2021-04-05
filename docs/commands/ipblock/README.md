---
description: IpBlock Operations
---

# Ipblock

## Usage

```text
ionosctl ipblock [command]
```

## Aliases

```text
[ip]
```

## Description

The sub-commands of `ionosctl ipblock` allow you to create/reserve, list, get, update, delete IpBlocks.

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
| [ionosctl ipblock create](create.md) | Create/Reserve an IpBlock |
| [ionosctl ipblock delete](delete.md) | Delete an IpBlock |
| [ionosctl ipblock get](get.md) | Get an IpBlock |
| [ionosctl ipblock list](list.md) | List IpBlocks |
| [ionosctl ipblock update](update.md) | Update an IpBlock |

