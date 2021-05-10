---
description: Get a Server
---

# ServerGet

## Usage

```text
ionosctl server get [flags]
```

## Description

Use this command to get information about a specified Server from a Virtual Data Center. You can also wait for Server to get in AVAILABLE state using `--wait-for-state` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [ServerId,Name,AvailabilityZone,State,Cores,Ram,CpuFamily])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --force                  Force command to execute without user input
  -h, --help                   help for get
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
      --timeout int            Timeout option for waiting for Server to be in AVAILABLE state [seconds] (default 60)
      --wait-for-state         Wait for specified Server to be in AVAILABLE state
```

## Examples

```text
ionosctl server get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id f45f435e-8d6c-4170-ab90-858b59dab9ff 
ServerId                               Name         AvailabilityZone   State       Cores   Ram     CpuFamily
f45f435e-8d6c-4170-ab90-858b59dab9ff   demoServer   AUTO               AVAILABLE   4       256MB   AMD_OPTERON
```

