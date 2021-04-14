---
description: Detach a NIC from a Load Balancer
---

# DetachNetworkInterface

## Usage

```text
ionosctl loadbalancer detach-nic [flags]
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
      --cols strings             Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id [Required flag]
  -h, --help                     help for detach-nic
      --ignore-stdin             Force command to execute without user input
      --loadbalancer-id string   The unique Load Balancer Id [Required flag]
      --nic-id string            The unique NIC Id [Required flag]
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --timeout int              Timeout option for NIC to be detached from a Load Balancer [seconds] (default 60)
      --wait                     Wait for NIC to detach from a Load Balancer
```

## Examples

```text
ionosctl loadbalancer detach-nic --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --loadbalancer-id 4450e35a-e89d-4769-af60-4957c3deaf33 --nic-id 6e8faa79-1e7e-4e99-be76-f3b3179ed3c3 
Warning: Are you sure you want to detach nic from loadbalancer (y/N) ? 
y
RequestId: a2a136cd-7bce-40fe-ae53-ad0d7b322387
Status: Command loadbalancer detach-nic has been successfully executed
```

