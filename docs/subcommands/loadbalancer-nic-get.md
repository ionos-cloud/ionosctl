---
description: Get an attached NIC to a Load Balancer
---

# LoadbalancerNicGet

## Usage

```text
ionosctl loadbalancer nic get [flags]
```

## Description

Use this command to retrieve the attributes of a given load balanced NIC.

Required values to run the command:

* Data Center Id
* Load Balancer Id
* NIC Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
      --force                    Force command to execute without user input
  -h, --help                     help for get
      --loadbalancer-id string   The unique Load Balancer Id (required)
      --nic-id string            The unique NIC Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
```

## Examples

```text
ionosctl loadbalancer nic get --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --loadbalancer-id 4450e35a-e89d-4769-af60-4957c3deaf33 --nic-id 6e8faa79-1e7e-4e99-be76-f3b3179ed3c3 
NicId                                  Name   Dhcp   LanId   Ips
6e8faa79-1e7e-4e99-be76-f3b3179ed3c3   test   true   2       []
```

