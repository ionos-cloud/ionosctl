---
description: Get a NAT Gateway FlowLog
---

# NatgatewayFlowlogGet

## Usage

```text
ionosctl natgateway flowlog get [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `flowlog` command:

```text
[f fl]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified NAT Gateway FlowLog from a NAT Gateway.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway FlowLog Id

## Options

```text
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --flowlog-id string      The unique FlowLog Id (required)
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             When using text output, don't print headers
```

## Examples

```text
ionosctl natgateway flowlog get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID
```

