---
description: Get an attached Volume from a Server
---

# ServerVolumeGet

## Usage

```text
ionosctl server volume get [flags]
```

## Aliases

For `server` command:
```text
[s svr]
```

For `volume` command:
```text
[v vol]
```

For `get` command:
```text
[g]
```

## Description

Use this command to retrieve information about an attached Volume on Server.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId SshKeys DeviceNumber UserData] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -i, --volume-id string       The unique Volume Id (required)
```

## Examples

```text
ionosctl server volume get --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --server-id 2bf04e0d-86e4-4f13-b405-442363b25e28 --volume-id 1ceb4b02-ed41-4651-a90b-9a30bc284e74 
VolumeId                               Name   Size   Type   LicenceType   State       Image
1ceb4b02-ed41-4651-a90b-9a30bc284e74   test   10GB   HDD    LINUX         AVAILABLE
```

