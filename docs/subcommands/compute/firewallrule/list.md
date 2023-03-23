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
      --datacenter-id string   The unique Data Center Id (required)
      --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32      The maximum number of elements to return
      --nic-id string          The unique NIC Id (required)
      --no-headers             When using text output, don't print headers
      --order-by string        Limits results to those containing a matching value for a specific property
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl firewallrule list --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID
```

