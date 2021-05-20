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
[v]
```

## Description

Use this command to retrieve information about a Volume using its ID.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -F, --format strings         Collection of fields to be printed on output (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -h, --help                   help for get
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --volume-id string       The unique Volume Id (required)
```

## Examples

```text
ionosctl volume get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --volume-id ce510144-9bc6-4115-bd3d-b9cd232dd422 
VolumeId                               Name         Size   Type   LicenceType   State       Image
ce510144-9bc6-4115-bd3d-b9cd232dd422   demoVolume   20GB   HDD    LINUX         AVAILABLE
```

