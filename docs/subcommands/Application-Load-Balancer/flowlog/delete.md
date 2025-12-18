---
description: "Delete an Application Load Balancer FlowLog"
---

# ApplicationloadbalancerFlowlogDelete

## Usage

```text
ionosctl applicationloadbalancer flowlog delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Application Load Balancer FlowLog from an Application Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Application Load Balancer FlowLog Id

## Options

```text
  -a, --all                                 Delete all Application Load Balancer FlowLogs
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan PrivateIps State] (default [ApplicationLoadBalancerId,Name,ListenerLan,Ips,TargetLan,PrivateIps,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int                           Level of detail for response objects (default 1)
  -F, --filters strings                     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -i, --flowlog-id string                   The unique FlowLog Id (required)
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --limit int                           Maximum number of items to return per request (default 50)
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Number of items to skip before starting to collect the results
      --order-by string                     Property to order the results by
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
      --query string                        JMESPath query string to filter the output
  -q, --quiet                               Quiet output
  -t, --timeout int                         Timeout option for Request for Application Load Balancer FlowLog deletion [seconds] (default 300)
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer FlowLog deletion to be executed
```

## Examples

```text
ionosctl applicationloadbalancer flowlog delete --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FLOWLOG_ID
```

