---
description: Server Operations
---

# Server

## Usage

```text
ionosctl server [command]
```

## Description

The sub-commands of `ionosctl server` allow you to create, list, get, update, delete, start, stop, reboot Servers.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for server
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Enable verbose output
```

## See also

* [ionosctl](../)
* [ionosctl server create](create.md)
* [ionosctl server delete](delete.md)
* [ionosctl server get](get.md)
* [ionosctl server list](list.md)
* [ionosctl server reboot](reboot.md)
* [ionosctl server start](start.md)
* [ionosctl server stop](stop.md)
* [ionosctl server update](update.md)

