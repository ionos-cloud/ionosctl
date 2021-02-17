---
description: Detach a NIC from a Load Balancer
---

# Detach

## Usage

```text
ionosctl nic detach [flags]
```

## Description

Use this command to detach a NIC from a Load Balancer. 

You can wait for the action to be executed using `--wait` option.
You can force the command to execute without user input using `--ignore-stdin` option.

Required values to run command:
- Data Center Id
- Load Balancer Id
- NIC Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string     The unique Data Center Id
  -h, --help                     help for detach
      --ignore-stdin             Force command to execute without user input
      --loadbalancer-id string   The unique Load Balancer Id [Required flag]
      --nic-id string            The unique NIC Id [Required flag]
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id
      --timeout int              Timeout option [seconds] (default 60)
  -v, --verbose                  Enable verbose output
      --wait                     Wait for NIC to detach from a Load Balancer
```

## Examples

```text
ionosctl nic detach --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id f16dfcc1-9181-400b-a08d-7fe15ca0e9af --nic-id c7903181-daa1-4e16-a65a-e9b495c1b324 --wait 
⚠ Warning: Are you sure you want to detach nic (y/N) ? y
⧖ Waiting for request: ccfb93cb-1493-4a2c-980c-5427e15a4b74
✔ RequestId: ccfb93cb-1493-4a2c-980c-5427e15a4b74
✔ Status: Command nic detach and request have been successfully executed

ionosctl nic detach --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id f16dfcc1-9181-400b-a08d-7fe15ca0e9af --nic-id c7903181-daa1-4e16-a65a-e9b495c1b324 --wait --ignore-stdin 
⧖ Waiting for request: 1cffbd14-3d8c-4530-91d9-aa3f522a5df6
✔ RequestId: 1cffbd14-3d8c-4530-91d9-aa3f522a5df6
✔ Status: Command nic detach and request have been successfully executed
```

## See also

* [ionosctl nic](./)

