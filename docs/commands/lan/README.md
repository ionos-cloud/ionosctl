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
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for lan
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Enable verbose output
```

## See also

* [ionosctl](../)
* [ionosctl lan create](create.md)
* [ionosctl lan delete](delete.md)
* [ionosctl lan get](get.md)
* [ionosctl lan list](list.md)
* [ionosctl lan update](update.md)

