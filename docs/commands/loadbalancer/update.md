---
description: Update a Load Balancer
---

# Update

## Usage

```text
ionosctl loadbalancer update [flags]
```

## Description

Use this command to update the configuration of a specified Load Balancer.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -u, --api-url string             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings               Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string       The unique Data Center Id [Required flag]
  -h, --help                       help for update
      --ignore-stdin               Force command to execute without user input
      --loadbalancer-dhcp          Indicates if the Load Balancer will reserve an IP using DHCP (default true)
      --loadbalancer-id string     The unique Load Balancer Id [Required flag]
      --loadbalancer-ip string     The IP of the Load Balancer
      --loadbalancer-name string   Name of the Load Balancer
  -o, --output string              Desired output format [text|json] (default "text")
  -q, --quiet                      Quiet output
      --timeout int                Timeout option for Load Balancer to be updated [seconds] (default 60)
      --wait                       Wait for Load Balancer to be updated
```

## Examples

```text
ionosctl loadbalancer update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id 3f9f14a9-5fa8-4786-ba86-a91f9daded2c --loadbalancer-dhcp=false --wait
Waiting for request: 0a9279d8-9757-41e0-b64f-b4cd2baf4717
LoadbalancerId                         Name               Dhcp
3f9f14a9-5fa8-4786-ba86-a91f9daded2c   demoLoadBalancer   false
RequestId: 0a9279d8-9757-41e0-b64f-b4cd2baf4717
Status: Command loadbalancer update and request have been successfully executed
```

