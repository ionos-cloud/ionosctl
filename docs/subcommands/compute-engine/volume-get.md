---
description: Get a Volume
---

# VolumeGet

## Usage

```text
ionosctl volume get [flags]
```

## Aliases

For `volume` command:

```text
[v vol]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information about a Volume using its ID.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             When using text output, don't print headers
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
  -i, --volume-id string       The unique Volume Id (required)
```

## Examples

```text
ionosctl volume get --datacenter-id DATACENTER_ID --volume-id VOLUME_ID
```

