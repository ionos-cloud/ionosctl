---
description: List attached Volumes from a Server
---

# ServerVolumeList

## Usage

```text
ionosctl server volume list [flags]
```

## Aliases

For `server` command:
```text
[s]
```

For `volume` command:
```text
[v]
```

## Description

Use this command to retrieve a list of Volumes attached to the Server.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId SshKeys ImageAlias DeviceNumber UserData] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl server volume list --datacenter-id aa8e07a2-287a-4b45-b5e9-94761750a53c --server-id 1dc7c6a8-5ab3-4fa8-83e7-9d989bd52ffa 
VolumeId                               Name   Size   Type   LicenceType   State       Image
101291d1-2227-432a-9773-97b5ace7b8d3   test   10GB   HDD    LINUX         AVAILABLE
```

