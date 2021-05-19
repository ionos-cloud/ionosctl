---
description: Remove IP Failover group from a LAN
---

# IpfailoverRemove

## Usage

```text
ionosctl ipfailover remove [flags]
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
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [NicId,Ip])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --force                  Force command to execute without user input
  -h, --help                   help for remove
      --ip string              Allocated IP (required)
      --lan-id string          The unique LAN Id (required)
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
      --timeout int            Timeout option for Request for IP Failover deletion [seconds] (default 60)
      --wait-for-request       Wait for the Request for IP Failover deletion to be executed
```

## Examples

```text
ionosctl ipfailover remove --datacenter-id 2c08a329-dbe3-427a-8ef9-897e620fef3d --server-id 11c8ac02-224b-4bd0-833c-196719860fc1 --nic-id 5662f39c-b7cb-4840-b6ab-ae43cd0202cc --lan-id 1 --ip "x.x.x.x"
Warning: Are you sure you want to remove ip failover group from lan (y/N) ? 
y
RequestId: 0643462d-22c7-4396-b8e8-dd3c42fce83a
Status: Command ipfailover remove has been successfully executed
```

