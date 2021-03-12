---
description: Server Operations
---

# Ionosctl Server

## Usage

```text
ionosctl server [command]
```

## Aliases

```text
[svr]
```

## Description

The sub-commands of `ionosctl server` allow you to create, list, get, update, delete, start, stop, reboot Servers.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for server
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl server create](create.md) | Create a Server |
| [ionosctl server delete](delete.md) | Delete a Server |
| [ionosctl server get](get.md) | Get a Server |
| [ionosctl server list](list.md) | List Servers |
| [ionosctl server reboot](reboot.md) | Force a hard reboot of a Server |
| [ionosctl server start](start.md) | Start a Server |
| [ionosctl server stop](stop.md) | Stop a Server |
| [ionosctl server update](update.md) | Update a Server |

