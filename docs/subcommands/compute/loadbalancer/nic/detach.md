---
description: Detach a NIC from a Load Balancer
---

# LoadbalancerNicDetach

## Usage

```text
ionosctl loadbalancer nic detach [flags]
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

For `detach` command:

```text
[d]
```

## Description

Use this command to remove the association of a NIC with a Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id

## Options

```text
  -a, --all                      Detach all Nics.
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [NicId Name Dhcp LanId Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac] (default [NicId,Name,Dhcp,LanId,Ips,State])
      --datacenter-id string     The unique Data Center Id (required)
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -i, --nic-id string            The unique NIC Id (required)
  -t, --timeout int              Timeout option for Request for NIC detachment [seconds] (default 60)
  -w, --wait-for-request         Wait for the Request for NIC detachment to be executed
```

## Examples

```text
ionosctl loadbalancer nic detach --datacenter-id DATACENTER_ID--loadbalancer-id LOADBALANCER_ID --nic-id NIC_ID
```

