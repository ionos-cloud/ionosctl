---
description: List attached NICs from a Load Balancer
---

# LoadbalancerNicList

## Usage

```text
ionosctl loadbalancer nic list [flags]
```

## Aliases

For `loadbalancer` command:
```text
[lb]
```

For `nic` command:
```text
[n]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to get a list of attached NICs to a Load Balancer from a Data Center.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [NicId Name Dhcp LanId Ips State FirewallActive Mac] (default [NicId,Name,Dhcp,LanId,Ips,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     help for list
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
```

## Examples

```text
ionosctl loadbalancer nic list --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID
```

