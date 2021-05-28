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

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId SshKeys ImageAlias DeviceNumber UserData] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl volume list --datacenter-id DATACENTER_ID
```

