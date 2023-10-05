---
description: "List attached NICs from a Load Balancer"
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
  -u, --api-url string           Override default host url (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [NicId Name Dhcp LanId Ips IPv6Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac DHCPv6 IPv6CidrBlock] (default [NicId,Name,Dhcp,LanId,Ips,IPv6Ips,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
  -F, --filters strings          cloudapiv6.ArgOrderByDescription. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2.
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -M, --max-results int32        The maximum number of elements to return
      --order-by string          Limits results to those containing a matching value for a specific property
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
  -v, --verbose                  Print step-by-step process when running command
```

## Examples

```text
ionosctl loadbalancer nic list --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID
```

