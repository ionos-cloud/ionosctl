---
description: "Create a Network Load Balancer"
---

# NetworkloadbalancerCreate

## Usage

```text
ionosctl networkloadbalancer create [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a Network Load Balancer in a specified Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NetworkLoadBalancerId Name ListenerLan Ips TargetLan LbPrivateIps State] (default [NetworkLoadBalancerId,Name,ListenerLan,Ips,TargetLan,LbPrivateIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
      --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            Collection of IP addresses of the Network Load Balancer
      --limit int              Maximum number of items to return per request (default 50)
      --listener-lan int       Id of the listening LAN (default 2)
  -n, --name string            Name of the Network Load Balancer (default "Network Load Balancer")
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --private-ips strings    Collection of private IP addresses with subnet mask of the Network Load Balancer
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --target-lan int         Id of the balanced private target LAN (default 1)
  -t, --timeout int            Timeout option for Request for Network Load Balancer creation [seconds] (default 300)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for Network Load Balancer creation to be executed
```

## Examples

```text
ionosctl networkloadbalancer create --datacenter-id DATACENTER_ID
```

