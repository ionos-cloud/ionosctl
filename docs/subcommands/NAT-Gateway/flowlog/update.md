---
description: "Update a NAT Gateway FlowLog"
---

# NatgatewayFlowlogUpdate

## Usage

```text
ionosctl compute natgateway flowlog update [flags]
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

Use this command to change the name, capture patterns (`--action` / `--direction`) or destination bucket of an existing NAT Gateway flowlog. Only the flags you pass are changed. Any new `--s3bucket` must already exist in your IONOS Object Storage.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway FlowLog Id

## Options

```text
  -a, --action string          Which flows to log by outcome: ACCEPTED (allowed traffic), REJECTED (blocked traffic), or ALL (both)
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -d, --direction string       Which flows to log by direction relative to the gateway: INGRESS (inbound), EGRESS (outbound), or BIDIRECTIONAL (both)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -i, --flowlog-id string      The unique FlowLog Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            New human-friendly name for the flowlog
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -b, --s3bucket string        Name of an existing IONOS Object Storage (S3) bucket to deliver the flow records to. The bucket must already exist
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a flowlog
ionosctl compute natgateway flowlog update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --flowlog-id FLOWLOG_ID --name renamed-flowlog

# Narrow an existing flowlog to only rejected traffic
ionosctl compute natgateway flowlog update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --flowlog-id FLOWLOG_ID --action REJECTED
```

