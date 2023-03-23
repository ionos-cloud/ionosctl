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
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --server-id string       The unique Server Id (required)
  -i, --volume-id string       The unique Volume Id (required)
```

## Examples

```text
ionosctl server volume get --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID
```

