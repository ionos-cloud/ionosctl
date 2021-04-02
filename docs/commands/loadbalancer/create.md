---
description: Create a Load Balancer
---

# Create

## Usage

```text
ionosctl loadbalancer create [flags]
```

## Description

Use this command to create a new Load Balancer on your account. The name, IP and DHCP for the Load Balancer can be set.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings               Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string       The unique Data Center Id [Required flag]
  -h, --help                       help for create
      --ignore-stdin               Force command to execute without user input
      --loadbalancer-dhcp          Indicates if the Load Balancer will reserve an IP using DHCP (default true)
      --loadbalancer-name string   Name of the Load Balancer
  -o, --output string              Desired output format [text|json] (default "text")
  -q, --quiet                      Quiet output
      --timeout int                Timeout option for Load Balancer to be created [seconds] (default 60)
      --wait                       Wait for Load Balancer to be created
```

## Examples

```text
ionosctl loadbalancer create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-name demoLoadBalancer
LoadbalancerId                         Name               Dhcp
3f9f14a9-5fa8-4786-ba86-a91f9daded2c   demoLoadBalancer   true
RequestId: 74441964-1134-4009-8b81-d7189170885e
Status: Command loadbalancer create has been successfully executed
```

