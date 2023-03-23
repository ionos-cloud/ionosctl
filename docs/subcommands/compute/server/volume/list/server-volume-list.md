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
[s svr]
```

For `volume` command:

```text
[v vol]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of Volumes attached to the Server.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name type size availabilityZone image imagePassword imageAlias sshKeys bus licenceType cpuHotPlug ramHotPlug nicHotPlug nicHotUnplug discVirtioHotPlug discVirtioHotUnplug deviceNumber pciSlot backupunitId userData bootServer bootOrder]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             When using text output, don't print headers
      --order-by string        Limits results to those containing a matching value for a specific property
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl server volume list --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

