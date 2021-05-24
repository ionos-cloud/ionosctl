---
description: Create a Firewall Rule
---

# FirewallruleCreate

## Usage

```text
ionosctl firewallrule create [flags]
```

## Aliases

For `firewallrule` command:
```text
[f fr firewall]
```

## Description

Use this command to create/add a new Firewall Rule to the specified NIC. All Firewall Rules must be associated with a NIC.

NOTE: the Firewall Rule Protocol can only be set when creating a new Firewall Rule.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id 
* Protocol

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FirewallRuleId Name Protocol SourceMac SourceIP TargetIP PortRangeStart PortRangeEnd State] (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for create
      --icmp-code int          Define the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes.
      --icmp-type int          Define the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types.
  -n, --name string            The name for the Firewall Rule
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
      --port-range-end int     Define the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports. (default 1)
      --port-range-start int   Define the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports. (default 1)
      --protocol string        The Protocol for Firewall Rule: TCP, UDP, ICMP, ANY (required)
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
      --source-ip string       Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs.
      --source-mac string      Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Unset option allows all source MAC addresses.
      --target-ip string       In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target IPs.
  -t, --timeout int            Timeout option for Request for Firewall Rule creation [seconds] (default 60)
  -w, --wait-for-request       Wait for Request for Firewall Rule creation to be executed
```

## Examples

```text
ionosctl firewallrule create --datacenter-id f2d82ba9-7dc4-4945-89b6-3d194f6be29b --server-id d776e064-a3f9-4fbd-8729-93818b7459bb --nic-id 029c05a4-f5f7-4398-9469-2eb3d6db3460 --protocol TCP --name demo --port-range-start 2476 --port-range-end 2476
FirewallRuleId                         Name   Protocol   PortRangeStart   PortRangeEnd   State
4221e2c8-0316-447c-aeed-69ac92e585be   demo   TCP        2476             2476           BUSY
RequestId: 09a47137-e377-4a79-b2b9-16744e298ad5
Status: Command firewallrule create has been successfully executed
```

