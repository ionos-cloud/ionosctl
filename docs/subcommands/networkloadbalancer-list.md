---
description: List Network Load Balancers
---

# NetworkloadbalancerList

## Usage

```text
ionosctl networkloadbalancer list [flags]
```

## Aliases

For `networkloadbalancer` command:
```text
[nlb]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to list Network Load Balancers from a specified Virtual Data Center.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NetworkLoadBalancerId Name ListenerLan Ips TargetLan LbPrivateIps State] (default [NetworkLoadBalancerId,Name,ListenerLan,Ips,TargetLan,LbPrivateIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl networkloadbalancer list --datacenter-id DATACENTER_ID
```

