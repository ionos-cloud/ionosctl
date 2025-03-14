---
description: "Update a Load Balancer"
---

# LoadbalancerUpdate

## Usage

```text
ionosctl loadbalancer update [flags]
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

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -u, --api-url string           Override default host url (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [LoadBalancerId Name Dhcp State Ip DatacenterId] (default [LoadBalancerId,Name,Dhcp,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int32              Controls the detail depth of the response objects. Max depth is 10.
      --dhcp                     Indicates if the Load Balancer will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false (default true)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --ip ip                    The IP of the Load Balancer
  -i, --loadbalancer-id string   The unique Load Balancer Id (required)
  -n, --name string              Name of the Load Balancer
      --no-headers               Don't print table headers when table output is used
  -o, --output string            Desired output format [text|json|api-json] (default "text")
  -q, --quiet                    Quiet output
  -t, --timeout int              Timeout option for Request for Load Balancer update [seconds] (default 60)
  -v, --verbose                  Print step-by-step process when running command
  -w, --wait                     Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl loadbalancer update --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --dhcp=false --wait-for-request
```

