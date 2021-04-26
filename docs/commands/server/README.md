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
      --cols strings           Columns to be printed in the standard output (default [ServerId,Name,AvailabilityZone,State,Cores,Ram,CpuFamily])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
      --force                  Force command to execute without user input
  -h, --help                   help for server
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl server add-label](add-label.md) | Add a Label on a Server |
| [ionosctl server attach-volume](attach-volume.md) | Attach a Volume to a Server |
| [ionosctl server create](create.md) | Create a Server |
| [ionosctl server delete](delete.md) | Delete a Server |
| [ionosctl server detach-volume](detach-volume.md) | Detach a Volume from a Server |
| [ionosctl server get](get.md) | Get a Server |
| [ionosctl server get-label](get-label.md) | Get a Label from a Server |
| [ionosctl server get-volume](get-volume.md) | Get an attached Volume from a Server |
| [ionosctl server list](list.md) | List Servers |
| [ionosctl server list-labels](list-labels.md) | List Labels from a Server |
| [ionosctl server list-volumes](list-volumes.md) | List attached Volumes from a Server |
| [ionosctl server reboot](reboot.md) | Force a hard reboot of a Server |
| [ionosctl server remove-label](remove-label.md) | Remove a Label from a Server |
| [ionosctl server start](start.md) | Start a Server |
| [ionosctl server stop](stop.md) | Stop a Server |
| [ionosctl server update](update.md) | Update a Server |

