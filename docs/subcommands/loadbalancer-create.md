---
description: Create a Load Balancer
---

# LoadbalancerCreate

## Usage

```text
ionosctl loadbalancer create [flags]
```

## Aliases

For `loadbalancer` command:
```text
[lb]
```

For `create` command:
```text
[c]
```

## Description

Use this command to create a new Load Balancer within the Virtual Data Center. Load balancers can be used for public or private IP traffic. The name, IP and DHCP for the Load Balancer can be set.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LoadBalancerId Name Dhcp State Ip] (default [LoadBalancerId,Name,Dhcp,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --dhcp                   Indicates if the Load Balancer will reserve an IP using DHCP (default true)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for create
  -n, --name string            Name of the Load Balancer
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for Load Balancer creation [seconds] (default 60)
  -w, --wait-for-request       Wait for Request for Load Balancer creation to be executed
```

## Examples

```text
ionosctl loadbalancer create --datacenter-id DATACENTER_ID --name NAME
```

