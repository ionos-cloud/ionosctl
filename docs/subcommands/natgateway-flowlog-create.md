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
* Name
* Direction
* Bucket Name

## Options

```text
  -a, --action string          Specifies the traffic Action pattern (default "ALL")
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
  -b, --bucket-name string     S3 Bucket name of an existing IONOS Cloud S3 Bucket (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -d, --direction string       Specifies the traffic Direction pattern (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for create
  -n, --name string            The name for the FlowLog (required)
      --natgateway-id string   The unique NatGateway Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for NAT Gateway FlowLog creation [seconds] (default 60)
  -v, --verbose                see step by step process when running a command
  -w, --wait-for-request       Wait for the Request for NAT Gateway FlowLog creation to be executed
```

## Examples

```text
ionosctl natgateway flowlog create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME --ip IP_1 --source-subnet SOURCE_SUBNET --target-subnet TARGET_SUBNET
```

