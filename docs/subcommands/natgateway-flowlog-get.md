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
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -i, --flowlog-id string      The unique FlowLog Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get
      --natgateway-id string   The unique NatGateway Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl natgateway flowlog get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID
```

