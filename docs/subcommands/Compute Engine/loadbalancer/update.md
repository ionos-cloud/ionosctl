---
description: "Update a Load Balancer"
---

# LoadbalancerUpdate

## Usage

```text
ionosctl compute loadbalancer update [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update the configuration of a specified Load Balancer.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [LoadBalancerId Name Dhcp State Ip DatacenterId]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int                Level of detail for response objects (default 1)
      --dhcp                     Indicates if the Load Balancer will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false (default true)
  -F, --filters strings          Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --ip ip                    The IP of the Load Balancer
      --limit int                Maximum number of items to return per request (default 50)
  -i, --loadbalancer-id string   The unique Load Balancer Id (required)
  -n, --name string              Name of the Load Balancer
      --no-headers               Don't print table headers when table output is used
      --offset int               Number of items to skip before starting to collect the results
      --order-by string          Property to order the results by
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
  -t, --timeout int              Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                     Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute loadbalancer update --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID
```

