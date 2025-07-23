---
description: "Get an Application Load Balancer"
---

# ApplicationloadbalancerGet

## Usage

```text
ionosctl applicationloadbalancer get [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Application Load Balancer from a Virtual Data Center. You can also wait for Application Load Balancer to get in AVAILABLE state using `--wait-for-state` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id

## Options

```text
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud'|'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
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
  -t, --timeout int                         Timeout option for waiting for Application Load Balancer to be in AVAILABLE state [seconds] (default 300)
  -v, --verbose                             Print step-by-step process when running command
  -W, --wait-for-state                      Wait for specified Application Load Balancer to be in AVAILABLE state
```

## Examples

```text
ionosctl applicationloadbalancer get --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID
```

