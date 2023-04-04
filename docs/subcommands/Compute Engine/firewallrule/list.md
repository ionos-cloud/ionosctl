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

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name protocol sourceMac sourceIp targetIp icmpCode icmpType portRangeStart portRangeEnd type]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* Data Center Id
* Server Id
* Nic Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FirewallRuleId Name Protocol SourceMac SourceIP DestinationIP PortRangeStart PortRangeEnd IcmpCode IcmpType Direction State] (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,Direction,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -M, --max-results int32      The maximum number of elements to return
      --nic-id string          The unique NIC Id (required)
      --no-headers             When using text output, don't print headers
      --order-by string        Limits results to those containing a matching value for a specific property
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl firewallrule list --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID
```

