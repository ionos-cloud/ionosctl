---
description: List Snapshots
---

# SnapshotList

## Usage

```text
ionosctl snapshot list [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Snapshots.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name description location size secAuthProtection cpuHotPlug cpuHotUnplug ramHotPlug ramHotUnplug nicHotPlug nicHotUnplug discVirtioHotPlug discVirtioHotUnplug discScsiHotPlug discScsiHotUnplug licenceType]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
      --order-by string     Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl snapshot list
```

