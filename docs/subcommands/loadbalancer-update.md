---
description: Update a Load Balancer
---

# LoadbalancerUpdate

## Usage

```text
ionosctl loadbalancer update [flags]
```

## Aliases

For `loadbalancer` command:
```text
[lb]
```

## Description

Use this command to update the configuration of a specified Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option.

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
      --dhcp                     Indicates if the Load Balancer will reserve an IP using DHCP (default true)
  -f, --force                    Force command to execute without user input
  -h, --help                     help for update
      --ip string                The IP of the Load Balancer
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -n, --name string              Name of the Load Balancer
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
  -t, --timeout int              Timeout option for Request for Load Balancer update [seconds] (default 60)
  -w, --wait-for-request         Wait for Request for Load Balancer update to be executed
```

## Examples

```text
ionosctl loadbalancer update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id 3f9f14a9-5fa8-4786-ba86-a91f9daded2c --dhcp=false --wait-for-request
1.2s Waiting for request... DONE
LoadbalancerId                         Name               Dhcp
3f9f14a9-5fa8-4786-ba86-a91f9daded2c   demoLoadBalancer   false
RequestId: 0a9279d8-9757-41e0-b64f-b4cd2baf4717
Status: Command loadbalancer update & wait have been successfully executed
```

