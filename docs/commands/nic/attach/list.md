---
description: List attached NICs from a Load Balancer
---

# List

## Usage

```text
ionosctl nic attach list [flags]
```

## Description

Use this command to get a list of attached NICs to a Load Balancer from a Data Center.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id [Required flag]
  -h, --help                     help for list
      --ignore-stdin             Force command to execute without user input
      --loadbalancer-id string   The unique Load Balancer Id [Required flag]
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id
```

## Examples

```text
ionosctl nic attach list --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id f16dfcc1-9181-400b-a08d-7fe15ca0e9af 
NicId                                  Name      Dhcp   LanId   Ips
c7903181-daa1-4e16-a65a-e9b495c1b324   demoNIC   true   2       []
```

