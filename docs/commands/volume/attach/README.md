---
description: Attach a Volume to a Server
---

# Attach

## Usage

```text
ionosctl volume attach [flags]
```

```text
ionosctl volume attach [command]
```

## Description

Use this command to attach a Volume to a Server from a Data Center.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id

The sub-commands of `ionosctl volume attach` allow you to retrieve information about attached Volumes or about a specified attached Volume.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for attach
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
      --timeout int            Timeout option [seconds] (default 60)
  -v, --verbose                Enable verbose output
      --volume-id string       The unique Volume Id [Required flag]
      --wait                   Wait for Volume to attach to Server
```

## Examples

```text
ionosctl volume attach --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --volume-id 15546173-a100-4851-8bc4-872ec6bbee32 --wait 
Waiting for request: dfd826bb-aace-4b1e-9ae2-17901e3cc792
VolumeId                               Name         Size   Type   LicenseType   State   Image
15546173-a100-4851-8bc4-872ec6bbee32   demoVolume   10GB   HDD    LINUX         BUSY    
RequestId: dfd826bb-aace-4b1e-9ae2-17901e3cc792
Status: Command volume attach and request have been successfully executed
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl volume attach get](get.md) | Get an attached Volume from a Server |
| [ionosctl volume attach list](list.md) | List attached Volumes from a Server |

