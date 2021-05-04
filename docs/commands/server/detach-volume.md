---
description: Detach a Volume from a Server
---

# DetachVolume

## Usage

```text
ionosctl server detach-volume [flags]
```

## Description

This will detach the Volume from the Server. Depending on the Volume HotUnplug settings, this may result in the Server being rebooted. This will NOT delete the Volume from your Virtual Data Center. You will need to use a separate command to delete a Volume.

You can wait for the action to be executed using `--wait` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [ServerId,Name,AvailabilityZone,State,Cores,Ram,CpuFamily])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --force                  Force command to execute without user input
  -h, --help                   help for detach-volume
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
      --timeout int            Timeout option for Server to be detached from a Server [seconds] (default 60)
      --volume-id string       The unique Volume Id (required)
      --wait                   Wait for Volume to detach from Server
```

## Examples

```text
ionosctl server detach-volume --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --server-id 2bf04e0d-86e4-4f13-b405-442363b25e28 --volume-id 1ceb4b02-ed41-4651-a90b-9a30bc284e74 
Warning: Are you sure you want to detach volume from server (y/N) ? 
y
RequestId: 0fd9d6eb-25a1-496c-b0c9-bbe18a989f18
Status: Command server detach-volume has been successfully executed
```

