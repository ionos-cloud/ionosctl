---
description: Attach a NIC to a Load Balancer
---

# LoadbalancerNicAttach

## Usage

```text
ionosctl loadbalancer nic attach [flags]
```

## Description

Use this command to add/attach a specified NIC to a Load Balancer on your account.
You can wait for the action to be executed using `--wait` option.
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
  -h, --help                     help for attach
      --loadbalancer-id string   The unique Load Balancer Id (required)
      --nic-id string            The unique NIC Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id on which NIC is build on. Not required, but it helps in autocompletion
      --timeout int              Timeout option for NIC to be added to a Load Balancer [seconds] (default 60)
      --wait                     Wait for NIC to be added to a Load Balancer
```

## Examples

```text
ionosctl loadbalancer nic attach --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --server-id 2bf04e0d-86e4-4f13-b405-442363b25e28 --nic-id 6e8faa79-1e7e-4e99-be76-f3b3179ed3c3 --loadbalancer-id 4450e35a-e89d-4769-af60-4957c3deaf33 
NicId                                  Name   Dhcp   LanId   Ips
6e8faa79-1e7e-4e99-be76-f3b3179ed3c3   test   true   1       []
RequestId: 01b8468f-b489-40af-a4fd-3606d06da8d7
Status: Command nic attach has been successfully executed
```

