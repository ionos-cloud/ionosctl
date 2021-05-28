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

Use this command to create a Server in a specified Virtual Data Center. The name, cores, ram, cpu-family and availability zone options can be set.

You can wait for the Request to be executed using `--wait-for-request` option. You can also wait for Server to be in AVAILABLE state using `--wait-for-state` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -z, --availability-zone string   Availability zone of the Server (default "AUTO")
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [ServerId Name AvailabilityZone Cores Ram CpuFamily VmState State] (default [ServerId,Name,AvailabilityZone,Cores,Ram,CpuFamily,VmState,State])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                  Cores option of the Server (default 2)
      --cpu-family string          CPU Family for the Server (default "AMD_OPTERON")
      --datacenter-id string       The unique Data Center Id (required)
  -f, --force                      Force command to execute without user input
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
ionosctl server create --datacenter-id DATACENTER_ID --name NAME --cores 2 --ram 512MB -w -W
```

