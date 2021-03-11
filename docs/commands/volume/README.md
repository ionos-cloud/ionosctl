---
description: Volume Operations
---

# Volume

## Usage

```text
ionosctl volume [command]
```

## Aliases

```text
[vol]
```

## Description

The sub-commands of `ionosctl volume` manage your block storage volumes by creating, updating, getting specific information, deleting Volumes or attaching, detaching a Volume from a Server.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for volume
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Enable verbose output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl volume attach](attach/) | Attach a Volume to a Server |
| [ionosctl volume create](create.md) | Create a Volume |
| [ionosctl volume delete](delete.md) | Delete a Volume |
| [ionosctl volume detach](detach.md) | Detach a Volume from a Server |
| [ionosctl volume get](get.md) | Get a Volume |
| [ionosctl volume list](list.md) | List Volumes |
| [ionosctl volume update](update.md) | Update a Volume |

