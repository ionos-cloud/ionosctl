---
description: "Create a Network Load Balancer FlowLog"
---

# NetworkloadbalancerFlowlogCreate

## Usage

```text
ionosctl compute networkloadbalancer flowlog create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create a flowlog on a Network Load Balancer. It records the NLB's connection traffic - filtered by --action (ACCEPTED / REJECTED / ALL) and --direction (INGRESS / EGRESS / BIDIRECTIONAL) - into an existing IONOS S3 bucket named by --s3bucket.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Bucket Name

## Options

```text
  -a, --action string                   Which connections to log: ACCEPTED, REJECTED, or ALL (default "ALL")
  -u, --api-url string                  Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [FlowLogId Name Action Direction Bucket State]
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int                       Level of detail for response objects (default 1)
  -d, --direction string                Which flow direction to log: INGRESS (inbound), EGRESS (outbound), or BIDIRECTIONAL (default "BIDIRECTIONAL")
  -F, --filters strings                 Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --limit int                       Maximum number of items to return per request (default 50)
  -n, --name string                     The name for the FlowLog (default "Unnamed FlowLog")
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      Don't print table headers when table output is used
      --offset int                      Number of items to skip before starting to collect the results
      --order-by string                 Property to order the results by
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
      --query string                    JMESPath query string to filter the output
  -q, --quiet                           Quiet output
  -b, --s3bucket string                 Name of an existing IONOS S3 bucket where the flow logs will be written (required)
  -t, --timeout int                     Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                   Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                            Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Log all traffic in both directions to an existing bucket
ionosctl compute networkloadbalancer flowlog create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --name "audit-all" --s3bucket my-log-bucket

# Log only rejected inbound connections (security triage)
ionosctl compute networkloadbalancer flowlog create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --name "rejected-in" --action REJECTED --direction INGRESS --s3bucket my-log-bucket
```

