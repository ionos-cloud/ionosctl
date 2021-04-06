---
description: Load Balancer Operations
---

# LoadBalancer

## Usage

```text
ionosctl loadbalancer [command]
```

## Aliases

```text
[lb]
```

## Description

The sub-commands of `ionosctl loadbalancer` manage your Load Balancers on your account. With the `ionosctl loadbalancer` command, you can list, create, delete Load Balancers and manage their configuration details.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
  -h, --help                   help for loadbalancer
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl loadbalancer create](create.md) | Create a Load Balancer |
| [ionosctl loadbalancer delete](delete.md) | Delete a Load Balancer |
| [ionosctl loadbalancer get](get.md) | Get a Load Balancer |
| [ionosctl loadbalancer list](list.md) | List Load Balancers |
| [ionosctl loadbalancer update](update.md) | Update a Load Balancer |

