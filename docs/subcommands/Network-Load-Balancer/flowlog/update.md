---
description: "Update a Network Load Balancer FlowLog"
---

# NetworkloadbalancerFlowlogUpdate

## Usage

```text
ionosctl compute networkloadbalancer flowlog update [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
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

Use this command to update a flowlog's name, capture filters (--action, --direction), or destination bucket (--s3bucket). Only the flags you pass are changed.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Network Load Balancer FlowLog Id

## Options

```text
  -a, --action string                   Which connections to log: ACCEPTED, REJECTED, or ALL
  -u, --api-url string                  Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [FlowLogId Name Action Direction Bucket State]
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int                       Level of detail for response objects (default 1)
  -d, --direction string                Which flow direction to log: INGRESS (inbound), EGRESS (outbound), or BIDIRECTIONAL
  -F, --filters strings                 Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -i, --flowlog-id string               The unique FlowLog Id (required)
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --limit int                       Maximum number of items to return per request (default 50)
  -n, --name string                     Name of the Network Load Balancer FlowLog
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      Don't print table headers when table output is used
      --offset int                      Number of items to skip before starting to collect the results
      --order-by string                 Property to order the results by
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
      --query string                    JMESPath query string to filter the output
  -q, --quiet                           Quiet output
  -b, --s3bucket string                 Name of an existing IONOS S3 bucket where the flow logs will be written
  -t, --timeout int                     Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                   Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                            Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Narrow an existing flowlog to log only rejected connections
ionosctl compute networkloadbalancer flowlog update --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID --action REJECTED

# Point the flowlog at a different bucket
ionosctl compute networkloadbalancer flowlog update --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID --s3bucket new-log-bucket
```

