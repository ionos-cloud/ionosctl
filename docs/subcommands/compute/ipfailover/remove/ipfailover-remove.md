---
description: Remove IP Failover group from a LAN
---

# IpfailoverRemove

## Usage

```text
ionosctl ipfailover remove [flags]
```

## Aliases

For `ipfailover` command:

```text
[ipf]
```

For `remove` command:

```text
[r]
```

## Description

Use this command to remove an IP Failover group from a LAN.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address

## Options

```text
  -a, --all                    Remove all IP Failovers.
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --ip ip                  Allocated IP (required)
      --lan-id string          The unique LAN Id (required)
      --nic-id string          The unique NIC Id (required)
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for IP Failover deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for IP Failover deletion to be executed
```

## Examples

```text
ionosctl ipfailover remove --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id LAN_ID --ip "x.x.x.x"
```

