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
[s svr]
```

For `create` command:
```text
[c]
```

## Description

Use this command to create a Server in a specified Virtual Data Center. It is required that the number of cores for the Server and the amount of memory for the Server to be set.

The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways: 

* providing only the value, e.g.`--ram 256` equals 256MB.
* providing both the value and the unit, e.g.`--ram 1GB`.

You can wait for the Request to be executed using `--wait-for-request` option. You can also wait for Server to be in AVAILABLE state using `--wait-for-state` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id
* Cores
* RAM

## Options

```text
  -u, --api-url string             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -z, --availability-zone string   Availability zone of the Server (default "AUTO")
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [ServerId Name AvailabilityZone Cores Ram CpuFamily VmState State] (default [ServerId,Name,AvailabilityZone,Cores,Ram,CpuFamily,VmState,State])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                  The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits (required)
      --cpu-family string          CPU Family for the Server (default "AMD_OPTERON")
      --datacenter-id string       The unique Data Center Id (required)
  -f, --force                      Force command to execute without user input
  -h, --help                       help for create
  -n, --name string                Name of the Server
  -o, --output string              Desired output format [text|json] (default "text")
  -q, --quiet                      Quiet output
      --ram string                 The amount of memory for the Server. Size must be specified in multiples of 256 (required)
  -t, --timeout int                Timeout option for Request for Server creation/for Server to be in AVAILABLE state [seconds] (default 60)
  -w, --wait-for-request           Wait for the Request for Server creation to be executed
  -W, --wait-for-state             Wait for new Server to be in AVAILABLE state
```

## Examples

```text
ionosctl server create --datacenter-id 5d5dfefe-32cd-4f07-aa0f-4da7d503d8dd --name test --cores 2 --ram 1024MB -w -W
7.1s Waiting for request.... DONE                                                                                                                                                                          
200ms Waiting for state. DONE                                                                                                                                                                              
ServerId                               Name   AvailabilityZone   Cores   Ram      CpuFamily    VmState   State
77538d4e-1772-4b3b-a9a5-550fa118b0bd   test   AUTO               2       1024MB   INTEL_XEON   RUNNING   AVAILABLE
RequestId: b0fadc81-b4c2-41a1-b98e-4e048a2e41cf
Status: Command server create & wait have been successfully executed
```

