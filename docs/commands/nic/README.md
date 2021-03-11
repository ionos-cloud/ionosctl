---
description: Network Interfaces Operations
---

# Nic

## Usage

```text
ionosctl nic [command]
```

## Description

The sub-commands of `ionosctl nic` allow you to create, list, get, update, delete NICs or attach, detach a NIC from a Load Balancer.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for nic
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id
  -v, --verbose                Enable verbose output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl nic attach](attach/) | Attach a NIC to a Load Balancer |
| [ionosctl nic create](create.md) | Create a NIC |
| [ionosctl nic delete](delete.md) | Delete a NIC |
| [ionosctl nic detach](detach.md) | Detach a NIC from a Load Balancer |
| [ionosctl nic get](get.md) | Get a NIC |
| [ionosctl nic list](list.md) | List NICs |
| [ionosctl nic update](update.md) | Update a NIC |

