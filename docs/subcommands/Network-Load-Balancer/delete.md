---
description: "Delete a Network Load Balancer"
---

# NetworkloadbalancerDelete

## Usage

```text
ionosctl networkloadbalancer delete [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Network Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id

## Options

```text
  -a, --all                             Delete all Network Load Balancers.
  -u, --api-url string                  Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [NetworkLoadBalancerId Name ListenerLan Ips TargetLan LbPrivateIps State] (default [NetworkLoadBalancerId,Name,ListenerLan,Ips,TargetLan,LbPrivateIps,State])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --limit int                       Pagination limit: Maximum number of items to return per request (default 50)
  -i, --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      Don't print table headers when table output is used
      --offset int                      Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
  -q, --quiet                           Quiet output
  -t, --timeout int                     Timeout option for Request for Network Load Balancer deletion [seconds] (default 300)
  -v, --verbose count                   Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request                Wait for the Request for Network Load Balancer deletion to be executed
```

## Examples

```text
ionosctl networkloadbalancer delete --datacenter-id DATACENTER_ID -i NETWORKLOADBALANCER_ID
```

