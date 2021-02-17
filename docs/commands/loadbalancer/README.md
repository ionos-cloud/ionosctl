---
description: Load Balancer Operations
---

# Loadbalancer

## Usage

```text
ionosctl loadbalancer [command]
```

## Description

The sub-commands of `ionosctl loadbalancer` manage your Load Balancers on your account.
With the `ionosctl loadbalancer` command, you can list, create, delete Load Balancers and manage their configuration details.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for loadbalancer
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Enable verbose output
```

## See also

* [ionosctl](../)
* [ionosctl loadbalancer create](create.md)
* [ionosctl loadbalancer delete](delete.md)
* [ionosctl loadbalancer get](get.md)
* [ionosctl loadbalancer list](list.md)
* [ionosctl loadbalancer update](update.md)

