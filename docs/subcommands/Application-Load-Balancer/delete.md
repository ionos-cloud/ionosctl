---
description: "Delete an Application Load Balancer"
---

# ApplicationloadbalancerDelete

## Usage

```text
ionosctl applicationloadbalancer delete [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Application Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` or `-w` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id

## Options

```text
  -a, --all                                 Delete all Application Load Balancers
  -u, --api-url string                      Override default host URL. If set, this will be preferred over the config file override. If unset, the default will only be used as a fallback (default "https://api.ionos.com")
  -i, --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan PrivateIps State] (default [ApplicationLoadBalancerId,Name,ListenerLan,Ips,TargetLan,PrivateIps,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --no-headers                          Don't print table headers when table output is used
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -q, --quiet                               Quiet output
  -t, --timeout int                         Timeout option for Request for Application Load Balancer deletion [seconds] (default 300)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer deletion to be executed
```

## Examples

```text
ionosctl applicationloadbalancer delete --datacenter-id DATACENTER_ID -i APPLICATIONLOADBALANCER_ID
```

