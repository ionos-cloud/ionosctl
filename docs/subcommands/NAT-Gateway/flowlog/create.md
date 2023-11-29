---
description: "Create a NAT Gateway FlowLog"
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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -d, --direction string       Specifies the traffic Direction pattern (default "BIDIRECTIONAL")
  -h, --help                   Print usage
  -n, --name string            The name for the FlowLog (default "Unnamed FlowLog")
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -b, --s3bucket string        S3 Bucket name of an existing IONOS Cloud S3 Bucket (required)
  -t, --timeout int            Timeout option for Request for NAT Gateway FlowLog creation [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait-for-request       Wait for the Request for NAT Gateway FlowLog creation to be executed
```

## Examples

```text
ionosctl natgateway flowlog create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME --ip IP_1 --source-subnet SOURCE_SUBNET --target-subnet TARGET_SUBNET
```

