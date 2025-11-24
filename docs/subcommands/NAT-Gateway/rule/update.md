---
description: "Update a NAT Gateway Rule"
---

# NatgatewayRuleUpdate

## Usage

```text
ionosctl natgateway rule update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified NAT Gateway Rule from a NAT Gateway.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway Rule Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayRuleId Name Type Protocol SourceSubnet PublicIp TargetSubnet TargetPortRangeStart TargetPortRangeEnd State] (default [NatGatewayRuleId,Name,Protocol,SourceSubnet,PublicIp,TargetSubnet,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ip ip                  Public IP address of the NAT Gateway Rule
      --limit int              Pagination limit: Maximum number of items to return per request (default 50)
  -n, --name string            Name of the NAT Gateway Rule
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --port-range-end int     Target port range end associated with the NAT Gateway Rule (default 1)
      --port-range-start int   Target port range start associated with the NAT Gateway Rule (default 1)
  -p, --protocol string        Protocol of the NAT Gateway Rule. If protocol is 'ICMP' then targetPortRange start and end cannot be set
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -i, --rule-id string         The unique Rule Id (required)
      --source-subnet string   Source subnet of the NAT Gateway Rule
      --target-subnet string   Target subnet or destination subnet of the NAT Gateway Rule
  -t, --timeout int            Timeout option for Request for NAT Gateway Rule update [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for NAT Gateway Rule update to be executed
```

## Examples

```text
ionosctl natgateway rule update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID --name NAME
```

