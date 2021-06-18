---
description: List IP Failovers groups from a LAN
---

# IpfailoverList

## Usage

```text
ionosctl ipfailover list [flags]
```

## Aliases

For `ipfailover` command:

```text
[ipf]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of IP Failovers groups from a LAN.

Required values to run command:

* Data Center Id
* Lan Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Ip] (default [NicId,Ip])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
      --lan-id string          The unique LAN Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl ipfailover list --datacenter-id DATACENTER_ID --lan-id LAN_ID
```

