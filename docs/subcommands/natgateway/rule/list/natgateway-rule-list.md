---
description: List NAT Gateway Rules
---

# NatgatewayRuleList

## Usage

```text
ionosctl natgateway rule list [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `rule` command:

```text
[r]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list NAT Gateway Rules from a specified NAT Gateway.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name type protocol sourceSubnet publicIp targetSubnet targetPortRange]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* Data Center Id
* NAT Gateway Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32      The maximum number of elements to return
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             When using text output, don't print headers
      --order-by string        Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl natgateway rule list --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID
```

