---
description: Create a Server
---

# Create

## Usage

```text
ionosctl server create [flags]
```

## Description

Use this command to create a Server in a specified Data Center. The name, cores, ram, cpu-family and availability zone options can be set.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings               Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string       The unique Data Center Id [Required flag]
      --force                      Force command to execute without user input
  -h, --help                       help for create
  -o, --output string              Desired output format [text|json] (default "text")
  -q, --quiet                      Quiet output
      --server-cores int           Cores option of the Server (default 2)
      --server-cpu-family string   CPU Family for the Server (default "AMD_OPTERON")
      --server-name string         Name of the Server
      --server-ram int             RAM[GB] option for the Server (default 256)
      --server-zone string         Availability zone of the Server (default "AUTO")
      --timeout int                Timeout option for Server to be created [seconds] (default 60)
      --wait                       Wait for Server to be created
```

## Examples

```text
ionosctl server create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-name demoServer
ServerId                               Name         AvailabilityZone   State   Cores   Ram     CpuFamily
f45f435e-8d6c-4170-ab90-858b59dab9ff   demoServer   AUTO               BUSY    2       256MB   AMD_OPTERON
RequestId: 07fd3682-8642-4a5e-a57a-056e909a2af8
Status: Command server create has been successfully executed

ionosctl server create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-name demoServer --wait 
Waiting for request: e9d12f57-3513-4ae3-ab39-179aacb8c072
ServerId                               Name         AvailabilityZone   State   Cores   Ram     CpuFamily
35201d04-0ea2-43e7-abc4-56f92737bb9d   demoServer                      BUSY    2       256MB   AMD_OPTERON
RequestId: e9d12f57-3513-4ae3-ab39-179aacb8c072
Status: Command server create and request have been successfully executed
```

