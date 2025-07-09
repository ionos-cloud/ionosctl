---
description: "List Volumes"
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
* filter by property: [name type size availabilityZone image imagePassword imageAlias sshKeys bus licenceType applicationType cpuHotPlug ramHotPlug nicHotPlug nicHotUnplug discVirtioHotPlug discVirtioHotUnplug exposeSerial requireLegacyBios deviceNumber pciSlot backupunitId userData bootServer bootOrder]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* Data Center Id

## Options

```text
  -a, --all                    List all resources without the need of specifying parent ID name.
  -u, --api-url string         Override default host URL. Preferred over the config file override 'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             Don't print table headers when table output is used
      --order-by string        Limits results to those containing a matching value for a specific property
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl volume list --datacenter-id DATACENTER_ID
```

