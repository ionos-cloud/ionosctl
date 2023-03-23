---
description: Attach a NIC to a Load Balancer
---

# LoadbalancerNicAttach

## Usage

```text
ionosctl loadbalancer nic attach [flags]
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

For `attach` command:

```text
[a]
```

## Description

Use this command to associate a NIC to a Load Balancer, enabling the NIC to participate in load-balancing.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

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
      --server-id string         The unique Server Id on which NIC is build on. Not required, but it helps in autocompletion
  -t, --timeout int              Timeout option for Request for NIC attachment [seconds] (default 60)
  -w, --wait-for-request         Wait for the Request for NIC attachment to be executed
```

## Examples

```text
ionosctl loadbalancer nic attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --loadbalancer-id LOADBALANCER_ID
```

