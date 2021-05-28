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
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId SshKeys ImageAlias DeviceNumber UserData] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -i, --volume-id string       The unique Volume Id (required)
```

## Examples

```text
ionosctl volume get --datacenter-id DATACENTER_ID --volume-id VOLUME_ID
```

