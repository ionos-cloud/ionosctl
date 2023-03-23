---
description: List Volumes
---

# VolumeList

## Usage

```text
ionosctl volume list [flags]
```

## Aliases

For `volume` command:

```text
[v vol]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list all Volumes from a Data Center on your account.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name type size availabilityZone image imagePassword imageAlias sshKeys bus licenceType cpuHotPlug ramHotPlug nicHotPlug nicHotUnplug discVirtioHotPlug discVirtioHotUnplug deviceNumber pciSlot backupunitId userData bootServer bootOrder]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* Data Center Id

## Options

```text
  -a, --all                    List all resources without the need of specifying parent ID name.
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             When using text output, don't print headers
      --order-by string        Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl volume list --datacenter-id DATACENTER_ID
```

