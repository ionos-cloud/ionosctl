---
description: "Get an Application Load Balancer FlowLog"
---

# ApplicationloadbalancerFlowlogGet

## Usage

```text
ionosctl applicationloadbalancer flowlog get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Application Load Balancer FlowLog from an Application Load Balancer.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Application Load Balancer FlowLog Id

## Options

```text
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -i, --flowlog-id string                   The unique FlowLog Id (required)
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --limit int                           Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -q, --quiet                               Quiet output
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl applicationloadbalancer flowlog get --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FLOWLOG_ID
```

