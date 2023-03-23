---
description: Create a NAT Gateway FlowLog
---

# NatgatewayFlowlogCreate

## Usage

```text
ionosctl natgateway flowlog create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create a NAT Gateway FlowLog in a specified NAT Gateway.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Bucket Name

## Options

```text
  -a, --action string          Specifies the traffic Action pattern (default "ALL")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -d, --direction string       Specifies the traffic Direction pattern (default "BIDIRECTIONAL")
  -n, --name string            The name for the FlowLog (default "Unnamed FlowLog")
      --natgateway-id string   The unique NatGateway Id (required)
  -b, --s3bucket string        S3 Bucket name of an existing IONOS Cloud S3 Bucket (required)
  -t, --timeout int            Timeout option for Request for NAT Gateway FlowLog creation [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NAT Gateway FlowLog creation to be executed
```

## Examples

```text
ionosctl natgateway flowlog create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME --ip IP_1 --source-subnet SOURCE_SUBNET --target-subnet TARGET_SUBNET
```

