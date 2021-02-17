---
description: Get an attached NIC to a Load Balancer
---

# Get

## Usage

```text
ionosctl nic attach get [flags]
```

## Description

Use this command to retrieve information about an attached NIC.

Required values to run the command:
- Data Center Id
- Load Balancer Id
- NIC Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string     The unique Data Center Id
  -h, --help                     help for get
      --ignore-stdin             Force command to execute without user input
      --loadbalancer-id string   The unique Load Balancer Id [Required flag]
      --nic-id string            The unique NIC Id [Required flag]
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id
  -v, --verbose                  Enable verbose output
```

## Examples

```text
ionosctl nic attach get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id f16dfcc1-9181-400b-a08d-7fe15ca0e9af --nic-id c7903181-daa1-4e16-a65a-e9b495c1b324
NicId                                  Name      Dhcp   LanId   Ips
c7903181-daa1-4e16-a65a-e9b495c1b324   demoNIC   true   2       []
```

## See also

* [ionosctl nic attach](./)

