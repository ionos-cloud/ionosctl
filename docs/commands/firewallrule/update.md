---
description: Update a FirewallRule
---

# Update

## Usage

```text
ionosctl firewallrule update [flags]
```

## Description

Use this command to update a specified Firewall Rule.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id

## Options

```text
  -u, --api-url string                      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings                        Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id [Required flag]
      --firewallrule-icmp-code string       Redefine the allowed code (from 0 to 254) if protocol ICMP is chosen. Unset option allows all codes.
      --firewallrule-icmp-type string       Redefine the allowed type (from 0 to 254) if the protocol ICMP is chosen. Unset option allows all types.
      --firewallrule-id string              The unique FirewallRule Id [Required flag]
      --firewallrule-name string            The name for the Firewall Rule
      --firewallrule-port-range-end int     Redefine the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd unset to allow all ports.
      --firewallrule-port-range-start int   Redefine the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd unset to allow all ports.
      --firewallrule-source-ip string       Only traffic originating from the respective IPv4 address is allowed. Unset option allows all source IPs.
      --firewallrule-source-mac string      Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Unset option allows all source MAC addresses.
      --firewallrule-target-ip string       In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Unset option allows all target IPs.
  -h, --help                                help for update
      --ignore-stdin                        Force command to execute without user input
      --nic-id string                       The unique NIC Id [Required flag]
  -o, --output string                       Desired output format [text|json] (default "text")
  -q, --quiet                               Quiet output
      --server-id string                    The unique Server Id [Required flag]
      --timeout int                         Timeout option for Firewall Rule to be updated [seconds] (default 60)
      --wait                                Wait for Firewall Rule to be updated
```

## Examples

```text
ionosctl firewallrule update --datacenter-id f2d82ba9-7dc4-4945-89b6-3d194f6be29b --server-id d776e064-a3f9-4fbd-8729-93818b7459bb --nic-id 029c05a4-f5f7-4398-9469-2eb3d6db3460 --firewallrule-id 4221e2c8-0316-447c-aeed-69ac92e585be --firewallrule-name new-test --wait 
Waiting for request: 2e3d6e81-2830-4d68-82ff-daee6f115864
FirewallRuleId                         Name       Protocol   PortRangeStart   PortRangeEnd   State
4221e2c8-0316-447c-aeed-69ac92e585be   new-test   TCP        2476             2476           BUSY
RequestId: 2e3d6e81-2830-4d68-82ff-daee6f115864
Status: Command firewallrule update and request have been successfully executed
```

