---
description: "Create an Application Load Balancer FlowLog"
---

# ApplicationloadbalancerFlowlogCreate

## Usage

```text
ionosctl compute applicationloadbalancer flowlog create [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
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

Use this command to create a flow log on a specified Application Load Balancer. The flow log streams connection metadata for the traffic the ALB handles to an existing IONOS Object Storage (S3) bucket. Choose which connections to capture with --action (ACCEPTED / REJECTED / ALL) and which direction with --direction (INGRESS / EGRESS / BIDIRECTIONAL).

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Bucket Name

## Options

```text
  -a, --action string                       Which connections to log by their disposition: ACCEPTED (allowed), REJECTED (denied), or ALL. Defaults to ALL. (default "ALL")
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [FlowLogId Name Action Direction Bucket State]
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int                           Level of detail for response objects (default 1)
  -d, --direction string                    Which traffic direction to log relative to the ALB: INGRESS (inbound), EGRESS (outbound), or BIDIRECTIONAL (both). Defaults to INGRESS. (default "INGRESS")
  -F, --filters strings                     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --limit int                           Maximum number of items to return per request (default 50)
  -n, --name string                         The name of the Application Load Balancer FlowLog. (default "Unnamed ALB Flow Log")
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Number of items to skip before starting to collect the results
      --order-by string                     Property to order the results by
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
      --query string                        JMESPath query string to filter the output
  -q, --quiet                               Quiet output
  -b, --s3bucket string                     The name of an existing IONOS Object Storage (S3) bucket that will receive the flow log records. (required)
  -t, --timeout int                         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Log all inbound traffic to an S3 bucket
ionosctl compute applicationloadbalancer flowlog create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --name "alb-ingress" --action ALL --direction INGRESS --s3bucket my-logs-bucket

# Log only rejected connections in both directions
ionosctl compute applicationloadbalancer flowlog create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --name "alb-denied" --action REJECTED --direction BIDIRECTIONAL --s3bucket my-logs-bucket --wait
```

