---
description: "Delete a Volume"
---

# VolumeDelete

## Usage

```text
ionosctl volume delete [flags]
```

## Aliases

For `volume` command:

```text
[v vol]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete specified Volume. This will result in the Volume being removed from your Virtual Data Center. Please use this with caution!

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -a, --all                    Delete all Volumes from a virtual Datacenter.
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for Volume deletion [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -i, --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Volume deletion to be executed
```

## Examples

```text
ionosctl volume delete --datacenter-id DATACENTER_ID --volume-id VOLUME_ID
```

