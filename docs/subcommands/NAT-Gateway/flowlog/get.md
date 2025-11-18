---
description: "Get a NAT Gateway FlowLog"
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
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -i, --flowlog-id string      The unique FlowLog Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl natgateway flowlog get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID
```

