---
description: Get an attached NIC to a Load Balancer
---

# LoadbalancerNicGet

## Usage

```text
ionosctl loadbalancer nic get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve the attributes of a given load balanced NIC.

Required values to run the command:

* Data Center Id
* Load Balancer Id
* NIC Id

## Options

```text
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [NicId Name Dhcp LanId Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac] (default [NicId,Name,Dhcp,LanId,Ips,State])
      --datacenter-id string     The unique Data Center Id (required)
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -i, --nic-id string            The unique NIC Id (required)
```

## Examples

```text
ionosctl loadbalancer nic get --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --nic-id NIC_ID
```

