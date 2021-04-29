---
description: Network Interfaces Operations
---

# NetworkInterface

## Usage

```text
ionosctl nic [command]
```

## Description

The sub-commands of `ionosctl nic` allow you to create, list, get, update, delete NICs. To attach a NIC to a Load Balancer, use the Load Balancer command `ionosctl loadbalancer attach-nic`.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --force                  Force command to execute without user input
  -h, --help                   help for nic
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl nic create](create.md) | Create a NIC |
| [ionosctl nic delete](delete.md) | Delete a NIC |
| [ionosctl nic get](get.md) | Get a NIC |
| [ionosctl nic list](list.md) | List NICs |
| [ionosctl nic update](update.md) | Update a NIC |

