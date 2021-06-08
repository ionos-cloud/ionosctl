---
description: List NICs
---

# NicList

## Usage

```text
ionosctl nic list [flags]
```

## Aliases

For `nic` command:
```text
[n]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to get a list of NICs on your account.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Name Dhcp LanId Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac] (default [NicId,Name,Dhcp,LanId,Ips,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id
```

## Examples

```text
ionosctl nic list --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

