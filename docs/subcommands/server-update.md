---
description: Update a Server
---

# ServerUpdate

## Usage

```text
ionosctl server update [flags]
```

## Aliases

For `server` command:
```text
[s svr]
```

## Description

Use this command to update a specified Server from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can also wait for Server to be in AVAILABLE state using `--wait-for-state` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -z, --availability-zone string   Availability zone of the Server
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [ServerId Name AvailabilityZone Cores Ram CpuFamily VmState State] (default [ServerId,Name,AvailabilityZone,Cores,Ram,CpuFamily,VmState,State])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                  Cores option of the Server (default 2)
      --cpu-family string          CPU Family of the Server (default "AMD_OPTERON")
      --datacenter-id string       The unique Data Center Id (required)
  -f, --force                      Force command to execute without user input
  -h, --help                       help for update
  -n, --name string                Name of the Server
  -o, --output string              Desired output format [text|json] (default "text")
  -q, --quiet                      Quiet output
      --ram-size int               RAM[GB] option for the Server (default 256)
  -i, --server-id string           The unique Server Id (required)
  -t, --timeout int                Timeout option for Request for Server update/for Server to be in AVAILABLE state [seconds] (default 60)
  -w, --wait-for-request           Wait for the Request for Server update to be executed
  -W, --wait-for-state             Wait for the updated Server to be in AVAILABLE state
```

## Examples

```text
ionosctl server update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id f45f435e-8d6c-4170-ab90-858b59dab9ff --cores 4
ServerId                               Name         AvailabilityZone   State   Cores   Ram     CpuFamily
f45f435e-8d6c-4170-ab90-858b59dab9ff   demoServer   AUTO               BUSY    4       256MB   AMD_OPTERON
RequestId: 571a1bbb-26b3-449d-9885-a20e50dc3b95
Status: Command server update has been successfully executed
```

