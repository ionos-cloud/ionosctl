---
description: Update a Server
---

# Update

## Usage

```text
ionosctl server update [flags]
```

## Description

Use this command to update a specified Server from a Data Center.

You can wait for the action to be executed using `--wait` option.

Required values to run command:
- Data Center Id
- Server Id

## Options

```text
  -u, --api-url string            Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings              Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string      The unique Data Center Id
  -h, --help                      help for update
      --ignore-stdin              Force command to execute without user input
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
      --server-cores int          Cores option of the Server (default 2)
      --server-cpufamily string   CPU Family of the Server (default "AMD_OPTERON")
      --server-id string          The unique Server Id [Required flag]
      --server-name string        Name of the Server
      --server-ram int            RAM[GB] option for the Server (default 256)
      --server-zone string        Availability zone of the Server
      --timeout int               Timeout option [seconds] (default 60)
  -v, --verbose                   Enable verbose output
      --wait                      Wait for Server to be updated
```

## Examples

```text
ionosctl server update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id f45f435e-8d6c-4170-ab90-858b59dab9ff --server-cores 4
ServerId                               Name         AvailabilityZone   State   Cores   Ram     CpuFamily
f45f435e-8d6c-4170-ab90-858b59dab9ff   demoServer   AUTO               BUSY    4       256MB   AMD_OPTERON
RequestId: 571a1bbb-26b3-449d-9885-a20e50dc3b95
Status: Command server update has been successfully executed
```

