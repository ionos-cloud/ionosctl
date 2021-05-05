---
description: Delete a Load Balancer
---

# LoadbalancerDelete

## Usage

```text
ionosctl loadbalancer delete [flags]
```

## Description

Use this command to delete the specified Load Balancer.

You can wait for the action to be executed using `--wait` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [LoadBalancerId,Name,Dhcp])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
      --force                    Force command to execute without user input
  -h, --help                     help for delete
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --timeout int              Timeout option for Load Balancer to be deleted [seconds] (default 60)
      --wait                     Wait for Load Balancer to be deleted
```

## Examples

```text
ionosctl loadbalancer delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id 3f9f14a9-5fa8-4786-ba86-a91f9daded2c --force --wait 
Waiting for request: 29c4e7bb-8ce8-4153-8b42-3734d8ede034
RequestId: 29c4e7bb-8ce8-4153-8b42-3734d8ede034
Status: Command loadbalancer delete and request have been successfully executed
```

