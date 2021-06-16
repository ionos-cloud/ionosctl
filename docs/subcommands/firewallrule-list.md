---
description: List Firewall Rules
---

# FirewallruleList

## Usage

```text
ionosctl firewallrule list [flags]
```

## Aliases

For `firewallrule` command:
```text
[f fr firewall]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to get a list of Firewall Rules from a specified NIC from a Server.

Required values to run command:

* Data Center Id
* Server Id
*Nic Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FirewallRuleId Name Protocol SourceMac SourceIP TargetIP PortRangeStart PortRangeEnd IcmpCode IcmpType Type State] (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,Type,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl firewallrule list --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID
```

