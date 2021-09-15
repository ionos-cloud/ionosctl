---
description: List Application Load Balancers
---

# ApplicationloadbalancerList

## Usage

```text
ionosctl applicationloadbalancer list [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Application Load Balancers from a specified Virtual Data Center.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan LbPrivateIps State] (default [ApplicationLoadBalancerId,Name,ListenerLan,Ips,TargetLan,LbPrivateIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl applicationloadbalancer list --datacenter-id DATACENTER_ID
```

