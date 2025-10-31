---
description: "List Application Load Balancers"
---

# ApplicationloadbalancerList

## Usage

```text
ionosctl applicationloadbalancer list [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Application Load Balancers from a specified Virtual Data Center.

Required values to run command:

* Data Center Id

## Options

```text
  -a, --all                    List all resources without the need of specifying parent ID name.
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan PrivateIps State] (default [ApplicationLoadBalancerId,Name,ListenerLan,Ips,TargetLan,PrivateIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              pagination limit: Maximum number of items to return per request (default 50)
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             Don't print table headers when table output is used
      --offset int             pagination offset: Number of items to skip before starting to collect the results
      --order-by string        Limits results to those containing a matching value for a specific property
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl applicationloadbalancer list --datacenter-id DATACENTER_ID
```

