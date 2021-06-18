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
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayRuleId Name Type Protocol SourceSubnet PublicIp TargetSubnet TargetPortRangeStart TargetPortRangeEnd State] (default [NatGatewayRuleId,Name,Protocol,SourceSubnet,PublicIp,TargetSubnet,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for delete
      --natgateway-id string   The unique NatGateway Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -i, --rule-id string         The unique Rule Id (required)
  -t, --timeout int            Timeout option for Request for NAT Gateway Rule deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NAT Gateway Rule deletion to be executed
```

## Examples

```text
ionosctl natgateway rule delete --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID
```

