---
description: Attach a NIC to a Load Balancer
---

# Attach

## Usage

```text
ionosctl nic attach [flags]
```

```text
ionosctl nic attach [command]
```

## Description

Use this command to attach a specified NIC to a Load Balancer on your account.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id

The sub-commands of `ionosctl nic attach` allow you to retrieve information about attached NICs or about a specified attached NIC.

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id [Required flag]
  -h, --help                     help for attach
      --ignore-stdin             Force command to execute without user input
      --loadbalancer-id string   The unique Load Balance Id [Required flag]
      --nic-id string            The unique NIC Id [Required flag]
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id
      --timeout int              Timeout option for NIC to be attached to a Load Balancer [seconds] (default 60)
      --wait                     Wait for NIC to attach to a Load Balancer
```

## Examples

```text
ionosctl nic attach --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --loadbalancer-id f16dfcc1-9181-400b-a08d-7fe15ca0e9af --nic-id c7903181-daa1-4e16-a65a-e9b495c1b324 
NicId                                  Name      Dhcp   LanId   Ips
c7903181-daa1-4e16-a65a-e9b495c1b324   demoNIC   true   1       []
RequestId: 5d892b7c-69e3-4983-ac18-a75081145d32
Status: Command nic attach has been successfully executed
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl nic attach get](get.md) | Get an attached NIC to a Load Balancer |
| [ionosctl nic attach list](list.md) | List attached NICs from a Load Balancer |

