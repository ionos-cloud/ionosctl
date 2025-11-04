---
description: "Get a Load Balancer"
---

# LoadbalancerGet

## Usage

```text
ionosctl loadbalancer get [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information about a Load Balancer instance.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [LoadBalancerId Name Dhcp State Ip DatacenterId] (default [LoadBalancerId,Name,Dhcp,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int32              Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --limit int                Pagination limit: Maximum number of items to return per request (default 50)
  -i, --loadbalancer-id string   The unique Load Balancer Id (required)
      --no-headers               Don't print table headers when table output is used
      --offset int               Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string            Desired output format [text|json|api-json] (default "text")
  -q, --quiet                    Quiet output
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl loadbalancer get --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID
```

