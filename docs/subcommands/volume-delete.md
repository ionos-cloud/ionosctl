---
description: Delete a Volume
---

# VolumeDelete

## Usage

```text
ionosctl volume delete [flags]
```

## Aliases

For `volume` command:
```text
[v vol]
```

## Description

Use this command to delete specified Volume. This will result in the Volume being removed from your Virtual Data Center. Please use this with caution!

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId SshKeys ImageAlias DeviceNumber UserData] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for delete
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for Volume deletion [seconds] (default 60)
  -i, --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Volume deletion to be executed
```

## Examples

```text
ionosctl volume delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --volume-id ce510144-9bc6-4115-bd3d-b9cd232dd422 
Warning: Are you sure you want to delete volume (y/N) ? y
RequestId: 6958b90b-54fa-4967-8be2-e32412559f9c
Status: Command volume delete has been successfully executed
```

