---
description: Delete a NAT Gateway Rule
---

# NatgatewayRuleDelete

## Usage

```text
ionosctl natgateway rule delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified NAT Gateway Rule from a NAT Gateway.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway Rule Id

## Options

```text
  -a, --all                    Delete all NAT Gateway Rules.
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --natgateway-id string   The unique NatGateway Id (required)
  -i, --rule-id string         The unique Rule Id (required)
  -t, --timeout int            Timeout option for Request for NAT Gateway Rule deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NAT Gateway Rule deletion to be executed
```

## Examples

```text
ionosctl natgateway rule delete --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID
```

