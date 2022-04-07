---
description: Get a Network Load Balancer
---

# NetworkloadbalancerGet

## Usage

```text
ionosctl networkloadbalancer get [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Network Load Balancer from a Virtual Data Center. You can also wait for Network Load Balancer to get in AVAILABLE state using `--wait-for-state` option.

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
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
  -i, --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      When using text output, don't print headers
  -o, --output string                   Desired output format [text|json] (default "text")
  -q, --quiet                           Quiet output
  -t, --timeout int                     Timeout option for waiting for Network Load Balancer to be in AVAILABLE state [seconds] (default 300)
  -v, --verbose                         Print step-by-step process when running command
  -W, --wait-for-state                  Wait for specified Network Load Balancer to be in AVAILABLE state
```

## Examples

```text
ionosctl networkloadbalancer get --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID
```

