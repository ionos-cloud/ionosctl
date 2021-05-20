---
description: Create a Server
---

# ServerCreate

## Usage

```text
ionosctl server create [flags]
```

## Aliases

For `server` command:
```text
[s]
```

## Description

Use this command to create a Server in a specified Virtual Data Center. The name, cores, ram, cpu-family and availability zone options can be set.

You can wait for the Request to be executed using `--wait-for-request` option. You can also wait for Server to be in AVAILABLE state using `--wait-for-state` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -z, --availability-zone string   Availability zone of the Server (default "AUTO")
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                  Cores option of the Server (default 2)
      --cpu-family string          CPU Family for the Server (default "AMD_OPTERON")
      --datacenter-id string       The unique Data Center Id (required)
  -f, --force                      Force command to execute without user input
  -F, --format strings             Collection of fields to be printed on output (default [ServerId,Name,AvailabilityZone,Cores,Ram,CpuFamily,VmState,State])
  -h, --help                       help for create
  -n, --name string                Name of the Server
  -o, --output string              Desired output format [text|json] (default "text")
  -q, --quiet                      Quiet output
      --ram-size int               RAM[GB] option for the Server (default 256)
  -t, --timeout int                Timeout option for Request for Server creation/for Server to be in AVAILABLE state [seconds] (default 60)
  -w, --wait-for-request           Wait for the Request for Server creation to be executed
  -W, --wait-for-state             Wait for new Server to be in AVAILABLE state
```

## Examples

```text
ionosctl server create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --name demoServer
ServerId                               Name         AvailabilityZone   State   Cores   Ram     CpuFamily
f45f435e-8d6c-4170-ab90-858b59dab9ff   demoServer   AUTO               BUSY    2       256MB   AMD_OPTERON
RequestId: 07fd3682-8642-4a5e-a57a-056e909a2af8
Status: Command server create has been successfully executed

ionosctl server create --datacenter-id 3087bf8b-3c84-405f-8b22-1978a36aa933 --name testing --wait-for-request --wait-for-state 
6.2s Waiting for request... DONE                                                                                                                                                                           
100ms Waiting for state. DONE                                                                                                                                                                              
ServerId                               Name      AvailabilityZone   State       Cores   Ram     CpuFamily
af960bf3-1585-4040-9c14-343a368339ac   testing   AUTO               AVAILABLE   2       256MB   AMD_OPTERON
RequestId: 9e6db134-284b-41a4-b581-c567c744b874
Status: Command server create & wait have been successfully executed
```

