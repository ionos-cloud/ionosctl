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
      --datacenter-id string   The unique Data Center Id (required)
      --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -D, --destination-ip -D      In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target/destination IPs. WARNING: This short-hand flag -D is deprecated.
  -d, --direction string       The type/direction of Firewall Rule (default "INGRESS")
      --icmp-code int          Define the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes
      --icmp-type int          Define the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types
  -n, --name string            The name for the Firewall Rule (default "Unnamed Rule")
      --nic-id string          The unique NIC Id (required)
      --port-range-end int     Define the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports (default 1)
      --port-range-start int   Define the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports (default 1)
      --protocol string        The Protocol for Firewall Rule: TCP, UDP, ICMP, ANY (required)
      --server-id string       The unique Server Id (required)
      --source-ip ip           Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs
      --source-mac string      Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Unset option allows all source MAC addresses
  -t, --timeout int            Timeout option for Request for Firewall Rule creation [seconds] (default 60)
  -w, --wait-for-request       Wait for Request for Firewall Rule creation to be executed
```

## Examples

```text
ionosctl firewallrule create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --protocol PROTOCOL --direction DIRECTION --destination-ip DESTINATION_IP
```

