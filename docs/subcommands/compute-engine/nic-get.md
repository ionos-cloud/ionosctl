---
description: Get a NIC
---

# NicGet

## Usage

```text
ionosctl nic get [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified NIC from specified Data Center and Server.

Required values to run command:

* Data Center Id
* Server Id
* NIC Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Name Dhcp LanId Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac] (default [NicId,Name,Dhcp,LanId,Ips,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -i, --nic-id string          The unique NIC Id (required)
      --no-headers             When using text output, don't print headers
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl nic get --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID
```

