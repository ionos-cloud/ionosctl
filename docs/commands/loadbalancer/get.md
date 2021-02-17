---
description: Get a Load Balancer
---

# Get

## Usage

```text
ionosctl loadbalancer get [flags]
```

## Description

Use this command to retrieve information about a Load Balancer instance.

Required values to run command:
- Data Center Id
- Load Balancer Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string     The unique Data Center Id
  -h, --help                     help for get
      --ignore-stdin             Force command to execute without user input
      --loadbalancer-id string   The unique Load Balancer Id [Required flag]
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
  -v, --verbose                  Enable verbose output
```

## Examples

```text
ionosctl loadbalancer get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id 3f9f14a9-5fa8-4786-ba86-a91f9daded2c 
LoadbalancerId                         Name               Dhcp
3f9f14a9-5fa8-4786-ba86-a91f9daded2c   demoLoadBalancer   false
```

