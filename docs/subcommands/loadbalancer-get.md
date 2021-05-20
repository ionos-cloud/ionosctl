---
description: Get a Load Balancer
---

# LoadbalancerGet

## Usage

```text
ionosctl loadbalancer get [flags]
```

## Aliases

For `loadbalancer` command:
```text
[lb]
```

## Description

Use this command to retrieve information about a Load Balancer instance.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -C, --cols strings             Set of columns to be printed on output 
                                 Available columns: [LoadBalancerId Name Dhcp State Ip] (default [LoadBalancerId,Name,Dhcp,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     help for get
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
```

## Examples

```text
ionosctl loadbalancer get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id 3f9f14a9-5fa8-4786-ba86-a91f9daded2c 
LoadbalancerId                         Name               Dhcp
3f9f14a9-5fa8-4786-ba86-a91f9daded2c   demoLoadBalancer   false
```

