---
description: Update a NAT Gateway FlowLog
---

# NatgatewayFlowlogUpdate

## Usage

```text
ionosctl natgateway flowlog update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified NAT Gateway FlowLog from a NAT Gateway.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway FlowLog Id

## Options

```text
  -a, --action string          Specifies the traffic Action pattern
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -d, --direction string       Specifies the traffic Direction pattern
  -i, --flowlog-id string      The unique FlowLog Id (required)
  -n, --name string            Name of the NAT Gateway FlowLog
      --natgateway-id string   The unique NatGateway Id (required)
  -b, --s3bucket string        S3 Bucket name of an existing IONOS Cloud S3 Bucket
  -t, --timeout int            Timeout option for Request for NAT Gateway FlowLog update [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NAT Gateway FlowLog update to be executed
```

## Examples

```text
ionosctl natgateway flowlog update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID --name NAME
```

