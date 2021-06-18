---
description: List Load Balancers
---

# LoadbalancerList

## Usage

```text
ionosctl loadbalancer list [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of Load Balancers within a Virtual Data Center on your account.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LoadBalancerId Name Dhcp State Ip] (default [LoadBalancerId,Name,Dhcp,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl loadbalancer list --datacenter-id DATACENTER_ID
```

