---
description: Detach a Volume from a Server
---

# ServerVolumeDetach

## Usage

```text
ionosctl server volume detach [flags]
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

For `detach` command:

```text
[d]
```

## Description

This will detach the Volume from the Server. Depending on the Volume HotUnplug settings, this may result in the Server being rebooted. This will NOT delete the Volume from your Virtual Data Center. You will need to use a separate command to delete a Volume.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id

## Options

```text
  -a, --all                    Detach all Volumes.
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Volume detachment [seconds] (default 60)
  -i, --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Volume detachment to be executed
```

## Examples

```text
ionosctl server volume detach --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID
```
