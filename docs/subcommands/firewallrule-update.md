---
description: Update a FirewallRule
---

# FirewallruleUpdate

## Usage

```text
ionosctl firewallrule update [flags]
```

## Aliases

For `firewallrule` command:

```text
[f fr firewall]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified Firewall Rule.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [FirewallRuleId Name Protocol SourceMac SourceIP TargetIP PortRangeStart PortRangeEnd IcmpCode IcmpType Type State] (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,Type,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
  -i, --firewallrule-id string   The unique FirewallRule Id (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     help for update
      --icmp-code int            Redefine the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes
      --icmp-type int            Redefine the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types
  -n, --name string              The name for the Firewall Rule
      --nic-id string            The unique NIC Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
      --port-range-end int       Redefine the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports (default 1)
      --port-range-start int     Redefine the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports (default 1)
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id (required)
      --source-ip string         Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs
      --source-mac string        Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Not setting option allows all source MAC addresses
      --target-ip string         In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target IPs
  -t, --timeout int              Timeout option for Request for Firewall Rule update [seconds] (default 60)
      --type string              The type of Firewall Rule
  -v, --verbose                  see step by step process when running a command
  -w, --wait-for-request         Wait for Request for Firewall Rule update to be executed
```

## Examples

```text
ionosctl firewallrule update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID --name NAME --wait-for-request
```

