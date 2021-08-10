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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Ip] (default [NicId,Ip])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for remove
      --ip string              Allocated IP (required)
      --lan-id string          The unique LAN Id (required)
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for IP Failover deletion [seconds] (default 60)
  -v, --verbose                see step by step process when running a command
  -w, --wait-for-request       Wait for the Request for IP Failover deletion to be executed
```

## Examples

```text
ionosctl ipfailover remove --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id LAN_ID --ip "x.x.x.x"
```

