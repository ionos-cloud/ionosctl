---
description: List Network Load Balancer Forwarding Rules
---

# NetworkloadbalancerRuleList

## Usage

```text
ionosctl networkloadbalancer rule list [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `rule` command:

```text
[r forwardingrule]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Network Load Balancer Forwarding Rules from a specified Network Load Balancer.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name algorithm protocol listenerIp listenerPort healthCheck]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* Data Center Id
* Network Load Balancer Id

## Options

```text
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings                 Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32               The maximum number of elements to return
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      When using text output, don't print headers
      --order-by string                 Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl networkloadbalancer rule list --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID
```

