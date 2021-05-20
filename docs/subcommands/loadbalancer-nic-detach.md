---
description: Detach a NIC from a Load Balancer
---

# LoadbalancerNicDetach

## Usage

```text
ionosctl loadbalancer nic detach [flags]
```

## Description

Use this command to remove the association of a NIC with a Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
  -f, --force                    Force command to execute without user input
  -F, --format strings           Set of fields to be printed on output (default [NicId,Name,Dhcp,LanId,Ips,State])
  -h, --help                     help for detach
      --loadbalancer-id string   The unique Load Balancer Id (required)
      --nic-id string            The unique NIC Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
  -t, --timeout int              Timeout option for Request for NIC detachment [seconds] (default 60)
  -w, --wait-for-request         Wait for the Request for NIC detachment to be executed
```

## Examples

```text
ionosctl loadbalancer nic detach --datacenter-id aa8e07a2-287a-4b45-b5e9-94761750a53c --loadbalancer-id de044efe-cfe1-41b8-9a21-966a9c03d240 --nic-id ba36c888-e966-480d-800c-77c93ec31083 
Warning: Are you sure you want to detach nic from loadbalancer (y/N) ? 
y
RequestId: 91065943-d4af-4427-aff6-ddf6a0f4ec80
Status: Command nic detach has been successfully executed
```

