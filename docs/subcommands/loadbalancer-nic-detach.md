---
description: Detach a NIC from a Load Balancer
---

# LoadbalancerNicDetach

## Usage

```text
ionosctl loadbalancer nic detach [flags]
```

## Description

Use this command to detach a NIC from a Load Balancer.
You can wait for the action to be executed using `--wait` option. You can force the command to execute without user input using `--ignore-stdin` option.
Required values to run command:
* Data Center Id
* Load Balancer Id
* NIC Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Server Id on which NIC is build on. Not required, but it helps in autocompletion
      --force                    Force command to execute without user input
  -h, --help                     help for detach
      --loadbalancer-id string   The unique Load Balancer Id (required)
      --nic-id string            The unique NIC Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --timeout int              Timeout option for NIC to be removed from a Load Balancer [seconds] (default 60)
      --wait                     Wait for NIC to be removed from a Load Balancer
```

## Examples

```text
ionosctl loadbalancer nic detach --datacenter-id aa8e07a2-287a-4b45-b5e9-94761750a53c --loadbalancer-id de044efe-cfe1-41b8-9a21-966a9c03d240 --nic-id ba36c888-e966-480d-800c-77c93ec31083 
Warning: Are you sure you want to detach nic from loadbalancer (y/N) ? 
y
RequestId: 91065943-d4af-4427-aff6-ddf6a0f4ec80
Status: Command nic detach has been successfully executed
```

