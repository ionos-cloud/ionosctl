---
description: List Servers
---

# ServerList

## Usage

```text
ionosctl server list [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Servers from a specified Virtual Data Center.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [templateUuid name cores ram placementGroupId availabilityZone vmState bootCdrom bootVolume cpuFamily type]
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
ionosctl server list --datacenter-id DATACENTER_ID
```
