---
description: Get a Application Load Balancer
---

# ApplicationloadbalancerGet

## Usage

```text
ionosctl applicationloadbalancer get [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Application Load Balancer from a Virtual Data Center. You can also wait for Application Load Balancer to get in AVAILABLE state using `--wait-for-state` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id

## Options

```text
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
  -i, --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan LbPrivateIps State] (default [ApplicationLoadBalancerId,Name,ListenerLan,Ips,TargetLan,LbPrivateIps,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -f, --force                               Force command to execute without user input
  -h, --help                                help for get
  -o, --output string                       Desired output format [text|json] (default "text")
  -q, --quiet                               Quiet output
  -t, --timeout int                         Timeout option for waiting for Application Load Balancer to be in AVAILABLE state [seconds] (default 300)
  -v, --verbose                             see step by step process when running a command
  -W, --wait-for-state                      Wait for specified Application Load Balancer to be in AVAILABLE state
```

## Examples

```text
ionosctl applicationloadbalancer get --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID
```

