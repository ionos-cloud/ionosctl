---
description: Get a NAT Gateway Rule
---

# NatgatewayRuleGet

## Usage

```text
ionosctl natgateway rule get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified NAT Gateway Rule from a NAT Gateway.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway Rule Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             When using text output, don't print headers
  -i, --rule-id string         The unique Rule Id (required)
```

## Examples

```text
ionosctl natgateway rule get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID
```

