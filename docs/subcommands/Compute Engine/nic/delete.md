---
description: "Delete a NIC"
---

# NicDelete

## Usage

```text
ionosctl nic delete [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `delete` command:

```text
[d]
```

## Description

This command deletes a specified NIC.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* NIC Id

## Options

```text
  -a, --all                    Delete all the Nics from a Server.
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Name Dhcp LanId Ips IPv6Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac DHCPv6 IPv6CidrBlock] (default [NicId,Name,Dhcp,LanId,Ips,IPv6Ips,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -i, --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for NIC deletion [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait-for-request       Wait for the Request for NIC deletion to be executed
```

## Examples

```text
ionosctl nic delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --force
```

