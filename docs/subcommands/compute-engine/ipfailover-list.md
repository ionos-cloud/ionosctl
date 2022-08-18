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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Ip] (default [NicId,Ip])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --lan-id string          The unique LAN Id (required)
      --no-headers             When using text output, don't print headers
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl ipfailover list --datacenter-id DATACENTER_ID --lan-id LAN_ID
```

