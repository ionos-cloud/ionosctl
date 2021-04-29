---
description: List attached NICs from a Load Balancer
---

# ListNetworkInterfaces

## Usage

```text
ionosctl loadbalancer list-nics [flags]
```

## Description

Use this command to get a list of attached NICs to a Load Balancer from a Data Center.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [LoadBalancerId,Name,Dhcp])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
      --force                    Force command to execute without user input
  -h, --help                     help for list-nics
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
```

## Examples

```text
ionosctl loadbalancer list-nics --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --loadbalancer-id 4450e35a-e89d-4769-af60-4957c3deaf33 
NicId                                  Name   Dhcp   LanId   Ips
6e8faa79-1e7e-4e99-be76-f3b3179ed3c3   test   true   2       []
```

