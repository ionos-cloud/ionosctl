---
description: "Update a NAT Gateway FlowLog"
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
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -d, --direction string       Specifies the traffic Direction pattern
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -i, --flowlog-id string      The unique FlowLog Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            Name of the NAT Gateway FlowLog
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -b, --s3bucket string        S3 Bucket name of an existing IONOS Cloud S3 Bucket
  -t, --timeout int            Timeout option for Request for NAT Gateway FlowLog update [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for NAT Gateway FlowLog update to be executed
```

## Examples

```text
ionosctl natgateway flowlog update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID --name NAME
```

