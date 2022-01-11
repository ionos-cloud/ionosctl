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

For `create` command:

```text
[c]
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
  -u, --api-url string          Override default host url (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [FirewallRuleId Name Protocol SourceMac SourceIP DestinationIP PortRangeStart PortRangeEnd IcmpCode IcmpType Direction State] (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,Direction,State])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string    The unique Data Center Id (required)
  -D, --destination-ip string   In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target/destination IPs
  -d, --direction string        The type/direction of Firewall Rule (default "INGRESS")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --icmp-code int           Define the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes
      --icmp-type int           Define the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types
  -n, --name string             The name for the Firewall Rule (default "Unnamed Rule")
      --nic-id string           The unique NIC Id (required)
  -o, --output string           Desired output format [text|json] (default "text")
      --port-range-end int      Define the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports (default 1)
      --port-range-start int    Define the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports (default 1)
      --protocol string         The Protocol for Firewall Rule: TCP, UDP, ICMP, ANY (required)
  -q, --quiet                   Quiet output
      --server-id string        The unique Server Id (required)
      --source-ip string        Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs
      --source-mac string       Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Unset option allows all source MAC addresses
  -t, --timeout int             Timeout option for Request for Firewall Rule creation [seconds] (default 60)
  -v, --verbose                 Print step-by-step process when running command
  -w, --wait-for-request        Wait for Request for Firewall Rule creation to be executed
```

## Examples

```text
ionosctl firewallrule create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --protocol PROTOCOL --direction DIRECTION --destination-ip DESTINATION_IP
```

