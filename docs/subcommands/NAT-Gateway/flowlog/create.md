---
description: "Create a NAT Gateway FlowLog"
---

# NatgatewayFlowlogCreate

## Usage

```text
ionosctl compute natgateway flowlog create [flags]
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
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -d, --direction string       Specifies the traffic Direction pattern (default "BIDIRECTIONAL")
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            The name for the FlowLog (default "Unnamed FlowLog")
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -b, --s3bucket string        S3 Bucket name of an existing IONOS Cloud S3 Bucket (required)
  -t, --timeout int            Timeout option for Request for NAT Gateway FlowLog creation [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for NAT Gateway FlowLog creation to be executed
```

## Examples

```text
ionosctl compute natgateway flowlog create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME --s3bucket BUCKET_NAME
```

