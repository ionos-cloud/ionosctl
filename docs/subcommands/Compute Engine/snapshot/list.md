---
description: "List Snapshots"
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
* filter by property: [name description location size secAuthProtection cpuHotPlug cpuHotUnplug ramHotPlug ramHotUnplug nicHotPlug nicHotUnplug discVirtioHotPlug discVirtioHotUnplug discScsiHotPlug exposeSerial requireLegacyBios discScsiHotUnplug licenceType applicationType]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [SnapshotId Name LicenceType Size State] (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --order-by string     Limits results to those containing a matching value for a specific property
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl snapshot list
```

