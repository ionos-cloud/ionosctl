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

Required values to run command:

* Data Center Id
* NAT Gateway Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayRuleId Name Type Protocol SourceSubnet PublicIp TargetSubnet TargetPortRangeStart TargetPortRangeEnd State] (default [NatGatewayRuleId,Name,Protocol,SourceSubnet,PublicIp,TargetSubnet,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
      --natgateway-id string   The unique NatGateway Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl natgateway rule list --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID
```

