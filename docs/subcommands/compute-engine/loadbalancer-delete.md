---
description: Delete a Load Balancer
---

# LoadbalancerDelete

## Usage

```text
ionosctl loadbalancer delete [flags]
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
  -u, --api-url string           Override default host url (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [LoadBalancerId Name Dhcp State Ip] (default [LoadBalancerId,Name,Dhcp,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
  -i, --loadbalancer-id string   The unique Load Balancer Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
  -t, --timeout int              Timeout option for Request for Load Balancer deletion [seconds] (default 60)
  -v, --verbose                  Print step-by-step process when running command
  -w, --wait-for-request         Wait for Request for Load Balancer deletion to be executed
```

## Examples

```text
ionosctl loadbalancer delete --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --force --wait-for-request
```

