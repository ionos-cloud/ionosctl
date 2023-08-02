---
description: "Update a Network Load Balancer"
---

# NetworkloadbalancerUpdate

## Usage

```text
ionosctl networkloadbalancer update [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified Network Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id

## Options

```text
  -u, --api-url string                  Override default host url (default "https://api.ionos.com")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [NetworkLoadBalancerId Name ListenerLan Ips TargetLan LbPrivateIps State] (default [NetworkLoadBalancerId,Name,ListenerLan,Ips,TargetLan,LbPrivateIps,State])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --ips strings                     Collection of IP addresses of the Network Load Balancer
      --listener-lan int                Id of the listening LAN (default 2)
  -n, --name string                     Name of the Network Load Balancer (default "Network Load Balancer")
  -i, --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -o, --output string                   Desired output format [text|json] (default "text")
      --private-ips strings             Collection of private IP addresses with subnet mask of the Network Load Balancer
  -q, --quiet                           Quiet output
      --target-lan int                  Id of the balanced private target LAN (default 1)
  -t, --timeout int                     Timeout option for Request for Network Load Balancer update [seconds] (default 300)
  -v, --verbose                         Print step-by-step process when running command
  -w, --wait-for-request                Wait for the Request for Network Load Balancer update to be executed
```

## Examples

```text
ionosctl networkloadbalancer update --datacenter-id DATACENTER_ID -i NETWORKLOADBALANCER_ID --name NAME
```

