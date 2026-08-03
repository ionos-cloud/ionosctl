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

Use this command to start capturing the traffic passing through a NAT Gateway into an IONOS Object Storage (S3) bucket. Choose which flows to record with `--action` (by outcome) and `--direction` (relative to the gateway).

The bucket named by `--s3bucket` must already exist in your IONOS Object Storage; it is not created for you.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Bucket Name

## Options

```text
  -a, --action string          Which flows to log by outcome: ACCEPTED (allowed traffic), REJECTED (blocked traffic), or ALL (both) (default "ALL")
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -d, --direction string       Which flows to log by direction relative to the gateway: INGRESS (inbound), EGRESS (outbound), or BIDIRECTIONAL (both) (default "BIDIRECTIONAL")
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            Human-friendly name for the flowlog (default "Unnamed FlowLog")
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -b, --s3bucket string        Name of an existing IONOS Object Storage (S3) bucket where the flow records are delivered. The bucket must already exist; it is not created for you (required)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Log all traffic in both directions to an existing bucket
ionosctl compute natgateway flowlog create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name audit-all --s3bucket my-logs-bucket

# Log only rejected egress traffic
ionosctl compute natgateway flowlog create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name rejected-egress --action REJECTED --direction EGRESS --s3bucket my-logs-bucket
```

