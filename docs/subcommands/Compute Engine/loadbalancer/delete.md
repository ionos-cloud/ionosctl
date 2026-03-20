---
description: "Delete a Load Balancer"
---

# LoadbalancerDelete

## Usage

```text
ionosctl compute loadbalancer delete [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete the specified Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -a, --all                      Delete all the LoadBlancers from a virtual Datacenter.
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [LoadBalancerId Name Dhcp State Ip DatacenterId] (default [LoadBalancerId,Name,Dhcp,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int                Level of detail for response objects (default 1)
  -F, --filters strings          Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --limit int                Maximum number of items to return per request (default 50)
  -i, --loadbalancer-id string   The unique Load Balancer Id (required)
      --no-headers               Don't print table headers when table output is used
      --offset int               Number of items to skip before starting to collect the results
      --order-by string          Property to order the results by
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
  -t, --timeout int              Timeout option for Request for Load Balancer deletion [seconds] (default 60)
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request         Wait for Request for Load Balancer deletion to be executed
```

## Examples

```text
ionosctl compute loadbalancer delete --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID
```

